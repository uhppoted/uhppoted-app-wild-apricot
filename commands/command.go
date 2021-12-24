package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/config"
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
	keys := []uint32{}
	for id, _ := range conf.Devices {
		keys = append(keys, id)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

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

	devices := []uhppote.Device{}
	for _, id := range keys {
		d := conf.Devices[id]

		if device := uhppote.NewDevice(d.Name, id, d.Address, d.Doors); device != nil {
			devices = append(devices, *device)
		}
	}

	u := uhppote.NewUHPPOTE(bind, broadcast, listen, conf.Timeout, devices, debug)

	return u, devices
}

func write(file string, bytes []byte) error {
	tmp, err := ioutil.TempFile(os.TempDir(), "ACL")
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

	if err := os.Rename(tmp.Name(), file); err != nil {
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

func debug(msg string) {
	log.Printf("%-5s %s", "DEBUG", msg)
}

func info(msg string) {
	log.Printf("%-5s %s", "INFO", msg)
}

func warn(msg string) {
	log.Printf("%-5s %s", "WARN", msg)
}

func fatal(msg string) {
	log.Printf("%-5s %s", "ERROR", msg)
}
