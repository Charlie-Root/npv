package parser

import (
	"io/ioutil"
	"net"
	"sort"
	"strconv"

	"github.com/Charlie-Root/mtrview/pkg/db"
	"github.com/Charlie-Root/mtrview/pkg/logging"
	"github.com/Charlie-Root/mtrview/pkg/mtr"
	"github.com/bitly/go-simplejson"
	"github.com/jamesog/iptoasn"
)

var logger = logging.NewLogger("parser")

func ParseJsonFile(jsonFile string) []string {

	iplist := []string{}

	file, err := ioutil.ReadFile("hosts.json")
	if err != nil {
		logger.Error(err.Error())
	}
	js, err := simplejson.NewJson(file)

	hostsArr := js.Get("hosts").MustArray()

	for _, ele := range hostsArr {
		t, _ := ele.(string)
		iplist = append(iplist, string(t))
	}

	return iplist

}

func ParseResults(m mtr.MTR, database db.DB) {

	// Get the keys from the map and sort them.
	keys := make([]int, 0, len(m.Statistic))
	for key := range m.Statistic {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for _, key := range keys {

		if len(m.Statistic[key].Targets) == 1 {

			var searchHostname string
			if m.Statistic[key].Targets[0] == "" {
				searchHostname = "***"
			} else {
				searchHostname = m.Statistic[key].Targets[0]
			}

			_, host_address, host_count := database.GetHostByAddress(searchHostname, m.Statistic[key].TTL)
			//ug(host_address)

			if host_address != "" {
				//logger.Warning("update hostcount")
				database.UpdateHostCount(m.Statistic[key].Targets[0], host_count)
			} else {
				//logger.Warning("add host to db")
				//logger.Debug(searchHostname)

				host := db.Host{
					Name:      searchHostname,
					Address:   searchHostname,
					HostPTR:   getPTR(searchHostname),
					HostASN:   getASN(searchHostname),
					HostTTL:   m.Statistic[key].TTL,
					HostCount: 1,
				}

				if err := database.InsertHost(host); err != nil {
					logger.Error(err.Error())
					return
				}
			}

		} else {
			logger.Error(strconv.Itoa(m.Statistic[key].TTL) + " - Looks like we found multiple possible targets ??")

		}
		if value, ok := m.Statistic[key-1]; ok {
			curAddress := m.Statistic[key].Targets[0]
			prevAddress := value.Targets[0]

			curTTL := m.Statistic[key].TTL
			//prevTTL := value.TTL

			if curAddress == "" {
				curAddress = "***"
			}
			if prevAddress == "" {
				prevAddress = "***"
			}

			//cur_host_id, _, _ := database.GetHostByAddress(curAddress, curTTL)
			//prev_host_id, _, _ := database.GetHostByAddress(prevAddress, prevTTL)

			host_from, host_to, _, linkCounter := database.GetLink(prevAddress, curAddress, curTTL)

			if host_from == "" && host_to == "" {
				// add the link here
				//logger.Warning("Add the link with TTL: " + strconv.Itoa(curTTL))

				link := db.Link{
					Source:      prevAddress,
					Target:      curAddress,
					TargetTTL:   curTTL,
					TargetLoss:  m.Statistic[key].Lost,
					TargetSNT:   m.Statistic[key].Sent,
					TargetLast:  int(m.Statistic[key].Last.Elapsed),
					TargetAVG:   int(m.Statistic[key].Avg()),
					TargetBest:  int(m.Statistic[key].Best.Elapsed),
					TargetWRST:  int(m.Statistic[key].Worst.Elapsed),
					TargetStDev: int(m.Statistic[key].Stdev()),
					LinkCount:   1,
				}

				if err := database.InsertLink(link); err != nil {
					logger.Error(err.Error())
					return
				}
			} else {
				database.UpdateLinkCount(prevAddress, curAddress, curTTL, linkCounter)
				//logger.Warning("Update the link")
			}

		}

	}

}

func getPTR(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return ip
	}
	if len(names) == 0 {
		return ip
	}
	return names[0]
}

func getASN(test string) string {
	ip, err := iptoasn.LookupIP(test)

	asn := strconv.FormatUint(uint64(ip.ASNum), 10)

	if err != nil {
		return "0"
	}
	return asn
}
