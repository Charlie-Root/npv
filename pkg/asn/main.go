package asn

import (
	"log"
	"net"

	"github.com/Charlie-Root/npv/pkg/logging"
)
var logger = logging.NewLogger("asn")

func GenerateFile(asn int) {
	// Create a new asn client
	client := NewClient()

	// Get the IPv4 subnets for ASN 1234
	subnets, err := client.SubnetsForASN(asn)
	if err != nil {
		logger.Error(err.Error())
	}

	// Convert the subnets into hosts
	var hosts Hosts
	for _, subnet := range subnets {
		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			log.Fatal(err)
		}
		hostsForSubnet, err := SubnetToHosts(*ipNet)
		if err != nil {
			log.Fatal(err)
		}
		hosts.Addresses = append(hosts.Addresses, hostsForSubnet.Addresses...)
	}

	// Save the hosts to a JSON file
	if err := SaveHostsToJSON(hosts, "hosts.json"); err != nil {
		log.Fatal(err)
	}
}
