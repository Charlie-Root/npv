package parser

import (
	"io/ioutil"
	"net"
	"sort"
	"strconv"

	"github.com/Charlie-Root/npv/pkg/apiclient"
	"github.com/Charlie-Root/npv/pkg/config"
	"github.com/Charlie-Root/npv/pkg/db"
	"github.com/Charlie-Root/npv/pkg/logging"
	"github.com/Charlie-Root/npv/pkg/mtr"
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
    var c, _ = config.LoadYAML("config.yaml")

    // Initialize the API client if the API is enabled
    var hostClient *apiclient.APIClient
    var linkClient *apiclient.APIClient

    if c.Api {
        logger.Warning("API is enabled")
        hostClient = apiclient.NewAPIClient(c)
        linkClient = apiclient.NewAPIClient(c)
    }

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

            _, hostAddress, hostCount := database.GetHostByAddress(searchHostname, m.Statistic[key].TTL)
            if hostAddress != "" {
                if c.Api {
                    logger.Warning("Writing to Remote API")
                    // Call apiclient.AddHost with the appropriate data
                    host := apiclient.Host{
                        ID:        "example",
                        Name:      searchHostname,
                        Address:   hostAddress,
                        HostPTR:   getPTR(searchHostname),
                        HostASN:   getASN(searchHostname),
                        HostTTL:   m.Statistic[key].TTL,
                    }
                    // AddHost returns an error, handle it if necessary
                    err := hostClient.AddHost(host)
                    if err != nil {
                        logger.Warning("Something went wrong")
                    }
                } else {
                    logger.Warning("Writing to local DB")
                    database.UpdateHostCount(searchHostname, hostCount)
                }
            } else {
                // Create a new host entry
                host := db.Host{
                    Name:      searchHostname,
                    Address:   searchHostname,
                    HostPTR:   getPTR(searchHostname),
                    HostASN:   getASN(searchHostname),
                    HostTTL:   m.Statistic[key].TTL,
                    HostCount: 1,
                }
                // Insert the host into the database
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

            if curAddress == "" {
                curAddress = "***"
            }
            if prevAddress == "" {
                prevAddress = "***"
            }

            hostFrom, hostTo, _, linkCounter := database.GetLink(prevAddress, curAddress, curTTL)

            if hostFrom == "" && hostTo == "" {
                if c.Api {
                    logger.Warning("Writing link to Remote API")
                    // Call apiclient.AddLink with the appropriate data
                    link := apiclient.Link{
                        Source:      prevAddress,
                        Target:      curAddress,
                        TargetTTL:   curTTL,
                        // Add other fields as needed
                    }
                    // AddLink returns an error, handle it if necessary
                    err := linkClient.AddLink(link)
                    if err != nil {
                        logger.Warning("Something went wrong")
                    }
                } else {
                    logger.Warning("Writing link to local DB")
                    // Add the link
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
                    // Insert the link into the database
                    if err := database.InsertLink(link); err != nil {
                        logger.Error(err.Error())
                        return
                    }
                }
            } else {
                database.UpdateLinkCount(prevAddress, curAddress, curTTL, linkCounter)
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
