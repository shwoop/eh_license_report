package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type ReportHost struct {
	Host            string
	Cores           int
	Vcores          int
	Servers         int
	SqlRoundedCores int
}

type Report struct {
	Licenses map[string]map[string]ReportHost
}

func NewReport(licenses *[]string, hosts *[]Host) *Report {
	r := Report{}
	r.Licenses = make(map[string]map[string]ReportHost)
	for _, l := range *licenses {
		emptyHosts := make(map[string]ReportHost)
		for _, h := range *hosts {
			emptyHosts[h.Host] = ReportHost{h.Host, h.Cores, 0, 0, 0}
		}
		r.Licenses[l] = emptyHosts
	}
	return &r
}

func (r *Report) sqlRounding(c int) int {
	if c < 4 {
		return 4
	} else if c%2 == 1 {
		return c + 1
	}
	return c
}

func (r *Report) UpdateHost(l, h string, c int) {
	rh := r.Licenses[l][h]
	rh.Servers++
	rh.Vcores += c
	rh.SqlRoundedCores += r.sqlRounding(c)
	r.Licenses[l][h] = rh
}

func stdOutput(line string) {
	fmt.Printf(line)
}

func makeCsvOutput(filename string) (func(string), func()) {
	file, _ := os.Create(filename)
	writer := csv.NewWriter(file)
	return func(line string) {
			writer.Write([]string{line})
		}, func() {
			file.Close()
			writer.Flush()
		}
}

func (r *Report) PrintReport(filename *string) {
	var closer func()
	var output func(string)
	if filename == nil {
		output = stdOutput
		closer = nil
	} else {
		output, closer = makeCsvOutput(*filename)
	}
	// 	fmt.Println("license,host,cores,vcores,servers,sql_adjusted\n")
	output("license,host,cores,vcores,servers,sql_adjusted\n")
	for license, reportHost := range r.Licenses {
		for _, host := range reportHost {
			output(fmt.Sprintf(
				"%s,%s,%d,%d,%d,%d\n",
				license,
				host.Host,
				host.Cores,
				host.Vcores,
				host.Servers,
				host.SqlRoundedCores,
			))
			// fmt.Printf(
			// 	"%s,%s,%d,%d,%d,%d\n",
			// 	license,
			// 	host.Host,
			// 	host.Cores,
			// 	host.Vcores,
			// 	host.Servers,
			// 	host.SqlRoundedCores,
			// )
		}
	}
	if closer != nil {
		closer()
	}
}

// TODO:  Csv lib
