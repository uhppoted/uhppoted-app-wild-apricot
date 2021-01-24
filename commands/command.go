package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

const APP = "uhppoted-app-wild-apricot"

type credentials struct {
	Account uint32 `json:"account"`
	APIKey  string `json:"api-key"`
}

func getCredentials(file string) (*credentials, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve credentials (%v)", err)
	}

	c := credentials{}

	if err := json.Unmarshal(bytes, &c); err != nil {
		return nil, fmt.Errorf("Unable to retrieve credentials (%v)", err)
	}

	return &c, nil
}

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
