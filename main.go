package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// get our api query object
	api, err := NewApiClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %s\n", err)
		os.Exit(1)
	}

	// query the api for our drives, hosts, and servers
	cDrive := make(chan *[]Drive)
	cHost := make(chan *[]Host)
	cServer := make(chan *[]Server)
	go api.GoGetDrives(cDrive)
	go api.GoGetHosts(cHost)
	go api.GoGetServers(cServer)
	drives, hosts, servers := <-cDrive, <-cHost, <-cServer
	// fmt.Printf("goDrives are %s\n", drives)
	// fmt.Printf("goHosts are %s\n", hosts)
	// fmt.Printf("goServers are %s\n", servers)

	drivesMap := make(map[string]string)
	for _, d := range *drives {
		drivesMap[d.Drive] = d.Licenses
	}

	// populate report
	licenses := &[]string{
		"msft_lwa_00135",
		"msft_p73_04837",
		"msft_p72_04169",
		"msft_p72_04169",
		"msft_6wc_00002",
		"msft_tfa_00009",
		"msft_228_03159",
		"cpanel_vps_1m",
	}
	r := NewReport(licenses, hosts)

	// fmt.Printf("\n\n")
	for _, s := range *servers {
		for _, dm := range drivesMap {
			if dm == "" {
				continue
			}
			for _, dl := range strings.Split(dm, " ") {
				r.UpdateHost(dl, s.Host, s.Cores)
			}
		}
	}
	r.PrintReport("fuckity.csv")
}
