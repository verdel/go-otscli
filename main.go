package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Secret struct {
	State               string   `json:"state"`
	Custid              string   `json:"custid"`
	Metadata_ttl        int64    `json:"metadata_ttl"`
	Ttl                 int64    `json:"ttl"`
	Secret_ttl          int64    `json:"secret_ttl"`
	Updated             float64  `json:"updated"`
	Created             float64  `json:"created"`
	Recipient           []string `json:"recipient"`
	Passphrase_required bool     `json:"passphrase_required"`
	Metadata_key        string   `json:"metadata_key"`
	Secret_key          string   `json:"secret_key"`
}

func getSecret(body []byte) (*Secret, error) {
	var s = new(Secret)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

func main() {
	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Usage = Usage
	urlPtr := flag.String("u", "https://pwd.rkomi.ru", "ots service url")
	secretPtr := flag.String("s", "", "the secret value which is encrypted before being stored.")
	passwordPtr := flag.String("p", "", "a string that the recipient must know to view the secret.")
	flag.Parse()

	if *urlPtr != "" && *secretPtr != "" {
		resource := "/api/v1/share"
		data := url.Values{}
		data.Set("secret", *secretPtr)
		if *passwordPtr != "" {
			data.Add("passphrase", *passwordPtr)
		}
		u, _ := url.ParseRequestURI(*urlPtr)
		u.Path = resource
		urlStr := fmt.Sprintf("%v", u)

		client := &http.Client{}
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		resp, err := client.Do(r)
		if err != nil {
			panic(err.Error())
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}

		s, err := getSecret(body)
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("%s/secret/%s\n", *urlPtr, s.Secret_key)
	}
}
