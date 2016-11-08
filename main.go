package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fileName := flag.String("f", "", "csv filename rather than stdout")
	flag.Parse()
	fmt.Println("filename: ", *fileName)

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

	// Pack slice of drives into Drive-License map
	driveLicenses := make(map[string]string)
	for _, d := range *drives {
		driveLicenses[d.Drive] = d.Licenses
	}

	// Populate report
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
	r.PopulateReport(servers, &driveLicenses)
	r.PrintReport(*fileName)
}
