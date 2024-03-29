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

	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

var GetMembersCmd = GetMembers{
	workdir:     DEFAULT_WORKDIR,
	credentials: filepath.Join(DEFAULT_CONFIG_DIR, ".wild-apricot", "credentials.json"),
	withPIN:     false,
	debug:       false,
}

type GetMembers struct {
	workdir     string
	credentials string
	file        string
	withPIN     bool
	debug       bool
}

func (cmd *GetMembers) Name() string {
	return "get-members"
}

func (cmd *GetMembers) Description() string {
	return "Retrieves a tabular member list from a Wild Apricot member database and stores it to a file"
}

func (cmd *GetMembers) Usage() string {
	return "--credentials <file> --file <file>"
}

func (cmd *GetMembers) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s [--debug] [--config <file>] get-members [--credentials <file>] [--with-pin] [--file <file>]\n", APP)
	fmt.Println()
	fmt.Println("  Downloads the members list from a Wild Apricot member database and (optionally) stores it to a TSV file")
	fmt.Println()

	helpOptions(cmd.FlagSet())

	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    uhppote-app-wild-apricot --debug get-members --credentials ".credentials/wild-apricot.json" \"`)
	fmt.Println(`                                                 --file "members.tsv"`)
	fmt.Println()
}

func (cmd *GetMembers) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("get-members", flag.ExitOnError)

	flagset.StringVar(&cmd.workdir, "workdir", cmd.workdir, "Directory for working files (tokens, revisions, etc)'")
	flagset.StringVar(&cmd.credentials, "credentials", cmd.credentials, "Path for the 'credentials.json' file. Defaults to "+cmd.credentials)
	flagset.StringVar(&cmd.file, "file", cmd.file, "TSV file name. Defaults to stdout if not supplied")
	flagset.BoolVar(&cmd.withPIN, "with-pin", cmd.withPIN, "Include card keypad PIN code in retrieved membmer information")

	return flagset
}

func (cmd *GetMembers) Execute(args ...interface{}) error {
	options := args[0].(*Options)

	cmd.debug = options.Debug

	// ... check parameters
	if strings.TrimSpace(cmd.credentials) == "" {
		return fmt.Errorf("invalid credentials file")
	}

	// ... get contacts list and member groups
	conf := config.NewConfig()
	if err := conf.Load(options.Config); err != nil {
		return fmt.Errorf("could not load configuration (%v)", err)
	}

	credentials, err := getCredentials(cmd.credentials)
	if err != nil {
		return err
	}

	members, err := getMembers(conf, credentials)
	if err != nil {
		return err
	}

	// ... write to stdout
	asTable := func(m *types.Members) *lib.Table {
		if cmd.withPIN {
			return m.AsTableWithPIN()
		} else {
			return m.AsTable()
		}
	}

	if cmd.file == "" {
		fmt.Fprintln(os.Stdout, string(asTable(members).MarshalTextIndent("  ", " ")))
		return nil
	}

	// ... write to TSV file
	asTSV := func(m *types.Members, w io.Writer) error {
		if cmd.withPIN {
			return m.ToTSVWithPIN(w)
		} else {
			return m.ToTSV(w)
		}
	}

	var b bytes.Buffer
	if err := asTSV(members, &b); err != nil {
		return fmt.Errorf("error creating TSV file (%v)", err)
	}

	if err := write(cmd.file, b.Bytes()); err != nil {
		return err
	}

	infof("Retrieved member list to file %s", cmd.file)

	return nil
}
