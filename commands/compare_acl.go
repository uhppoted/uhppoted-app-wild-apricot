package commands

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
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

type summary struct {
	header []string
	data   [][]string
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
	flagset.StringVar(&cmd.file, "report", cmd.file, "Report file name. Defaults to stdout")

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

	// ... summary?
	if cmd.summary {
		if err := cmd.summarize(*diff); err != nil {
			return err
		}

		info(fmt.Sprintf("ACL compare report summary saved to %s\n", cmd.file))

		return nil
	}

	// ... write report
	//	var b bytes.Buffer
	//
	//	if err := cmd.report(&b, *diff); err != nil {
	//		return err
	//	}
	//
	//	// ... write to stdout
	//
	//	if cmd.file == "" {
	//		//	text, err := acl.MarshalText()
	//		//	if err != nil {
	//		//		return fmt.Errorf("Error formatting ACL (%v)", err)
	//		//	}
	//
	//		fmt.Fprintln(os.Stdout, string(b.Bytes()))
	//
	//		return nil
	//	}
	//
	//	if err := write(cmd.file, b.Bytes()); err != nil {
	//		return fmt.Errorf("Error writing 'compare' report to file (%v)", err)
	//	}

	info(fmt.Sprintf("Compare report saved to file %s\n", cmd.file))

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

func (cmd *CompareACL) summarize(diff api.SystemDiff) error {
	rpt, err := summarize(diff)
	if err != nil {
		return err
	}

	if cmd.file == "" {
		text, err := rpt.MarshalTextIndent("  ")
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Printf("  ACL Compare Report %s\n", time.Now().Format("2006-01-02 15:03:04"))
		fmt.Println()
		fmt.Printf("%v\n", string(text))
		fmt.Println()
		return nil
	}

	var b bytes.Buffer
	if err := rpt.toTSV(&b); err != nil {
		return fmt.Errorf("Error creating TSV file from 'compare' report (%v)", err)
	}

	if err := write(cmd.file, b.Bytes()); err != nil {
		return err
	}

	return nil
}

func summarize(diff api.SystemDiff) (*summary, error) {
	keys := []uint32{}
	for k, _ := range diff {
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

	return &summary{header, data}, nil
}

func (rpt *summary) MarshalText() ([]byte, error) {
	return rpt.MarshalTextIndent("")
}

func (rpt *summary) MarshalTextIndent(indent string) ([]byte, error) {
	table := [][]string{}

	table = append(table, rpt.header)
	table = append(table, rpt.data...)

	var b bytes.Buffer

	if len(table) > 0 {
		widths := make([]int, len(table[0]))
		for _, row := range table {
			for i, field := range row {
				if len(field) > widths[i] {
					widths[i] = len(field)
				}
			}
		}

		for i := 1; i < len(widths); i++ {
			widths[i-1] += 2
		}

		// Print header
		for _, row := range table[0:1] {
			fmt.Fprintf(&b, "%s", indent)
			for i, field := range row {
				fmt.Fprintf(&b, "%-*v", widths[i], field)
			}
			fmt.Fprintln(&b)
		}

		for _, row := range table[0:1] {
			fmt.Fprintf(&b, "%s", indent)
			for i, field := range row {
				fmt.Fprintf(&b, "%-*v", widths[i], strings.Repeat("-", len(field)))
			}
			fmt.Fprintln(&b)
		}

		// Print data
		previous := ""
		for _, row := range table[1:] {
			if row[0] != previous {
				if previous != "" {
					fmt.Fprintln(&b)
				}
				previous = row[0]

				fmt.Fprintf(&b, "%s", indent)
				for i, field := range row {
					fmt.Fprintf(&b, "%-*v", widths[i], field)
				}
			} else {
				fmt.Fprintf(&b, "%s", indent)
				fmt.Fprintf(&b, "%-*v", widths[0], "")
				for i, field := range row[1:] {
					fmt.Fprintf(&b, "%-*v", widths[i+1], field)
				}
			}

			fmt.Fprintln(&b)
		}
	}

	return b.Bytes(), nil
}

func (rpt *summary) toTSV(f io.Writer) error {
	w := csv.NewWriter(f)
	w.Comma = '\t'

	w.Write(rpt.header)
	for _, row := range rpt.data {
		w.Write(row)
	}

	w.Flush()

	return nil
}
