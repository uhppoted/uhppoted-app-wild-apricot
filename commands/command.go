package commands

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
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

func fetch(uri string) ([]byte, error) {
	f := ioutil.ReadFile
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		f = fetchHTTP
	} else if strings.HasPrefix(uri, "file://") {
		f = fetchFile
	}

	return f(uri)
}

// Ref. https://stackoverflow.com/questions/18177419/download-public-file-from-google-drive-golang
//      Need to use https://drive.google.com/uc?export=download&id=<ID> for Google Drive shares.
func fetchHTTP(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var b bytes.Buffer
	if _, err = io.Copy(&b, response.Body); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func fetchFile(url string) ([]byte, error) {
	match := regexp.MustCompile("^file://(.*)").FindStringSubmatch(url)
	if len(match) != 2 {
		return nil, fmt.Errorf("Invalid file URI (%s)", url)
	}

	return ioutil.ReadFile(match[1])
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
