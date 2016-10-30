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

func (ac *ApiClient) Get(slug string) (*json.Decoder, error) {
	fmt.Printf("entering GET\n")
	uri := ac.uri + slug
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

func (ac *ApiClient) GetHosts() (*[]Host, error) {
	dec, err := ac.Get("hosts/info/full")
	if err != nil {
		return nil, err
	}
	var hosts []Host
	var h Host
	_, err = dec.Token()
	for dec.More() {
		if err = dec.Decode(&h); err != nil {
			return nil, err
		}
		hosts = append(hosts, h)
	}
	return &hosts, nil
}

func (ac *ApiClient) GetServers() (*[]Server, error) {
	dec, err := ac.Get("servers/info/full")
	if err != nil {
		return nil, err
	}
	// var servers []Server
	// var s Server
	// var s interface{}
	_, err = dec.Token()
	for dec.More() {
		var s map[string]interface{}
		if err = dec.Decode(&s); err != nil {
			return nil, err
		}
		fmt.Printf("name %s\n", s["name"])
		fmt.Printf("cores %s\n", s["smp:cores"])
		// fmt.Printf("server %s\n", s)

		// servers = append(servers, s)
	}
	//return &servers, nil
	return nil, nil
}
