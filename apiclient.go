package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "io"
	"net/http"
	"os"
	"strings"
)

type ApiClient struct {
	usr string
	pwd string
	uri string
	cli *http.Client
}

func NewApiClient() (*ApiClient, error) {
	ehuri := os.Getenv("EHURI")
	ehauth := os.Getenv("EHAUTH")
	if ehuri == "" || ehauth == "" {
		fmt.Fprintf(os.Stderr, "EHURI or EHAUTH not set\n")
		os.Exit(1)
	}
	spl := strings.Split(ehauth, ":")
	if len(spl) != 2 {
		fmt.Fprintf(os.Stderr, "EHAUTH not valid\n")
		os.Exit(1)
	}
	usr, pwd := spl[0], spl[1]
	// Ignore for testing
	// if usr != "00000000-0000-0000-0000-000000000000" {
	// 	fmt.Fprintf(os.Stderr, "EHAUTH not set as su, please amend.\n")
	// 	os.Exit(1)
	// }
	return &ApiClient{usr: usr, pwd: pwd, uri: ehuri, cli: &http.Client{}}, nil
}

func (ac *ApiClient) Get(endpoint string) (*json.Decoder, error) {
	uri := ac.uri + endpoint
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return &json.Decoder{}, err
	}
	req.SetBasicAuth(ac.usr, ac.pwd)
	req.Header.Set("Accept", "application/json")
	resp, err := ac.cli.Do(req)
	if err != nil {
		return &json.Decoder{}, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return json.NewDecoder(buf), nil
}

func (ac *ApiClient) GetEndpoint(endpoint string) (*[]map[string]interface{}, error) {
	dec, err := ac.Get(endpoint)
	if err != nil {
		return nil, err
	}
	if _, err = dec.Token(); err != nil {
		return nil, err
	}
	var objs []map[string]interface{}
	for dec.More() {
		var obj map[string]interface{}
		if err = dec.Decode(&obj); err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}

	return &objs, nil
}
