package main

import (
	"fmt"
	// 	"io"
	// 	"net/http"
	"os"
	// 	"strings"
)

func main() {
	api, err := NewApiClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %s\n", err)
		os.Exit(1)
	}
	hosts, err := api.GetHosts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("hosts are %s", hosts)
	servers, err := api.GetServers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("first server is %s", servers)
}
