package commands

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/uhppoted/uhppote-core/device"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/config"
	"github.com/uhppoted/uhppoted-app-wild-apricot/acl"
	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

var CompareACLCmd = CompareACL{
	workdir:     DEFAULT_WORKDIR,
	credentials: filepath.Join(DEFAULT_CONFIG_DIR, ".wild-apricot", "credentials.json"),
	rules:       filepath.Join(DEFAULT_CONFIG_DIR, "wild-apricot.grl"),
	file:        time.Now().Format("ACL 2006-01-02T150405.rpt"),
	summary:     false,
	debug:       false,
}

type CompareACL struct {
	workdir     string
	credentials string
	rules       string
	file        string
	summary     bool
	debug       bool
}

func (cmd *CompareACL) Name() string {
	return "compare-acl"
}

func (cmd *CompareACL) Description() string {
	return "Retrieves an access control list from a Wild Apricot member database and compares it to the current controllers card lists"
}

func (cmd *CompareACL) Usage() string {
	return "--credentials <file> --rules <url> --report <file>"
}

func (cmd *CompareACL) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s [--debug] [--config <file>] compare-acl [--credentials <file>] [--rules <url>] [--summary] [--file <file>]\n", APP)
	fmt.Println()
	fmt.Println("  Downloads an access control list from a Wild Apricot member database, applies the ACL rules and stores the generated")
	fmt.Println("  access control list to a TSV file")
	fmt.Println()

	helpOptions(cmd.FlagSet())

	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    uhppote-app-wild-apricot --debug --config uhppoted.conf compare-acl --credentials ".credentials/wild-apricot.json" \"`)
	fmt.Println(`                                                                         --rules "wild-apricot.grl" \`)
	fmt.Println(`                                                                         --file "example.tsv"`)
	fmt.Println()
}

func (cmd *CompareACL) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("compare-acl", flag.ExitOnError)

	flagset.StringVar(&cmd.workdir, "workdir", cmd.workdir, "Directory for working files (tokens, revisions, etc)'")
	flagset.StringVar(&cmd.credentials, "credentials", cmd.credentials, "Path for the 'credentials.json' file. Defaults to "+cmd.credentials)
	flagset.StringVar(&cmd.rules, "rules", cmd.rules, "URI for the 'grule' rules file. Support file path, HTTP and HTTPS. Defaults to "+cmd.rules)
	flagset.BoolVar(&cmd.summary, "summary", cmd.summary, "Report only a summary of the comparison. Defaults to "+fmt.Sprintf("%v", cmd.summary))
	flagset.StringVar(&cmd.file, "report", cmd.file, "Report file name. Defaults to 'ACL - <yyyy-mm-dd HHmmss>.rpt'")

	return flagset
}

func (cmd *CompareACL) Execute(args ...interface{}) error {
	options := args[0].(*Options)

	cmd.debug = options.Debug

	// ... check parameters
	if strings.TrimSpace(cmd.credentials) == "" {
		return fmt.Errorf("Invalid credentials file")
	}

	if strings.TrimSpace(cmd.rules) == "" {
		return fmt.Errorf("Invalid rules file")
	}

	if strings.TrimSpace(cmd.file) == "" {
		return fmt.Errorf("Invalid output file")
	}

	// ... load config
	conf := config.NewConfig()
	if err := conf.Load(options.Config); err != nil {
		return fmt.Errorf("Could not load configuration (%v)", err)
	}

	// ... load rules

	ruleset, err := fetch(cmd.rules)
	if err != nil {
		return err
	}

	rules, err := acl.NewRules(ruleset, cmd.debug)
	if err != nil {
		return err
	}

	// ... get contacts list and member groups

	credentials, err := getCredentials(cmd.credentials)
	if err != nil {
		return err
	}

	token, err := wildapricot.Authorize(credentials.APIKey)
	if err != nil {
		return err
	}

	contacts, err := wildapricot.GetContacts(credentials.Account, token)
	if err != nil {
		return err
	}

	groups, err := wildapricot.GetMemberGroups(credentials.Account, token)
	if err != nil {
		return err
	}

	members, err := types.MakeMemberList(contacts, groups)
	if err != nil {
		return err
	} else if members == nil {
		return fmt.Errorf("Invalid members list")
	}

	if cmd.debug {
		if text, err := members.MarshalTextIndent("  "); err == nil {
			fmt.Printf("MEMBERS:\n%s\n", string(text))
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

	acl, err := rules.MakeACL(*members, doors)
	if err != nil {
		return err
	}

	if cmd.debug {
		if text, err := acl.MarshalTextIndent("  "); err == nil {
			fmt.Printf("ACL:\n%s\n", string(text))
		}
	}

	// ... compare

	u, devices := getDevices(conf, cmd.debug)

	diff, err := compare(&u, devices, acl)
	if err != nil {
		return err
	}

	// ... summary
	if !cmd.summary {
		cmd.summarize(os.Stdout, *diff)
	}

	// ... save to report file
	tmp, err := ioutil.TempFile(os.TempDir(), "RPT")
	if err != nil {
		return err
	}

	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	keys := []uint32{}
	for k, _ := range *diff {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, k := range keys {
		if v, ok := (*diff)[k]; ok {
			fmt.Fprintf(tmp, "%v\n", k)

			if len(v.Updated) > 0 {
				for _, c := range v.Updated {
					fmt.Fprintf(tmp, "  UPDATED %-10v\n", c.CardNumber)
				}
			}

			if len(v.Added) > 0 {
				for _, c := range v.Added {
					fmt.Fprintf(tmp, "  ADDED   %-10v\n", c.CardNumber)
				}
			}

			if len(v.Deleted) > 0 {
				for _, c := range v.Deleted {
					fmt.Fprintf(tmp, "  DELETED %-10v\n", c.CardNumber)
				}
			}

			fmt.Fprintln(tmp)
		}
	}

	tmp.Close()

	dir := filepath.Dir(cmd.file)
	if err := os.MkdirAll(dir, 0770); err != nil {
		return err
	}

	if err := os.Rename(tmp.Name(), cmd.file); err != nil {
		return err
	}

	info(fmt.Sprintf("Compare report saved to file %s\n", cmd.report))

	return nil
}

func compare(u device.IDevice, devices []*uhppote.Device, cards *acl.ACL) (*api.SystemDiff, error) {
	current, err := api.GetACL(u, devices)
	if err != nil {
		return nil, err
	}

	table := cards.AsTable()

	acl, warnings, err := api.ParseTable(&table, devices, false)
	if err != nil {
		return nil, err
	}

	for _, w := range warnings {
		warn(w.Error())
	}

	if acl == nil {
		return nil, fmt.Errorf("Error creating ACL from cards (%v)", cards)
	}

	d, err := api.Compare(current, *acl)
	if err != nil {
		return nil, err
	}

	diff := api.SystemDiff(d)

	return &diff, nil
}

func (cmd *CompareACL) summarize(f io.Writer, diff api.SystemDiff) {
	keys := []uint32{}
	for k, _ := range diff {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, k := range keys {
		if v, ok := diff[k]; ok {
			fmt.Fprintf(f, "%v  updated:%-5v added:%-5v deleted:%-5v\n", k, len(v.Updated), len(v.Added), len(v.Deleted))
		}
	}
}

func (cmd *CompareACL) report(f io.Writer, diff api.SystemDiff) {
	keys := []uint32{}
	for k, _ := range diff {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, k := range keys {
		if v, ok := diff[k]; ok {
			fmt.Fprintf(f, "%v\n", k)

			if len(v.Updated) > 0 {
				for _, c := range v.Updated {
					fmt.Fprintf(f, "  UPDATED %-10v\n", c.CardNumber)
				}
			}

			if len(v.Added) > 0 {
				for _, c := range v.Added {
					fmt.Fprintf(f, "  ADDED   %-10v\n", c.CardNumber)
				}
			}

			if len(v.Deleted) > 0 {
				for _, c := range v.Deleted {
					fmt.Fprintf(f, "  DELETED %-10v\n", c.CardNumber)
				}
			}

			fmt.Fprintln(f)
		}
	}
}
