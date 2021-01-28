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

func fetch(uri string) ([]byte, error) {
	f := ioutil.ReadFile
	//	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
	//		f = fetchHTTP
	//	} else if strings.HasPrefix(uri, "s3://") {
	//		f = fetchS3
	//	} else if strings.HasPrefix(uri, "file://") {
	//		f = fetchFile
	//	}

	return f(uri)
}

//func fetchFile(url string) ([]byte, error) {
//	match := regexp.MustCompile("^file://(.*)").FindStringSubmatch(url)
//	if len(match) != 2 {
//		return nil, fmt.Errorf("Invalid file URI (%s)", url)
//	}
//
//	return ioutil.ReadFile(match[1])
//}

//func fetchHTTP(url string) ([]byte, error) {
//	response, err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//
//	defer response.Body.Close()
//
//	var b bytes.Buffer
//	if _, err = io.Copy(&b, response.Body); err != nil {
//		return nil, err
//	}
//
//	return b.Bytes(), nil
//}
//
//func fetchS3(url, config, profile, region string) ([]byte, error) {
//	match := regexp.MustCompile("^s3://(.*?)/(.*)").FindStringSubmatch(url)
//	if len(match) != 3 {
//		return nil, fmt.Errorf("Invalid S3 URI (%s)", url)
//	}
//
//	bucket := match[1]
//	key := match[2]
//	object := s3.GetObjectInput{
//		Bucket: aws.String(bucket),
//		Key:    aws.String(key),
//	}
//
//	cfg := aws.NewConfig().
//		WithCredentials(credentials.NewSharedCredentials(config, profile)).
//		WithRegion(region)
//
//	ss := session.Must(session.NewSession(cfg))
//
//	buffer := make([]byte, 1024)
//	b := aws.NewWriteAtBuffer(buffer)
//	if _, err := s3manager.NewDownloader(ss).Download(b, &object); err != nil {
//		return nil, err
//	}
//
//	return b.Bytes(), nil
//}

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
