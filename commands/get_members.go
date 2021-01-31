package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

var GetMembersCmd = GetMembers{
	workdir:     DEFAULT_WORKDIR,
	credentials: filepath.Join(DEFAULT_CONFIG_DIR, ".wild-apricot", "credentials.json"),
	file:        time.Now().Format("members 2006-01-02T150405.tsv"),
	debug:       false,
}

type GetMembers struct {
	workdir     string
	credentials string
	file        string
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
	fmt.Printf("  Usage: %s [--debug] [--config <file>] get-members [--credentials <file>] [--file <file>]\n", APP)
	fmt.Println()
	fmt.Println("  Downloads an access control list from a Wild Apricot member database, applies the ACL rules and")
	fmt.Println("  stores the generated access control list to a TSV file")
	fmt.Println()

	helpOptions(cmd.FlagSet())

	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    uhppote-app-wild-apricot --debug get-members --credentials ".credentials/wild-apricot.json" \"`)
	fmt.Println(`                                                 --file "example.tsv"`)
	fmt.Println()
}

func (cmd *GetMembers) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("get-members", flag.ExitOnError)

	flagset.StringVar(&cmd.workdir, "workdir", cmd.workdir, "Directory for working files (tokens, revisions, etc)'")
	flagset.StringVar(&cmd.credentials, "credentials", cmd.credentials, "Path for the 'credentials.json' file. Defaults to "+cmd.credentials)
	flagset.StringVar(&cmd.file, "file", cmd.file, "TSV file name. Defaults to 'ACL - <yyyy-mm-dd HHmmss>.tsv'")

	return flagset
}

func (cmd *GetMembers) Execute(args ...interface{}) error {
	options := args[0].(*Options)

	cmd.debug = options.Debug

	// ... check parameters
	if strings.TrimSpace(cmd.credentials) == "" {
		return fmt.Errorf("Invalid credentials file")
	}

	if strings.TrimSpace(cmd.file) == "" {
		return fmt.Errorf("Invalid output file")
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

	// ... save to TSV file
	tmp, err := ioutil.TempFile(os.TempDir(), "ACL")
	if err != nil {
		return err
	}

	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	if err := members.ToTSV(tmp); err != nil {
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

	info(fmt.Sprintf("Retrieved member list to file %s\n", cmd.file))

	return nil
}
