package asn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const bgpviewAPI = "https://api.bgpview.io/"

// Client is a client for the BGPView API
type Client struct {
	
}

// NewClient returns a new BGPView client with the given API key
func NewClient() *Client {
	return &Client{
	}
}

// SubnetsForASN returns all IPv4 subnets for the given ASN
func (c *Client) SubnetsForASN(asn int) ([]string, error) {
	// Build the request URL
	u, err := url.Parse(bgpviewAPI)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	
	u.Path = fmt.Sprintf("asn/%d/prefixes", asn)
	
	u.RawQuery = q.Encode()

	// Send the request
	res, err := http.Get(u.String())
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	defer res.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// Parse the response
	var response struct {
		Data struct {
			IPv4Prefixes []struct {
				Prefix string `json:"prefix"`
			} `json:"ipv4_prefixes"`
		} `json:"data"`
	}

	
	if err := json.Unmarshal(body, &response); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// Extract the subnets from the response
	var subnets []string
	for _, prefix := range response.Data.IPv4Prefixes {
		subnets = append(subnets, prefix.Prefix)
	}

	return subnets, nil
}
