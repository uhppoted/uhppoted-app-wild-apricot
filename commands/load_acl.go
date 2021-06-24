package commands

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
	api "github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/config"
)

var LoadACLCmd = LoadACL{
	workdir:     DEFAULT_WORKDIR,
	credentials: filepath.Join(DEFAULT_CONFIG_DIR, ".wild-apricot", "credentials.json"),
	rules:       filepath.Join(DEFAULT_CONFIG_DIR, "wild-apricot.grl"),
	force:       false,
	strict:      false,
	dryrun:      false,
	debug:       false,
}

type LoadACL struct {
	workdir     string
	credentials string
	rules       string
	force       bool
	strict      bool
	dryrun      bool
	logfile     string
	rptfile     string
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

	flagset.StringVar(&cmd.workdir, "workdir", cmd.workdir, "Directory for working files (tokens, revisions, etc)'")
	flagset.StringVar(&cmd.credentials, "credentials", cmd.credentials, "Path for the 'credentials.json' file. Defaults to "+cmd.credentials)
	flagset.StringVar(&cmd.rules, "rules", cmd.rules, "URI for the 'grule' rules file. Support file path, HTTP and HTTPS. Defaults to "+cmd.rules)
	flagset.BoolVar(&cmd.force, "force", cmd.force, "Forces an update, overriding the  version and compare logic")
	flagset.BoolVar(&cmd.strict, "strict", cmd.strict, "Fails with an error if the members list contains duplicate card numbers")
	flagset.BoolVar(&cmd.dryrun, "dry-run", cmd.dryrun, "Simulates a load-acl without making any changes to the access controllers")
	flagset.StringVar(&cmd.logfile, "log", cmd.logfile, "File to which the (optional) summary report is appended")
	flagset.StringVar(&cmd.rptfile, "report", cmd.rptfile, "File to which the detail report is written. Defaults to stdout if not provided")

	return flagset
}

func (cmd *LoadACL) Execute(args ...interface{}) error {
	timestamp := time.Now()
	options := args[0].(*Options)

	cmd.debug = options.Debug

	// ... check parameters
	if strings.TrimSpace(cmd.credentials) == "" {
		return fmt.Errorf("Invalid credentials file")
	}

	if strings.TrimSpace(cmd.rules) == "" {
		return fmt.Errorf("Invalid rules file")
	}

	// ... locked?
	lockfile, err := cmd.lock()
	if err != nil {
		return err
	}

	defer func() {
		info(fmt.Sprintf("Removing lockfile '%v'", lockfile))
		os.Remove(lockfile)
	}()

	// ... get config, credentials and version information
	conf := config.NewConfig()
	if err := conf.Load(options.Config); err != nil {
		return fmt.Errorf("Could not load configuration (%v)", err)
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
		fmt.Printf("MEMBERS:\n%s\n", string(members.AsTable().MarshalTextIndent("  ", " ")))
	}

	// ... get rules

	rules, err := getRules(cmd.rules, cmd.debug)
	if err != nil {
		return fmt.Errorf("Unable to create ruleset (%v)", err)
	}

	// ... updated?

	// NOTE: Wild Apricot's 'get updated profiles since' query is iffy at best.
	//       So just ignore errors and rely on the hashes for the members and rules
	updated, err := revised(conf, credentials, version.Timestamp)
	if err != nil {
		warn(fmt.Sprintf("Unable to get DB version information (%v)", err))
	}

	// ... make ACL

	doors, err := getDoors(conf)
	if err != nil {
		return err
	}

	if cmd.debug {
		fmt.Printf("DOORS:\n")
		for _, d := range doors {
			fmt.Printf("  %v\n", d)
		}
		fmt.Println()
	}

	acl, err := rules.MakeACL(*members, doors)
	if err != nil {
		return err
	}

	if cmd.debug {
		fmt.Printf("ACL:\n%s\n", string(acl.AsTable().MarshalTextIndent("  ", " ")))
	}

	// ... load
	u, devices := getDevices(conf, cmd.debug)
	cards := acl.AsTable()

	// different, err := cmd.compare(&u, devices, cards)
	// if err != nil {
	// 	return err
	// }

	if !cmd.force && !updated && !members.Updated(version.Hashes.Members) && !rules.Updated(version.Hashes.Rules) {
		info("Nothing to do")
		return nil
	}

	rpt, warnings, err := cmd.load(u, devices, cards)
	if err != nil {
		return err
	}

	if rpt != nil {
		if err := cmd.log(rpt, warnings); err != nil {
			warn(fmt.Sprintf("Error appending summary report to log file (%v)", err))
		}

		if err := cmd.report(rpt, *members); err != nil {
			warn(fmt.Sprintf("Error writing report file (%v)", err))
		}
	}

	if err := storeVersionInfo(cmd.workdir, credentials.AccountID, timestamp, members, rules, acl); err != nil {
		return fmt.Errorf("Failed to store updated version information (%v)", err)
	}

	return nil
}

func (cmd *LoadACL) lock() (string, error) {
	lockfile := filepath.Join(cmd.workdir, ".wild-apricot", "uhppoted-app-wild-apricot.lock")
	pid := fmt.Sprintf("%d\n", os.Getpid())

	if err := os.MkdirAll(filepath.Dir(lockfile), 0770); err != nil {
		return "", fmt.Errorf("Unable to create directory '%v' for lockfile (%v)", lockfile, err)
	}

	if _, err := os.Stat(lockfile); err == nil {
		return "", fmt.Errorf("Locked by '%v'", lockfile)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("Error checking PID lockfile '%v' (%v)", lockfile, err)
	}

	if err := ioutil.WriteFile(lockfile, []byte(pid), 0660); err != nil {
		return "", fmt.Errorf("Unable to create lockfile '%v' (%v)", lockfile, err)
	}

	return lockfile, nil
}

func (cmd *LoadACL) load(u uhppote.IUHPPOTE, devices []uhppote.Device, cards *api.Table) (map[uint32]api.Report, []error, error) {
	acl, warnings, err := api.ParseTable(cards, devices, cmd.strict)
	if err != nil {
		return nil, warnings, err
	}

	for _, w := range warnings {
		warn(w.Error())
	}

	rpt, errors := api.PutACL(u, *acl, cmd.dryrun)
	if len(errors) > 0 {
		return nil, warnings, fmt.Errorf("%v", errors)
	}

	for k, v := range rpt {
		for _, err := range v.Errors {
			fatal(fmt.Sprintf("%v  %v", k, err))
		}
	}

	summary := api.Summarize(rpt)
	format := "%v  unchanged:%v  updated:%v  added:%v  deleted:%v  failed:%v  errors:%v"
	for _, v := range summary {
		info(fmt.Sprintf(format, v.DeviceID, v.Unchanged, v.Updated, v.Added, v.Deleted, v.Failed, v.Errored+len(warnings)))
	}

	return rpt, warnings, nil
}

func (cmd *LoadACL) compare(u uhppote.IUHPPOTE, devices []uhppote.Device, cards *api.Table) (bool, error) {
	current, errors := api.GetACL(u, devices)
	if len(errors) > 0 {
		return false, fmt.Errorf("%v", errors)
	}

	acl, _, err := api.ParseTable(cards, devices, false)
	if err != nil {
		return false, err
	}

	if acl == nil {
		return false, fmt.Errorf("Error creating ACL from cards (%v)", cards)
	}

	diff, err := api.Compare(current, *acl)
	if err != nil {
		return false, err
	}

	for _, v := range diff {
		if v.HasChanges() {
			return true, nil
		}
	}

	return false, nil
}

func (cmd *LoadACL) log(rpt map[uint32]api.Report, warnings []error) error {
	if cmd.logfile != "" {
		var b bytes.Buffer
		summary := api.Summarize(rpt)
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

		fmt.Fprintf(f, "%s", string(b.Bytes()))

		return f.Close()
	}

	return nil
}

func (cmd *LoadACL) report(rpt map[uint32]api.Report, members types.Members) error {
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

	consolidated := api.Consolidate(rpt)

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
		table := api.Table{
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

		fmt.Fprintf(f, "%s", string(b.Bytes()))

		return f.Close()
	}

	fmt.Printf("%s\n", string(b.Bytes()))
	return nil
}
