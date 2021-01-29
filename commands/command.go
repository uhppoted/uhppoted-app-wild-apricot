package commands

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/uhppoted/uhppoted-api/config"
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

func getDoors(file string) ([]string, error) {
	conf := config.NewConfig()
	if err := conf.Load(file); err != nil {
		return nil, fmt.Errorf("WARN  Could not load configuration (%v)", err)
	}

	doors := map[string]string{}
	for _, device := range conf.Devices {
		for _, d := range device.Doors {
			if _, ok := doors[normalise(d)]; ok {
				return nil, fmt.Errorf("WARN  Duplicate door in configuration (%v)", d)
			}

			doors[normalise(d)] = clean(d)
		}
	}

	list := []string{}
	for _, d := range doors {
		list = append(list, d)
	}

	sort.SliceStable(list, func(i, j int) bool { return list[i] < list[j] })

	return list, nil
}

func normalise(v string) string {
	return strings.ToLower(strings.ReplaceAll(v, " ", ""))
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
