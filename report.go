package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type ReportHost struct {
	Host            string
	Cores           int
	Vcores          int
	Servers         int
	Licenses        int
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
			emptyHosts[h.Host] = ReportHost{h.Host, h.Cores, 0, 0, 0, 0}
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
	sl := strings.Split(l, ":")
	l = sl[0]
	ln := 1
	if len(sl) > 1 {
		ln, _ = ToInt(sl[1])
	}
	rh := r.Licenses[l][h]
	rh.Servers++
	rh.Licenses += ln
	rh.Vcores += c
	rh.SqlRoundedCores += r.sqlRounding(c)
	r.Licenses[l][h] = rh
}

func stdOutput(line string) {
	fmt.Println(line)
}

func makeCsvOutput(filename string) (func(string), func()) {
	file, _ := os.Create(filename)
	writer := csv.NewWriter(file)
	return func(line string) {
			err := writer.Write(strings.Split(line, ","))
			if err != nil {
				fmt.Println("got error:", err)
			}
		}, func() {
			writer.Flush()
			file.Close()
		}
}

func (r *Report) PrintReport(filename string) {
	var closer func()
	var output func(string)
	if filename == "" {
		output = stdOutput
		closer = func() {}
	} else {
		output, closer = makeCsvOutput(filename)
	}
	defer closer()

	output("license,host,cores,vcores,servers,licenses,sql_adjusted")
	for license, reportHost := range r.Licenses {
		for _, host := range reportHost {
			output(fmt.Sprintf(
				"%s,%s,%d,%d,%d,%d,%d",
				license,
				host.Host,
				host.Cores,
				host.Vcores,
				host.Servers,
				host.Licenses,
				host.SqlRoundedCores,
			))
		}
	}
}
