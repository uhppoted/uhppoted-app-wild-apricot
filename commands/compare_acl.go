package commands

import (
	"bytes"
	"flag"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/uhppoted/uhppote-core/uhppote"
	lib "github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/config"
	"github.com/uhppoted/uhppoted-lib/lockfile"

	"github.com/uhppoted/uhppoted-app-wild-apricot/acl"
	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

var CompareACLCmd = CompareACL{
	workdir:     DEFAULT_WORKDIR,
	credentials: filepath.Join(DEFAULT_CONFIG_DIR, ".wild-apricot", "credentials.json"),
	rules:       filepath.Join(DEFAULT_CONFIG_DIR, "wild-apricot.grl"),
	withPIN:     false,
	summary:     false,
	strict:      false,
	lockfile:    "",
	debug:       false,
}

type CompareACL struct {
	workdir     string
	credentials string
	rules       string
	file        string
	withPIN     bool
	summary     bool
	strict      bool
	lockfile    string
	debug       bool
}

func (cmd *CompareACL) Name() string {
	return "compare-acl"
}

func (cmd *CompareACL) Description() string {
	return "Retrieves an access control list from a Wild Apricot member database and compares card lists on the configured controllers"
}

func (cmd *CompareACL) Usage() string {
	return "--credentials <file> --rules <url> --report <file>"
}

func (cmd *CompareACL) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s [--debug] [--config <file>] compare-acl [--credentials <file>] [--rules <url>] [--with-pin] [--summary] [--report <file>]\n", APP)
	fmt.Println()
	fmt.Println("  Downloads an access control list from a Wild Apricot member database, applies the ACL rules and stores the generated")
	fmt.Println("  access control list to a TSV file")
	fmt.Println()

	helpOptions(cmd.FlagSet())

	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    uhppote-app-wild-apricot --debug --config uhppoted.conf compare-acl --credentials ".credentials/wild-apricot.json" \"`)
	fmt.Println(`                                                                         --rules "wild-apricot.grl" \`)
	fmt.Println(`                                                                         --report "example.tsv"`)
	fmt.Println()
}

func (cmd *CompareACL) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("compare-acl", flag.ExitOnError)
	lockfile := filepath.Join(cmd.workdir, ".wild-apricot", "uhppoted-app-wild-apricot.lock")

	flagset.StringVar(&cmd.workdir, "workdir", cmd.workdir, "Directory for working files (tokens, revisions, etc)'")
	flagset.StringVar(&cmd.credentials, "credentials", cmd.credentials, "Path for the 'credentials.json' file. Defaults to "+cmd.credentials)
	flagset.StringVar(&cmd.rules, "rules", cmd.rules, "URI for the 'grule' rules file. Support file path, HTTP and HTTPS. Defaults to "+cmd.rules)
	flagset.BoolVar(&cmd.withPIN, "with-pin", cmd.withPIN, "Include card keypad PIN code ACL comparison")
	flagset.BoolVar(&cmd.summary, "summary", cmd.summary, "Report only a summary of the comparison. Defaults to "+fmt.Sprintf("%v", cmd.summary))
	flagset.StringVar(&cmd.file, "report", cmd.file, "Report file name. Defaults to stdout")
	flagset.BoolVar(&cmd.strict, "strict", cmd.strict, "Fails with an error if the members list contains duplicate card numbers")
	flagset.StringVar(&cmd.lockfile, "lockfile", cmd.lockfile, fmt.Sprintf("Filepath for lock file. Defaults to %v", lockfile))

	return flagset
}

func (cmd *CompareACL) Execute(args ...interface{}) error {
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

	// ... get config, members and rules
	conf := config.NewConfig()
	if err := conf.Load(options.Config); err != nil {
		return fmt.Errorf("could not load configuration (%v)", err)
	}

	credentials, err := getCredentials(cmd.credentials)
	if err != nil {
		return err
	}

	rules, err := getRules(cmd.rules, cmd.workdir, cmd.debug)
	if err != nil {
		return err
	}

	members, err := getMembers(conf, credentials)
	if err != nil {
		return err
	}

	if cmd.debug {
		if cmd.withPIN {
			fmt.Printf("MEMBERS:\n%s\n", string(members.AsTableWithPIN().MarshalTextIndent("  ", " ")))
		} else {
			fmt.Printf("MEMBERS:\n%s\n", string(members.AsTable().MarshalTextIndent("  ", " ")))
		}
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

	makeACL := func(members types.Members, doors []string) (*acl.ACL, error) {
		if cmd.withPIN {
			return rules.MakeACLWithPIN(members, doors)
		} else {
			return rules.MakeACL(members, doors)
		}
	}

	acl, err := makeACL(*members, doors)
	if err != nil {
		return err
	}

	if cmd.debug {
		if cmd.withPIN {
			fmt.Printf("ACL:\n%s\n", string(acl.AsTableWithPIN().MarshalTextIndent("  ", " ")))
		} else {
			fmt.Printf("ACL:\n%s\n", string(acl.AsTable().MarshalTextIndent("  ", " ")))
		}
	}

	// ... compare

	u, devices := getDevices(conf, cmd.debug)

	diff, err := cmd.compare(u, devices, acl)
	if err != nil {
		return err
	}

	// ... summary?
	if cmd.summary {
		return cmd.summarize(*diff)
	}

	// ... detail report
	return cmd.report(*members, *diff)
}

func (cmd *CompareACL) compare(u uhppote.IUHPPOTE, devices []uhppote.Device, cards *acl.ACL) (*lib.SystemDiff, error) {
	current, errors := lib.GetACL(u, devices)
	if len(errors) > 0 {
		return nil, fmt.Errorf("%v", errors)
	}

	asTable := func(cards *acl.ACL) *lib.Table {
		if cmd.withPIN {
			return cards.AsTableWithPIN()
		} else {
			return cards.AsTable()
		}
	}

	acl, warnings, err := lib.ParseTable(asTable(cards), devices, cmd.strict)
	if err != nil {
		return nil, err
	}

	for _, w := range warnings {
		warnf("%v", w.Error())
	}

	if acl == nil {
		return nil, fmt.Errorf("error creating ACL from cards (%v)", cards)
	}

	compare := func(current, acl lib.ACL) (map[uint32]lib.Diff, error) {
		if cmd.withPIN {
			return lib.CompareWithPIN(current, acl)
		} else {
			return lib.Compare(current, acl)
		}
	}

	d, err := compare(current, *acl)
	if err != nil {
		return nil, err
	}

	diff := lib.SystemDiff(d)

	return &diff, nil
}

func (cmd *CompareACL) summarize(diff lib.SystemDiff) error {
	rpt := summarize(diff)

	if cmd.file == "" {
		fmt.Println()
		fmt.Printf("  ACL Compare Report %s\n", time.Now().Format("2006-01-02 15:03:04"))
		fmt.Println()
		fmt.Printf("%v\n", string(rpt.MarshalTextIndent("  ", " ")))
		fmt.Println()

		return nil
	}

	var b bytes.Buffer
	if err := rpt.ToTSV(&b); err != nil {
		return fmt.Errorf("error creating TSV file from 'compare' report (%v)", err)
	}

	if err := write(cmd.file, b.Bytes()); err != nil {
		return err
	}

	infof("ACL compare report summary saved to %s", cmd.file)

	return nil
}

func (cmd *CompareACL) report(members types.Members, diff lib.SystemDiff) error {
	rpt := detail(members, diff)

	if cmd.file == "" {
		if !diff.HasChanges() {
			fmt.Println()
			fmt.Printf("  ACL Compare Report %s\n", time.Now().Format("2006-01-02 15:03:04"))
			fmt.Println()
			fmt.Printf("%v\n", "  NO DIFFERENCES")
			fmt.Println()
		} else {
			fmt.Println()
			fmt.Printf("  ACL Compare Report %s\n", time.Now().Format("2006-01-02 15:03:04"))
			fmt.Println()
			fmt.Printf("%v\n", string(rpt.MarshalTextIndent("  ", " ")))
			fmt.Println()
		}
		return nil
	}

	var b bytes.Buffer
	if err := rpt.ToTSV(&b); err != nil {
		return fmt.Errorf("error creating TSV file from 'compare' report (%v)", err)
	}

	if err := write(cmd.file, b.Bytes()); err != nil {
		return err
	}

	infof("ACL compare report saved to %s", cmd.file)

	return nil
}

func summarize(diff lib.SystemDiff) *lib.Table {
	keys := []uint32{}
	for k := range diff {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	header := []string{"Controller", "Incorrect", "Missing", "Unexpected"}
	data := [][]string{}
	for _, k := range keys {
		if v, ok := diff[k]; ok {
			updated := []uint32{}
			added := []uint32{}
			deleted := []uint32{}

			for _, c := range v.Updated {
				updated = append(updated, c.CardNumber)
			}

			for _, c := range v.Added {
				added = append(added, c.CardNumber)
			}

			for _, c := range v.Deleted {
				deleted = append(deleted, c.CardNumber)
			}

			sort.Slice(updated, func(i, j int) bool { return updated[i] < updated[j] })
			sort.Slice(added, func(i, j int) bool { return added[i] < added[j] })
			sort.Slice(deleted, func(i, j int) bool { return deleted[i] < deleted[j] })

			N := len(updated)
			if len(added) > N {
				N = len(added)
			}
			if len(deleted) > N {
				N = len(deleted)
			}

			for i := 0; i < N; i++ {
				row := []string{
					fmt.Sprintf("%v", k),
					"",
					"",
					"",
				}

				if i < len(updated) {
					row[1] = fmt.Sprintf("%v", updated[i])
				}

				if i < len(added) {
					row[2] = fmt.Sprintf("%v", added[i])
				}

				if i < len(deleted) {
					row[3] = fmt.Sprintf("%v", deleted[i])
				}

				data = append(data, row)
			}
		}
	}

	table := lib.Table{
		Header:  header,
		Records: data,
	}

	return &table
}

func detail(members types.Members, diff lib.SystemDiff) *lib.Table {
	type card struct {
		cardnumber uint32
		action     string
	}

	cards := map[uint32]card{}

	for _, v := range diff {
		for _, c := range v.Updated {
			cards[c.CardNumber] = card{
				cardnumber: c.CardNumber,
				action:     "update",
			}
		}
	}

	for _, v := range diff {
		for _, c := range v.Added {
			if _, ok := cards[c.CardNumber]; !ok {
				cards[c.CardNumber] = card{
					cardnumber: c.CardNumber,
					action:     "add",
				}
			}
		}
	}

	for _, v := range diff {
		for _, c := range v.Deleted {
			cards[c.CardNumber] = card{
				cardnumber: c.CardNumber,
				action:     "delete",
			}
		}
	}

	names := map[uint32]string{}
	for _, v := range members.Members {
		if v.CardNumber != nil {
			names[uint32(*v.CardNumber)] = clean(v.Name)
		}
	}

	keys := []uint32{}
	for k := range cards {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	timestamp := time.Now().Format("2006-01-02 15:03:04")
	header := []string{"Timestamp", "Name", "Card Number", "Action"}
	data := [][]string{}
	for _, k := range keys {
		if card, ok := cards[k]; ok {
			data = append(data, []string{
				timestamp,
				names[card.cardnumber],
				fmt.Sprintf("%v", card.cardnumber),
				card.action,
			})
		}
	}

	table := lib.Table{
		Header:  header,
		Records: data,
	}

	return &table
}
