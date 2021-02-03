package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/uhppoted/uhppote-core/device"
	"github.com/uhppoted/uhppote-core/uhppote"
	api "github.com/uhppoted/uhppoted-api/acl"
	"github.com/uhppoted/uhppoted-api/config"
)

var LoadACLCmd = LoadACL{
	workdir:     DEFAULT_WORKDIR,
	credentials: filepath.Join(DEFAULT_CONFIG_DIR, ".wild-apricot", "credentials.json"),
	rules:       filepath.Join(DEFAULT_CONFIG_DIR, "wild-apricot.grl"),
	force:       false,
	strict:      false,
	dryrun:      false,
	debug:       false,
}

type LoadACL struct {
	workdir     string
	credentials string
	rules       string
	force       bool
	strict      bool
	dryrun      bool
	debug       bool
}

func (cmd *LoadACL) Name() string {
	return "load-acl"
}

func (cmd *LoadACL) Description() string {
	return "Retrieves an access control list from a Wild Apricot member database and updates the card lists on the configured controllers"
}

func (cmd *LoadACL) Usage() string {
	return "--credentials <file> --rules <url>"
}

func (cmd *LoadACL) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s [--debug] [--config <file>] load-acl [--credentials <file>] [--rules <url>]\n", APP)
	fmt.Println()
	fmt.Println("  Downloads an access control list from a Wild Apricot member database, applies the ACL rules and updates the card lists")
	fmt.Println("  on the configured controllers")
	fmt.Println()

	helpOptions(cmd.FlagSet())

	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println(`    uhppote-app-wild-apricot --debug --config uhppoted.conf load-acl --credentials ".credentials/wild-apricot.json"`)
	fmt.Println()
}

func (cmd *LoadACL) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("load-acl", flag.ExitOnError)

	flagset.StringVar(&cmd.workdir, "workdir", cmd.workdir, "Directory for working files (tokens, revisions, etc)'")
	flagset.StringVar(&cmd.credentials, "credentials", cmd.credentials, "Path for the 'credentials.json' file. Defaults to "+cmd.credentials)
	flagset.StringVar(&cmd.rules, "rules", cmd.rules, "URI for the 'grule' rules file. Support file path, HTTP and HTTPS. Defaults to "+cmd.rules)
	flagset.BoolVar(&cmd.force, "force", cmd.force, "Forces an update, overriding the  version and compare logic")
	flagset.BoolVar(&cmd.strict, "strict", cmd.strict, "Fails with an error if the members list contains duplicate card numbers")
	flagset.BoolVar(&cmd.dryrun, "dry-run", cmd.dryrun, "Simulates a load-acl without making any changes to the access controllers")

	return flagset
}

func (cmd *LoadACL) Execute(args ...interface{}) error {
	options := args[0].(*Options)

	cmd.debug = options.Debug

	// ... check parameters
	if strings.TrimSpace(cmd.credentials) == "" {
		return fmt.Errorf("Invalid credentials file")
	}

	if strings.TrimSpace(cmd.rules) == "" {
		return fmt.Errorf("Invalid rules file")
	}

	// ... locked?
	lockfile, err := cmd.lock()
	if err != nil {
		return err
	}

	defer func() {
		info(fmt.Sprintf("Removing lockfile '%v'", lockfile))
		os.Remove(lockfile)
	}()

	// ... get config, members and rules
	conf := config.NewConfig()
	if err := conf.Load(options.Config); err != nil {
		return fmt.Errorf("Could not load configuration (%v)", err)
	}

	cardNumberField := conf.WildApricot.Fields.CardNumber
	groupDisplayOrder := strings.Split(conf.WildApricot.DisplayOrder.Groups, ",")

	members, err := getMembers(cmd.credentials, cardNumberField, groupDisplayOrder)
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

	// ... load

	u, devices := getDevices(conf, cmd.debug)
	cards := acl.AsTable()

	updated, err := cmd.compare(&u, devices, &cards)
	if err != nil {
		return err
	}

	if cmd.force || updated {
		rpt, err := cmd.load(&u, devices, &cards)
		if err != nil {
			return err
		}

		if rpt != nil {
			for k, v := range rpt {
				for _, err := range v.Errors {
					fatal(fmt.Sprintf("%v  %v", k, err))
				}
			}

			//		if !cmd.nolog {
			//		}
			//
			//		if !cmd.noreport {
			//		}
		}
	} else {
		info("No changes - Nothing to do")
	}

	//	if version != nil {
	//		version.store(cmd.revisions)
	//	}

	return nil
}

func (cmd *LoadACL) compare(u device.IDevice, devices []*uhppote.Device, cards *api.Table) (bool, error) {
	current, err := api.GetACL(u, devices)
	if err != nil {
		return false, err
	}

	acl, _, err := api.ParseTable(cards, devices, cmd.strict)
	if err != nil {
		return false, err
	}

	if acl == nil {
		return false, fmt.Errorf("Error creating ACL from cards (%v)", cards)
	}

	diff, err := api.Compare(current, *acl)
	if err != nil {
		return false, err
	}

	for _, v := range diff {
		if v.HasChanges() {
			return true, nil
		}
	}

	return false, nil
}

func (cmd *LoadACL) load(u device.IDevice, devices []*uhppote.Device, cards *api.Table) (map[uint32]api.Report, error) {
	acl, warnings, err := api.ParseTable(cards, devices, cmd.strict)
	if err != nil {
		return nil, err
	}

	for _, w := range warnings {
		warn(w.Error())
	}

	rpt, err := api.PutACL(u, *acl, cmd.dryrun)
	if err != nil {
		return nil, err
	}

	summary := api.Summarize(rpt)
	format := "%v  unchanged:%v  updated:%v  added:%v  deleted:%v  failed:%v  errors:%v"
	for _, v := range summary {
		info(fmt.Sprintf(format, v.DeviceID, v.Unchanged, v.Updated, v.Added, v.Deleted, v.Failed, v.Errored+len(warnings)))
	}

	return rpt, nil
}

func (cmd *LoadACL) lock() (string, error) {
	lockfile := filepath.Join(cmd.workdir, ".wild-apricot", "uhppoted-app-wild-apricot.lock")
	pid := fmt.Sprintf("%d\n", os.Getpid())

	if err := os.MkdirAll(filepath.Dir(lockfile), 0770); err != nil {
		return "", fmt.Errorf("Unable to create directory '%v' for lockfile (%v)", lockfile, err)
	}

	if _, err := os.Stat(lockfile); err == nil {
		return "", fmt.Errorf("Locked by '%v'", lockfile)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("Error checking PID lockfile '%v' (%v)", lockfile, err)
	}

	if err := ioutil.WriteFile(lockfile, []byte(pid), 0660); err != nil {
		return "", fmt.Errorf("Unable to create lockfile '%v' (%v)", lockfile, err)
	}

	return lockfile, nil
}
