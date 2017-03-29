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

func (ac *ApiClient) GetHosts() (*[]Host, error) {
	dec, err := ac.Get("hosts/info/full")
	if err != nil {
		return nil, err
	}
	var hosts []Host
	_, err = dec.Token()
	for dec.More() {
		var h Host
		if err = dec.Decode(&h); err != nil {
			return nil, err
		}
		hosts = append(hosts, h)
	}
	return &hosts, nil
}

func (ac *ApiClient) GoGetHosts(c chan *[]Host) {
	hosts, err := ac.GetHosts()
	if err != nil {
		c <- nil
	} else {
		c <- hosts
	}
	close(c)
}

func (ac *ApiClient) GetDrives() (*[]Drive, error) {
	dec, err := ac.Get("drives/info/full")
	if err != nil {
		return nil, err
	}
	var drives []Drive
	_, err = dec.Token()
	for dec.More() {
		var d Drive
		if err = dec.Decode(&d); err != nil {
			return nil, err
		}
		drives = append(drives, d)
	}
	return &drives, nil
}

func (ac *ApiClient) GoGetDrives(c chan *[]Drive) {
	drives, err := ac.GetDrives()
	if err != nil {
		c <- nil
	} else {
		c <- drives
	}
	close(c)
}

func (ac *ApiClient) GetServers() (*[]Server, error) {
	objs, err := ac.GetEndpoint("servers/info/full")
	if err != nil {
		return nil, err
	}
	var drives []string
	var servers []Server
	for _, v := range *objs {
		if v["type"] != "vm" || v["status"] != "active" {
			continue
		}
		for key := range v {
			switch str := fmt.Sprint(v[key]); {
			case strings.HasPrefix(key, "ide") && ValidateUuid(str):
				drives = append(drives, str)
			case strings.HasPrefix(key, "block") && ValidateUuid(str):
				drives = append(drives, str)
			case strings.HasPrefix(key, "ata") && ValidateUuid(str):
				drives = append(drives, str)
			case strings.HasPrefix(key, "scsi") && ValidateUuid(str):
				drives = append(drives, str)
			}
		}
		cores, _ := ToInt(v["smp:cores"])
		servers = append(
			servers,
			Server{
				Server: fmt.Sprint(v["server"]),
				Host:   fmt.Sprint(v["host"]),
				Cores:  cores,
				Drives: drives,
			},
		)
		drives = nil
	}
	return &servers, err
}

func (ac *ApiClient) GoGetServers(c chan *[]Server) {
	servers, err := ac.GetServers()
	if err != nil {
		c <- nil
	} else {
		c <- servers
	}
	close(c)
}

func (ac *ApiClient) GetEndpoint(endpoint string) (
	*[]map[string]interface{},
	error,
) {
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

func (ac *ApiClient) GoGetEndpoint(
	endpoint string,
	c chan *[]map[string]interface{},
) {
	objs, err := ac.GetEndpoint(endpoint)
	if err != nil {
		c <- nil
	} else {
		c <- objs
	}
	close(c)
}
