package main

import "fmt"

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

func (r *Report) PrintReport() {
	fmt.Println("license,host,cores,vcores,servers,sql_adjusted\n")
	for license, reportHost := range r.Licenses {
		for _, host := range reportHost {
			fmt.Printf(
				"%s,%s,%d,%d,%d,%d\n",
				license,
				host.Host,
				host.Cores,
				host.Vcores,
				host.Servers,
				host.SqlRoundedCores,
			)
		}
	}
}

// TODO:  Csv lib
