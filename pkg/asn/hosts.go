package asn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/korylprince/ipnetgen"
)

// Hosts is a struct containing a list of IP addresses
type Hosts struct {
	Addresses []net.IP
}
// SubnetToHosts converts a subnet into a list of /32 hosts
func SubnetToHosts(subnet net.IPNet) (Hosts, error) {
	var hosts Hosts

	// Check if the subnet is an IPv4 subnet
	if subnet.IP.To4() == nil {
		return hosts, fmt.Errorf("not an IPv4 subnet: %v", subnet)
	}

	gen, err := ipnetgen.New(subnet.String())
	if err != nil {
		logger.Error(err.Error())
		//do something with err
	}
	for ip := gen.Next(); ip != nil; ip = gen.Next() {
		//do something with ip
		ipnetgen.Increment(ip)
		hosts.Addresses = append(hosts.Addresses, ip)
	}

	return hosts, nil
}


// isNetworkOrBroadcast returns true if the given IP is the network or broadcast address for the given subnet
func isNetworkOrBroadcast(ip net.IP, subnet net.IPNet) bool {
	if ip.Equal(subnet.IP) {
		return true
	}

	broadcast := subnet.IP.Mask(subnet.Mask)
	for i := range broadcast {
		broadcast[i] |= ^subnet.Mask[i]
	}

	return ip.Equal(broadcast)
}

// SaveHostsToJSON saves the given hosts to a JSON file with the following structure:
// { "hosts": [] }
func SaveHostsToJSON(hosts Hosts, filename string) error {
	// Convert the hosts to JSON
	data, err := json.Marshal(struct {
		Hosts []net.IP `json:"hosts"`
	}{
		Hosts: hosts.Addresses,
	})
	if err != nil {
		return err
	}

	// Write the JSON to a file
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	return nil
}
