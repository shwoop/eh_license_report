package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
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
	// map[LICENSES]map[HOST]ReportHost
	Report map[string]map[string]ReportHost
}

func NewReport(licenses *[]string, hosts *[]Host) *Report {
	r := Report{}
	r.Report = make(map[string]map[string]ReportHost)
	for _, l := range *licenses {
		emptyHosts := make(map[string]ReportHost)
		for _, h := range *hosts {
			emptyHosts[h.Host] = ReportHost{h.Host, h.Cores, 0, 0, 0, 0}
		}
		r.Report[l] = emptyHosts
	}
	return &r
}

func (r *Report) sqlRounding(x int) int {
	if x < 4 {
		return 4
	} else if x%2 == 1 {
		return x + 1
	}
	return x
}

func (r *Report) UpdateHost(license, host string, cores int) {
	sl := strings.Split(license, ":")
	license = sl[0]
	licenseCnt := 1
	if len(sl) > 1 {
		licenseCnt, _ = strconv.Atoi(sl[1])
	}
	if _, ok := r.Report[license]; !ok {
		fmt.Fprintf(os.Stderr, "Error: Unknown license encountered: ", license)
		return
	}
	if _, ok := r.Report[license][host]; !ok {
		fmt.Fprintf(os.Stderr, "Error: Unknown host encountered: ", host)
		return
	}
	rh := r.Report[license][host]
	rh.Servers++
	rh.Licenses += licenseCnt
	rh.Vcores += cores
	rh.SqlRoundedCores += r.sqlRounding(cores)
	r.Report[license][host] = rh
}

func stdOutput(line []string) {
	fmt.Println(strings.Join(line, ", "))
}

func makeCsvOutput(filename string) (output func([]string), closer func()) {
	file, _ := os.Create(filename)
	writer := csv.NewWriter(file)
	output = func(line []string) {
		err := writer.Write(line)
		if err != nil {
			fmt.Printf("Writer had error: %s", err)
		}
	}
	closer = func() {
		writer.Flush()
		file.Close()
	}
	return
}

func (r *Report) PrintReport(filename string) {
	var closer func()
	var output func([]string)
	if filename == "" {
		output = stdOutput
		closer = func() {}
	} else {
		output, closer = makeCsvOutput(filename)
	}
	defer closer()

	output([]string{
		"license",
		"host",
		"cores",
		"vcores",
		"servers",
		"licenses",
		"sql_adjusted",
	})
	for license, reportHost := range r.Report {
		for _, host := range reportHost {
			output([]string{
				license,
				host.Host,
				strconv.Itoa(host.Cores),
				strconv.Itoa(host.Vcores),
				strconv.Itoa(host.Servers),
				strconv.Itoa(host.Licenses),
				strconv.Itoa(host.SqlRoundedCores),
			})
		}
	}
}

func (r *Report) PopulateReport(
	servers *[]Server,
	driveLicenses *map[string]string,
) {
	for _, s := range *servers {
		for _, dl := range *driveLicenses {
			if dl == "" {
				continue
			}
			for _, l := range strings.Split(dl, " ") {
				r.UpdateHost(l, s.Host, s.Cores)
			}
		}
	}
}
