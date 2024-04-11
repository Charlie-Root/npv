package apiclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Charlie-Root/npv/pkg/config"
)

// Host represents a host in the API.
type Host struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Address   string `json:"address"`
    HostPTR   string `json:"hostPTR"`
    HostASN   string `json:"hostASN"`
    HostTTL   int    `json:"hostTTL"`
}

// Link represents a link in the API.
type Link struct {
    Source      string `json:"source"`
    Target      string `json:"target"`
    TargetTTL   int    `json:"targetTTL"`
    TargetLoss  int    `json:"targetLoss"`
    TargetSNT   int    `json:"targetSNT"`
    TargetLast  int    `json:"targetLast"`
    TargetAVG   int    `json:"targetAVG"`
    TargetBest  int    `json:"targetBest"`
    TargetWRST  int    `json:"targetWRST"`
    TargetStDev int    `json:"targetStDev"`
}

// APIClient represents the API client.
type APIClient struct {
    baseURL string
}

// NewAPIClient creates a new instance of APIClient.
func NewAPIClient(config config.Config) *APIClient {
	port := strconv.Itoa(config.ApiPort)
    
    // Construct the base URL from the configuration
    baseURL := "http://" + config.ApiServer + ":" + port


    return &APIClient{baseURL: baseURL}
}

// AddHost sends a POST request to add a host to the API.
func (c *APIClient) AddHost(host Host) error {
    url := c.baseURL + "/api/add-host"
    data, err := json.Marshal(host)
    if err != nil {
        return err
    }
    _, err = http.Post(url, "application/json", bytes.NewBuffer(data))
    return err
}

// AddLink sends a POST request to add a link to the API.
func (c *APIClient) AddLink(link Link) error {
    url := c.baseURL + "/api/add-link"
    data, err := json.Marshal(link)
    if err != nil {
        return err
    }
    _, err = http.Post(url, "application/json", bytes.NewBuffer(data))
    return err
}
