package commands

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	lib "github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/config"

	"github.com/uhppoted/uhppoted-app-wild-apricot/acl"
	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

var GetACLCmd = GetACL{
	workdir:     DEFAULT_WORKDIR,
	credentials: filepath.Join(DEFAULT_CONFIG_DIR, ".wild-apricot", "credentials.json"),
	rules:       filepath.Join(DEFAULT_CONFIG_DIR, "wild-apricot.grl"),
	withPIN:     false,
	debug:       false,
}

type GetACL struct {
	workdir     string
	credentials string
	rules       string
	file        string
	withPIN     bool
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
	fmt.Printf("  Usage: %s [--debug] [--config <file>] get-acl [--credentials <file>] [--with-pin] [--rules <url>] [--file <file>]\n", APP)
	fmt.Println()
	fmt.Println("  Downloads an access control list from a Wild Apricot member database, applies the ACL rules and")
	fmt.Println("  stores the generated access control list to a TSV file")
	fmt.Println()

	helpOptions(cmd.FlagSet())

	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    uhppote-app-wild-apricot --debug get-acl --credentials ".credentials/wild-apricot.json" \"`)
	fmt.Println(`                                             --rules "wild-apricot.grl" \`)
	fmt.Println(`                                             --file "ACL.tsv"`)
	fmt.Println()
}

func (cmd *GetACL) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("get-acl", flag.ExitOnError)

	flagset.StringVar(&cmd.workdir, "workdir", cmd.workdir, "Directory for working files (tokens, revisions, etc)'")
	flagset.StringVar(&cmd.credentials, "credentials", cmd.credentials, "Path for the 'credentials.json' file. Defaults to "+cmd.credentials)
	flagset.StringVar(&cmd.rules, "rules", cmd.rules, "URI for the 'grule' rules file. Support file path, HTTP and HTTPS. Defaults to "+cmd.rules)
	flagset.StringVar(&cmd.file, "file", cmd.file, "Output file name. Defaults to stdout")
	flagset.BoolVar(&cmd.withPIN, "with-pin", cmd.withPIN, "Include card keypad PIN code in retrieved ACL information")

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

	credentials, err := getCredentials(cmd.credentials)
	if err != nil {
		return err
	}

	members, err := getMembers(conf, credentials)
	if err != nil {
		return err
	}

	rules, err := getRules(cmd.rules, cmd.workdir, cmd.debug)
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

	makeACL := func(members types.Members, doors Doors) (*acl.ACL, error) {
		if cmd.withPIN {
			return rules.MakeACLWithPIN(members, doors)
		} else {
			return rules.MakeACL(members, doors)
		}
	}

	ACL, err := makeACL(*members, doors)
	if err != nil {
		return err
	}

	asTable := func(a *acl.ACL) *lib.Table {
		if cmd.withPIN {
			return a.AsTableWithPIN()
		} else {
			return a.AsTable()
		}
	}

	_, devices := getDevices(conf, cmd.debug)
	_, warnings, err := lib.ParseTable(asTable(ACL), devices, false)
	if err != nil {
		return err
	}

	if cmd.debug {
		if cmd.withPIN {
			fmt.Printf("ACL:\n%s\n", string(ACL.AsTableWithPIN().MarshalTextIndent("  ", " ")))
		} else {
			fmt.Printf("ACL:\n%s\n", string(ACL.AsTable().MarshalTextIndent("  ", " ")))
		}
	}

	for _, w := range warnings {
		warnf("%v", w.Error())
	}

	// ... write to stdout
	if cmd.file == "" {
		fmt.Fprintln(os.Stdout, string(asTable(ACL).MarshalTextIndent("  ", " ")))
		return nil
	}

	// ... write to TSV file
	asTSV := func(a *acl.ACL, w io.Writer) error {
		if cmd.withPIN {
			return a.ToTSVWithPIN(w)
		} else {
			return a.ToTSV(w)
		}
	}

	var b bytes.Buffer
	if err := asTSV(ACL, &b); err != nil {
		return fmt.Errorf("Error creating TSV file (%v)", err)
	}

	if err := write(cmd.file, b.Bytes()); err != nil {
		return err
	}

	infof("ACL saved to %s", cmd.file)

	return nil
}
