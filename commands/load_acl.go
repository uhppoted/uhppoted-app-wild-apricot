package commands

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/uhppoted/uhppote-core/uhppote"
	lib "github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/config"
	"github.com/uhppoted/uhppoted-lib/lockfile"

	"github.com/uhppoted/uhppoted-app-wild-apricot/acl"
	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

var LoadACLCmd = LoadACL{
	workdir:     DEFAULT_WORKDIR,
	credentials: filepath.Join(DEFAULT_CONFIG_DIR, ".wild-apricot", "credentials.json"),
	rules:       filepath.Join(DEFAULT_CONFIG_DIR, "wild-apricot.grl"),
	withPIN:     false,
	force:       false,
	strict:      false,
	dryrun:      false,
	lockfile:    "",
	debug:       false,
}

type LoadACL struct {
	workdir     string
	credentials string
	rules       string
	withPIN     bool
	force       bool
	strict      bool
	dryrun      bool
	logfile     string
	rptfile     string
	lockfile    string
	debug       bool
}

func (cmd *LoadACL) Name() string {
	return "load-acl"
}

func (cmd *LoadACL) Description() string {
	return "Retrieves an access control list from a Wild Apricot member database and updates the card lists on the configured controllers"
}

func (cmd *LoadACL) Usage() string {
	return "--credentials <file> --rules <url>"
}

func (cmd *LoadACL) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s [--debug] [--config <file>] load-acl [--credentials <file>] [--rules <url>] [--log <file>]\n", APP)
	fmt.Println()
	fmt.Println("  Downloads an access control list from a Wild Apricot member database, applies the ACL rules and updates the card lists")
	fmt.Println("  on the configured controllers")
	fmt.Println()

	helpOptions(cmd.FlagSet())

	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    uhppote-app-wild-apricot --debug --config uhppoted.conf load-acl --credentials ".credentials/wild-apricot.json"`)
	fmt.Println()
}

func (cmd *LoadACL) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("load-acl", flag.ExitOnError)
	lockfile := filepath.Join(cmd.workdir, ".wild-apricot", "uhppoted-app-wild-apricot.lock")

	flagset.StringVar(&cmd.workdir, "workdir", cmd.workdir, "Directory for working files (tokens, revisions, etc)'")
	flagset.StringVar(&cmd.credentials, "credentials", cmd.credentials, "Path for the 'credentials.json' file. Defaults to "+cmd.credentials)
	flagset.StringVar(&cmd.rules, "rules", cmd.rules, "URI for the 'grule' rules file. Support file path, HTTP and HTTPS. Defaults to "+cmd.rules)
	flagset.BoolVar(&cmd.withPIN, "with-pin", cmd.withPIN, "Updates the card keypad PIN code on the access controllers")
	flagset.BoolVar(&cmd.force, "force", cmd.force, "Forces an update, overriding the  version and compare logic")
	flagset.BoolVar(&cmd.strict, "strict", cmd.strict, "Fails with an error if the members list contains duplicate card numbers")
	flagset.BoolVar(&cmd.dryrun, "dry-run", cmd.dryrun, "Simulates a load-acl without making any changes to the access controllers")
	flagset.StringVar(&cmd.logfile, "log", cmd.logfile, "File to which the (optional) summary report is appended")
	flagset.StringVar(&cmd.rptfile, "report", cmd.rptfile, "File to which the detail report is written. Defaults to stdout if not provided")
	flagset.StringVar(&cmd.lockfile, "lockfile", cmd.lockfile, fmt.Sprintf("Filepath for lock file. Defaults to %v", lockfile))

	return flagset
}

func (cmd *LoadACL) Execute(args ...interface{}) error {
	timestamp := time.Now()
	options := args[0].(*Options)

	cmd.debug = options.Debug

	// ... check parameters
	if strings.TrimSpace(cmd.credentials) == "" {
		return fmt.Errorf("invalid credentials file")
	}

	if strings.TrimSpace(cmd.rules) == "" {
		return fmt.Errorf("invalid rules file")
	}

	// ... locked?
	lockFile := config.Lockfile{
		File:   filepath.Join(cmd.workdir, ".wild-apricot", "uhppoted-app-wild-apricot.lock"),
		Remove: lockfile.RemoveLockfile,
	}

	if cmd.lockfile != "" {
		lockFile = config.Lockfile{
			File:   cmd.lockfile,
			Remove: lockfile.RemoveLockfile,
		}
	}

	if kraken, err := lockfile.MakeLockFile(lockFile); err != nil {
		return err
	} else {
		defer func() {
			infof("Removing lockfile '%v'", lockFile.File)
			kraken.Release()
		}()
	}

	// ... get config, credentials and version information
	conf := config.NewConfig()
	if err := conf.Load(options.Config); err != nil {
		return fmt.Errorf("could not load configuration (%v)", err)
	}

	credentials, err := getCredentials(cmd.credentials)
	if err != nil {
		return err
	}

	version := getVersionInfo(cmd.workdir, credentials.AccountID)

	// ... get members
	members, err := getMembers(conf, credentials)
	if err != nil {
		return err
	}

	if cmd.debug {
		filename := time.Now().Format("MEMBERS 2006-01-02 15:04:05.tsv")
		path := filepath.Join(os.TempDir(), filename)
		if f, err := os.Create(path); err != nil {
			fmt.Printf("ERROR %v", err)
		} else {
			if cmd.withPIN {
				fmt.Fprintf(f, "%s\n", string(members.AsTableWithPIN().MarshalTextIndent("  ", " ")))
			} else {
				fmt.Fprintf(f, "%s\n", string(members.AsTable().MarshalTextIndent("  ", " ")))
			}

			f.Close()
			fmt.Printf("DEBUG stashed Wild Apricot members list in file %s\n", path)
		}
	}

	// ... get rules
	rules, err := getRules(cmd.rules, cmd.workdir, cmd.debug)
	if err != nil {
		return fmt.Errorf("failed to load ruleset (%v)", err)
	}

	// ... updated?
	// NOTE: Wild Apricot's 'get updated profiles since' query is iffy at best.
	//       So just ignore errors and rely on the hashes for the members and rules
	updated, err := revised(conf, credentials, version.Timestamp)
	if err != nil {
		warnf("Unable to get DB version information (%v)", err)
	}

	// ... make ACL
	doors, err := getDoors(conf)
	if err != nil {
		return err
	}

	if cmd.debug {
		filename := time.Now().Format("DOORS 2006-01-02 15:04:05.txt")
		path := filepath.Join(os.TempDir(), filename)
		if f, err := os.Create(path); err != nil {
			fmt.Printf("ERROR %v", err)
		} else {
			for _, d := range doors {
				fmt.Fprintf(f, "  %v\n", d)
			}
			f.Close()
			fmt.Printf("DEBUG stashed doors list in file %s\n", path)
		}
	}

	makeACL := func(members types.Members, doors []string) (*acl.ACL, error) {
		if cmd.withPIN {
			return rules.MakeACLWithPIN(members, doors)
		} else {
			return rules.MakeACL(members, doors)
		}
	}

	asTable := func(a *acl.ACL) *lib.Table {
		if cmd.withPIN {
			return a.AsTableWithPIN()
		} else {
			return a.AsTable()
		}
	}

	ACL, err := makeACL(*members, doors)
	if err != nil {
		return err
	}

	if cmd.debug {
		filename := time.Now().Format("ACL 2006-01-02 15:04:05.tsv")
		path := filepath.Join(os.TempDir(), filename)
		if f, err := os.Create(path); err != nil {
			fmt.Printf("ERROR %v", err)
		} else {
			fmt.Fprintf(f, "%s\n", string(asTable(ACL).MarshalTextIndent("  ", " ")))
			f.Close()
			fmt.Printf("DEBUG stashed Wild Apricot ACL in file %s\n", path)
		}
	}

	// ... load
	u, devices := getDevices(conf, cmd.debug)
	cards := asTable(ACL)

	// different, err := cmd.compare(&u, devices, cards)
	// if err != nil {
	// 	return err
	// }

	if !cmd.force && !updated && !members.Updated(version.Hashes.Members, cmd.withPIN) && !rules.Updated(version.Hashes.Rules) {
		infof("Nothing to do")
		return nil
	}

	rpt, warnings, err := cmd.load(u, devices, cards)
	if err != nil {
		return err
	}

	if rpt != nil {
		if err := cmd.log(rpt, warnings); err != nil {
			warnf("Error appending summary report to log file (%v)", err)
		}

		if err := cmd.report(rpt, *members); err != nil {
			warnf("Error writing report file (%v)", err)
		}
	}

	if cmd.withPIN {
		membersWithPIN := types.MembersWithPIN{
			Members: *members,
		}

		if err := storeVersionInfo(cmd.workdir, credentials.AccountID, timestamp, &membersWithPIN, rules, ACL); err != nil {
			return fmt.Errorf("failed to store updated version information (%v)", err)
		}
	} else {
		if err := storeVersionInfo(cmd.workdir, credentials.AccountID, timestamp, members, rules, ACL); err != nil {
			return fmt.Errorf("failed to store updated version information (%v)", err)
		}
	}

	return nil
}

// func (cmd *LoadACL) lock() (string, error) {
// 	lockfile := filepath.Join(cmd.workdir, ".wild-apricot", "uhppoted-app-wild-apricot.lock")
// 	pid := fmt.Sprintf("%d\n", os.Getpid())
//
// 	if err := os.MkdirAll(filepath.Dir(lockfile), 0770); err != nil {
// 		return "", fmt.Errorf("Unable to create directory '%v' for lockfile (%v)", lockfile, err)
// 	}
//
// 	if _, err := os.Stat(lockfile); err == nil {
// 		return "", fmt.Errorf("Locked by '%v'", lockfile)
// 	} else if !os.IsNotExist(err) {
// 		return "", fmt.Errorf("Error checking PID lockfile '%v' (%v)", lockfile, err)
// 	}
//
// 	if err := ioutil.WriteFile(lockfile, []byte(pid), 0660); err != nil {
// 		return "", fmt.Errorf("Unable to create lockfile '%v' (%v)", lockfile, err)
// 	}
//
// 	return lockfile, nil
// }

func (cmd *LoadACL) load(u uhppote.IUHPPOTE, devices []uhppote.Device, cards *lib.Table) (map[uint32]lib.Report, []error, error) {
	acl, warnings, err := lib.ParseTable(cards, devices, cmd.strict)
	if err != nil {
		return nil, warnings, err
	}

	for _, w := range warnings {
		warnf("%v", w.Error())
	}

	putACL := func(acl lib.ACL) (map[uint32]lib.Report, []error) {
		if cmd.withPIN {
			return lib.PutACLWithPIN(u, acl, cmd.dryrun)
		} else {
			return lib.PutACL(u, acl, cmd.dryrun)
		}
	}

	rpt, errors := putACL(*acl)
	if len(errors) > 0 {
		return nil, warnings, fmt.Errorf("%v", errors)
	}

	for k, v := range rpt {
		for _, err := range v.Errors {
			errorf("%v  %v", k, err)
		}
	}

	summary := lib.Summarize(rpt)
	format := "%v  unchanged:%v  updated:%v  added:%v  deleted:%v  failed:%v  errors:%v"
	for _, v := range summary {
		infof(format, v.DeviceID, v.Unchanged, v.Updated, v.Added, v.Deleted, v.Failed, v.Errored+len(warnings))
	}

	return rpt, warnings, nil
}

// func (cmd *LoadACL) compare(u uhppote.IUHPPOTE, devices []uhppote.Device, cards *lib.Table) (bool, error) {
// 	current, errors := lib.GetACL(u, devices)
// 	if len(errors) > 0 {
// 		return false, fmt.Errorf("%v", errors)
// 	}
//
// 	acl, _, err := lib.ParseTable(cards, devices, false)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	if acl == nil {
// 		return false, fmt.Errorf("error creating ACL from cards (%v)", cards)
// 	}
//
// 	compare := func(current lib.ACL, acl lib.ACL) (map[uint32]lib.Diff, error) {
// 		if cmd.withPIN {
// 			return lib.CompareWithPIN(current, acl)
// 		} else {
// 			return lib.Compare(current, acl)
// 		}
// 	}
//
// 	diff, err := compare(current, *acl)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	for _, v := range diff {
// 		if v.HasChanges() {
// 			return true, nil
// 		}
// 	}
//
// 	return false, nil
// }

func (cmd *LoadACL) log(rpt map[uint32]lib.Report, warnings []error) error {
	if cmd.logfile != "" {
		var b bytes.Buffer
		summary := lib.Summarize(rpt)
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		format := "%v  %v  unchanged:%v  updated:%v  added:%v  deleted:%v  failed:%v  errors:%v\n"
		if strings.HasSuffix(cmd.logfile, ".tsv") {
			format = "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n"
		}

		for _, v := range summary {
			fmt.Fprintf(&b, format, timestamp, v.DeviceID, v.Unchanged, v.Updated, v.Added, v.Deleted, v.Failed, v.Errored+len(warnings))
		}

		f, err := os.OpenFile(cmd.logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		fmt.Fprintf(f, "%s", b.String())

		return f.Close()
	}

	return nil
}

func (cmd *LoadACL) report(rpt map[uint32]lib.Report, members types.Members) error {
	// ... build card/name map
	names := map[uint32]string{}
	for _, m := range members.Members {
		if m.CardNumber != nil {
			names[uint32(*m.CardNumber)] = m.Name
		}
	}

	// ... build report
	header := []string{"Timestamp", "Action", "Card Number", "Name"}
	index := map[string]int{
		"timestamp":  0,
		"action":     1,
		"cardnumber": 2,
		"name":       3,
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	consolidated := lib.Consolidate(rpt)

	format := []struct {
		Cards  []uint32
		Action string
	}{
		{consolidated.Updated, "Updated"},
		{consolidated.Added, "Added"},
		{consolidated.Deleted, "Deleted"},
		{consolidated.Failed, "Failed"},
		{consolidated.Errored, "Error"},
	}

	rows := [][]string{}
	for _, f := range format {
		for _, card := range f.Cards {
			row := make([]string, len(header))

			for i := 0; i < len(row); i++ {
				row[i] = ""
			}

			if ix, ok := index["timestamp"]; ok {
				row[ix] = timestamp
			}

			if ix, ok := index["action"]; ok {
				row[ix] = f.Action
			}

			if ix, ok := index["cardnumber"]; ok {
				row[ix] = fmt.Sprintf("%v", card)
			}

			if ix, ok := index["name"]; ok {
				row[ix] = fmt.Sprintf("%v", names[card])
			}

			rows = append(rows, row)
		}
	}

	// ... write report
	var b bytes.Buffer
	if strings.HasSuffix(cmd.rptfile, ".tsv") {
		w := csv.NewWriter(&b)
		w.Comma = '\t'

		for _, row := range rows {
			w.Write(row)
		}

		w.Flush()
	} else {
		table := lib.Table{
			Header:  header,
			Records: rows,
		}

		fmt.Fprintf(&b, "%s\n", string(table.MarshalTextIndent("  ", " ")))
	}

	if cmd.rptfile != "" {
		f, err := os.OpenFile(cmd.rptfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		fmt.Fprintf(f, "%s", b.String())

		return f.Close()
	}

	fmt.Printf("%s\n", b.String())
	return nil
}
