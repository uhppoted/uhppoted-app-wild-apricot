package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-app-wild-apricot/log"
	"github.com/uhppoted/uhppoted-lib/config"
	lib "github.com/uhppoted/uhppoted-lib/os"
)

const APP = "uhppoted-app-wild-apricot"

type Options struct {
	Config string
	Debug  bool
}

func helpOptions(flagset *flag.FlagSet) {
	count := 0
	flag.VisitAll(func(f *flag.Flag) {
		count++
	})

	flagset.VisitAll(func(f *flag.Flag) {
		fmt.Printf("    --%-13s %s\n", f.Name, f.Usage)
	})

	if count > 0 {
		fmt.Println()
		fmt.Println("  Options:")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf("    --%-13s %s\n", f.Name, f.Usage)
		})
	}
}

func getDevices(conf *config.Config, debug bool) (uhppote.IUHPPOTE, []uhppote.Device) {
	bind, broadcast, listen := config.DefaultIpAddresses()

	if conf.BindAddress != nil {
		bind = *conf.BindAddress
	}

	if conf.BroadcastAddress != nil {
		broadcast = *conf.BroadcastAddress
	}

	if conf.ListenAddress != nil {
		listen = *conf.ListenAddress
	}

	controllers := conf.Devices.ToControllers()

	u := uhppote.NewUHPPOTE(bind, broadcast, listen, conf.Timeout, controllers, false)

	return u, controllers
}

func write(file string, bytes []byte) error {
	tmp, err := os.CreateTemp(os.TempDir(), "ACL")
	if err != nil {
		return err
	}

	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	fmt.Fprintf(tmp, "%s", string(bytes))
	tmp.Close()

	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0770); err != nil {
		return err
	}

	if err := lib.Rename(tmp.Name(), file); err != nil {
		return err
	}

	return nil
}

func normalise(v string) string {
	re := regexp.MustCompile(`[^a-z1-9]`)

	return re.ReplaceAllString(strings.ToLower(v), "")
}

func clean(v string) string {
	return strings.TrimSpace(v)
}

func debugf(format string, args ...any) {
	log.Debugf(format, args...)
}

func infof(format string, args ...any) {
	log.Infof(format, args...)
}

func warnf(format string, args ...any) {
	log.Warnf(format, args...)
}

func errorf(format string, args ...any) {
	log.Errorf(format, args...)
}
