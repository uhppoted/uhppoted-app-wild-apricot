package commands

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/uhppoted/uhppoted-api/config"
)

var GetDoorsCmd = GetDoors{
	workdir: DEFAULT_WORKDIR,
	debug:   false,
}

type GetDoors struct {
	workdir string
	file    string
	debug   bool
}

func (cmd *GetDoors) Name() string {
	return "get-doors"
}

func (cmd *GetDoors) Description() string {
	return "Retrieves a list of doors from the uhppoted.conf file and stores it to a file"
}

func (cmd *GetDoors) Usage() string {
	return "--credentials <file> --file <file>"
}

func (cmd *GetDoors) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s [--debug] [--config <file>] get-doors [--file <file>]\n", APP)
	fmt.Println()
	fmt.Println("  Downloads a list of doors from the uhppoted.conf file and (optionally) stores it to a TSV file")
	fmt.Println()

	helpOptions(cmd.FlagSet())

	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    uhppote-app-wild-apricot --debug get-doors --file "groups.tsv"`)
	fmt.Println()
}

func (cmd *GetDoors) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("get-doors", flag.ExitOnError)

	flagset.StringVar(&cmd.workdir, "workdir", cmd.workdir, "Directory for working files (tokens, revisions, etc)'")
	flagset.StringVar(&cmd.file, "file", cmd.file, "TSV file name. Defaults to stdout if not supplied")

	return flagset
}

func (cmd *GetDoors) Execute(args ...interface{}) error {
	options := args[0].(*Options)

	cmd.debug = options.Debug

	// ... get doors
	conf := config.NewConfig()
	if err := conf.Load(options.Config); err != nil {
		return fmt.Errorf("Could not load configuration (%v)", err)
	}

	doors, err := getDoors(conf)
	if err != nil {
		return err
	}

	// ... write to stdout
	if cmd.file == "" {
		fmt.Fprintln(os.Stdout, string(doors.AsTable().MarshalTextIndent("  ", " ")))
		return nil
	}

	// ... write to TSV file
	var b bytes.Buffer
	if err := doors.AsTable().ToTSV(&b); err != nil {
		return fmt.Errorf("Error creating TSV file (%v)", err)
	}

	if err := write(cmd.file, b.Bytes()); err != nil {
		return err
	}

	info(fmt.Sprintf("Extracted doors list to file %s\n", cmd.file))

	return nil
}
