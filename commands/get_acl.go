package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/uhppoted/uhppoted-app-wild-apricot/acl"
	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

var GetACLCmd = GetACL{
	workdir:     DEFAULT_WORKDIR,
	credentials: filepath.Join(DEFAULT_CONFIG_DIR, ".wild-apricot", "credentials.json"),
	rules:       filepath.Join(DEFAULT_CONFIG_DIR, "wild-apricot.grl"),
	file:        time.Now().Format("2006-01-02T150405.tsv"),
	debug:       false,
}

type GetACL struct {
	workdir     string
	credentials string
	rules       string
	file        string
	debug       bool
}

func (cmd *GetACL) Name() string {
	return "get-acl"
}

func (cmd *GetACL) Description() string {
	return "Retrieves an access control list from a Wild Apricot member database and stores it to a file"
}

func (cmd *GetACL) Usage() string {
	return "--credentials <file> --rules <url> --file <file>"
}

func (cmd *GetACL) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s [--debug] [--config <file>] get-acl [--credentials <file>] [--rules <url>] [--file <file>]\n", APP)
	fmt.Println()
	fmt.Println("  Downloads an access control list from a Wild Apricot member database, applies the ACL rules and")
	fmt.Println("  stores the generated access control list to a TSV file")
	fmt.Println()

	helpOptions(cmd.FlagSet())

	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    uhppote-app-wild-apricot --debug get-acl --credentials ".credentials/wild-apricot.json" \"`)
	fmt.Println(`                                             --rules "wild-apricot.grl" \`)
	fmt.Println(`                                             --file "example.tsv"`)
	fmt.Println()
}

func (cmd *GetACL) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("get-acl", flag.ExitOnError)

	flagset.StringVar(&cmd.workdir, "workdir", cmd.workdir, "Directory for working files (tokens, revisions, etc)'")
	//	flagset.StringVar(&cmd.config, "config", cmd.config, "Path for the configuration file. Defaults to "+cmd.config)
	flagset.StringVar(&cmd.credentials, "credentials", cmd.credentials, "Path for the 'credentials.json' file. Defaults to "+cmd.credentials)
	flagset.StringVar(&cmd.rules, "rules", cmd.rules, "URI for the 'grule' rules file. Support file path, HTTP and HTTPS. Defaults to "+cmd.rules)
	flagset.StringVar(&cmd.file, "file", cmd.file, "TSV file name. Defaults to 'ACL - <yyyy-mm-dd HHmmss>.tsv'")

	return flagset
}

func (cmd *GetACL) Execute(args ...interface{}) error {
	options := args[0].(*Options)

	cmd.debug = options.Debug

	// ... check parameters
	if strings.TrimSpace(cmd.credentials) == "" {
		return fmt.Errorf("Invalid credentials file")
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

	// ... create ACL

	doors, err := getDoors(options.Config)
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

	// ... save to TSV file
	tmp, err := ioutil.TempFile(os.TempDir(), "ACL")
	if err != nil {
		return err
	}

	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	if err := acl.ToTSV(tmp); err != nil {
		return fmt.Errorf("Error creating TSV file (%v)", err)
	}

	tmp.Close()

	dir := filepath.Dir(cmd.file)
	if err := os.MkdirAll(dir, 0770); err != nil {
		return err
	}

	if err := os.Rename(tmp.Name(), cmd.file); err != nil {
		return err
	}

	info(fmt.Sprintf("Retrieved ACL to file %s\n", cmd.file))

	return nil
}
