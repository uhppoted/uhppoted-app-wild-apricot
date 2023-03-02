package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func fetch(uri string) ([]byte, error) {
	f := os.ReadFile
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		f = fetchHTTP
	} else if strings.HasPrefix(uri, "file://") {
		f = fetchFile
	}

	return f(uri)
}

// Ref. https://stackoverflow.com/questions/18177419/download-public-file-from-google-drive-golang
//
//	Need to use https://drive.google.com/uc?export=download&id=<ID> for Google Drive shares.
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
		return nil, fmt.Errorf("invalid file URI (%s)", url)
	}

	return os.ReadFile(match[1])
}
