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
	hosts, err := api.GetEndpoint("hosts/info/full")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("hosts are %s", hosts)
	servers, err := api.GetEndpoint("servers/info/full")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("first server is %s", servers)
}
