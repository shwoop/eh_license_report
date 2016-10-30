package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	ehuri := os.Getenv("EHURI")
	ehauth := os.Getenv("EHAUTH")
	if ehuri == "" || ehauth == "" {
		fmt.Fprintf(os.Stderr, "EHURI or EHAUTH not set\n")
		os.Exit(1)
	}
	usr := strings.Split(ehauth, ":")[0]
	// if usr != "00000000-0000-0000-0000-000000000000" {
	// 	fmt.Fprintf(os.Stderr, "EHAUTH not set as su, please amend.\n")
	// 	os.Exit(1)
	// }
	fmt.Printf("ehuri = %s, ehauth = %s, usr = %s", ehuri, ehauth, usr)
	resp, err := http.Get(ehuri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query %s. %s\n", ehuri, err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to output response body. %s\n", err)
		os.Exit(1)
	}
}
