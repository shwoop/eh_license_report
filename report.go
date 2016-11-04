package main

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
	emptyHosts := make(map[string]ReportHost)
	for _, h := range *hosts {
		emptyHosts[h.Host] = ReportHost{h.Host, h.Cores, 0, 0, 0}
	}
	for _, l := range *licenses {
		tmp := emptyHosts
		r.Licenses[l] = tmp
	}
	return &r
}

// TODO:  UpdateHost (license, host, serverno, vcoreno, ...
// TODO:  PrintReport
// TODO:  Csv lib
