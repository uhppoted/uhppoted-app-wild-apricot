package commands

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/uhppoted/uhppoted-lib/config"
)

var GetGroupsCmd = GetGroups{
	workdir:     DEFAULT_WORKDIR,
	credentials: filepath.Join(DEFAULT_CONFIG_DIR, ".wild-apricot", "credentials.json"),
	debug:       false,
}

type GetGroups struct {
	workdir     string
	credentials string
	file        string
	debug       bool
}

func (cmd *GetGroups) Name() string {
	return "get-groups"
}

func (cmd *GetGroups) Description() string {
	return "Retrieves a list of groups from a Wild Apricot database and stores it to a file"
}

func (cmd *GetGroups) Usage() string {
	return "--credentials <file> --file <file>"
}

func (cmd *GetGroups) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s [--debug] [--config <file>] get-groups [--credentials <file>] [--file <file>]\n", APP)
	fmt.Println()
	fmt.Println("  Downloads a list of member groups from a Wild Apricot member database and (optionally) stores it to a TSV file")
	fmt.Println()

	helpOptions(cmd.FlagSet())

	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    uhppote-app-wild-apricot --debug get-groups --credentials ".credentials/wild-apricot.json" --file "groups.tsv"`)
	fmt.Println()
}

func (cmd *GetGroups) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("get-groups", flag.ExitOnError)

	flagset.StringVar(&cmd.workdir, "workdir", cmd.workdir, "Directory for working files (tokens, revisions, etc)'")
	flagset.StringVar(&cmd.credentials, "credentials", cmd.credentials, "Path for the 'credentials.json' file. Defaults to "+cmd.credentials)
	flagset.StringVar(&cmd.file, "file", cmd.file, "TSV file name. Defaults to stdout if not supplied")

	return flagset
}

func (cmd *GetGroups) Execute(args ...interface{}) error {
	options := args[0].(*Options)

	cmd.debug = options.Debug

	// ... check parameters
	if strings.TrimSpace(cmd.credentials) == "" {
		return fmt.Errorf("invalid credentials file")
	}

	// ... get member groups
	conf := config.NewConfig()
	if err := conf.Load(options.Config); err != nil {
		return fmt.Errorf("could not load configuration (%v)", err)
	}

	credentials, err := getCredentials(cmd.credentials)
	if err != nil {
		return err
	}

	groups, err := getGroups(conf, credentials)
	if err != nil {
		return err
	}

	// ... write to stdout
	if cmd.file == "" {
		fmt.Fprintln(os.Stdout, string(groups.AsTable().MarshalTextIndent("  ", " ")))
		return nil
	}

	// ... write to TSV file
	var b bytes.Buffer
	if err := groups.AsTable().ToTSV(&b); err != nil {
		return fmt.Errorf("error creating TSV file (%v)", err)
	}

	if err := write(cmd.file, b.Bytes()); err != nil {
		return err
	}

	infof("Retrieved groups list to file %s", cmd.file)

	return nil
}
