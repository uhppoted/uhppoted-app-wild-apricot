package commands

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/uhppoted/uhppoted-api/config"
)

var GetACLCmd = GetACL{
	workdir:     DEFAULT_WORKDIR,
	credentials: filepath.Join(DEFAULT_CONFIG_DIR, ".wild-apricot", "credentials.json"),
	rules:       filepath.Join(DEFAULT_CONFIG_DIR, "wild-apricot.grl"),
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
	flagset.StringVar(&cmd.credentials, "credentials", cmd.credentials, "Path for the 'credentials.json' file. Defaults to "+cmd.credentials)
	flagset.StringVar(&cmd.rules, "rules", cmd.rules, "URI for the 'grule' rules file. Support file path, HTTP and HTTPS. Defaults to "+cmd.rules)
	flagset.StringVar(&cmd.file, "file", cmd.file, "Output file name. Defaults to stdout")

	return flagset
}

func (cmd *GetACL) Execute(args ...interface{}) error {
	options := args[0].(*Options)

	cmd.debug = options.Debug

	// ... check parameters
	if strings.TrimSpace(cmd.credentials) == "" {
		return fmt.Errorf("Invalid credentials file")
	}

	if strings.TrimSpace(cmd.rules) == "" {
		return fmt.Errorf("Invalid rules file")
	}

	// ... get config, members and rules
	conf := config.NewConfig()
	if err := conf.Load(options.Config); err != nil {
		return fmt.Errorf("Could not load configuration (%v)", err)
	}

	members, err := getMembers(cmd.credentials)
	if err != nil {
		return err
	}

	rules, err := getRules(cmd.rules, cmd.debug)
	if err != nil {
		return err
	}

	if cmd.debug {
		if text, err := members.MarshalTextIndent("  "); err == nil {
			fmt.Printf("MEMBERS:\n%s\n", string(text))
		}
	}

	// ... create ACL
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

	// ... write to stdout

	if cmd.file == "" {
		text, err := acl.MarshalText()
		if err != nil {
			return fmt.Errorf("Error formatting ACL (%v)", err)
		}

		fmt.Fprintln(os.Stdout, string(text))

		return nil
	}

	// ... write to TSV file

	var b bytes.Buffer
	if err := acl.ToTSV(&b); err != nil {
		return fmt.Errorf("Error creating TSV file (%v)", err)
	}

	if err := write(cmd.file, b.Bytes()); err != nil {
		return err
	}

	info(fmt.Sprintf("ACL saved to %s\n", cmd.file))

	return nil
}
