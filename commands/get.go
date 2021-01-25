package commands

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

var GetCmd = Get{
	workdir:     DEFAULT_WORKDIR,
	credentials: filepath.Join(DEFAULT_WORKDIR, ".wild-apricot", "credentials.json"),
	file:        time.Now().Format("2006-01-02T150405.tsv"),
	debug:       false,
}

type Get struct {
	workdir     string
	credentials string
	file        string
	debug       bool
}

func (cmd *Get) Name() string {
	return "get"
}

func (cmd *Get) Description() string {
	return "Retrieves an access control list from a Wild Apricot member database and stores it to a file"
}

func (cmd *Get) Usage() string {
	return "--credentials <file> --file <file>"
}

func (cmd *Get) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s [--debug] [--config <file>] get [--credentials <file>] [--file <file>]\n", APP)
	fmt.Println()
	fmt.Println("  Downloads an access control list from a Wild Apricot member database and stores it to a TSV file")
	fmt.Println()

	helpOptions(cmd.FlagSet())

	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    uhppote-app-wild-apricot --debug get --credentials ".credentials/wild-apricot.json" --file "example.tsv"`)
	fmt.Println()
}

func (cmd *Get) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("get", flag.ExitOnError)

	flagset.StringVar(&cmd.workdir, "workdir", cmd.workdir, "Directory for working files (tokens, revisions, etc)'")
	flagset.StringVar(&cmd.credentials, "credentials", cmd.credentials, "Path for the 'credentials.json' file. Defaults to "+cmd.credentials)
	flagset.StringVar(&cmd.file, "file", cmd.file, "TSV file name. Defaults to 'ACL - <yyyy-mm-dd HHmmss>.tsv'")

	return flagset
}

func (cmd *Get) Execute(args ...interface{}) error {
	options := args[0].(*Options)

	cmd.debug = options.Debug

	// ... check parameters
	if strings.TrimSpace(cmd.credentials) == "" {
		return fmt.Errorf("Invalid credentials file")
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

	members, err := makeMemberList(contacts, groups)
	if err != nil {
		return err
	}

	//	fmt.Printf("GROUPS:\n")
	//	for _, g := range groups {
	//		fmt.Printf("  %v\n", g)
	//	}

	if cmd.debug {
		if text, err := members.MarshalTextIndent("  "); err == nil {
			fmt.Printf("MEMBERS:\n%s\n", string(text))
		}
	}

	// ... save to TSV file
	//	tmp, err := ioutil.TempFile(os.TempDir(), "ACL")
	//	if err != nil {
	//		return err
	//	}
	//
	//	defer func() {
	//		tmp.Close()
	//		os.Remove(tmp.Name())
	//	}()
	//
	//	if err := sheetToTSV(tmp, response); err != nil {
	//		return fmt.Errorf("Error creating TSV file (%v)", err)
	//	}
	//
	//	tmp.Close()
	//
	//	dir := filepath.Dir(cmd.file)
	//	if err := os.MkdirAll(dir, 0770); err != nil {
	//		return err
	//	}
	//
	//	if err := os.Rename(tmp.Name(), cmd.file); err != nil {
	//		return err
	//	}
	//
	//	info(fmt.Sprintf("Retrieved ACL to file %s\n", cmd.file))
	//
	//	return nil
	return fmt.Errorf("NOT IMPLEMENTED")
}
