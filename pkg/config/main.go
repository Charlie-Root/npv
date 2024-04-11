package config

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
    // Add fields for storing configuration options here.
    // For example:
    DbType string `json:"dbtype" yaml:"dbtype"`
    DbFilename string `json:"dbfilename" yaml:"dbfilename"`
    Api bool `json:"api" yaml:"api"`
    ApiPort int `json:"apiport" yaml:"apiport"`
    ApiServer string `json:"apiserver" yaml:"apiserver"`
    ApiSecure bool `json:"apisecure" yaml:"apisecure"`
    MTRView MTRView
}
type MTRView struct {
    Count int `json:"count" yaml:"count"`
    Timeout int `json:"timeout" yaml:"timeout"`
    Interval int `json:"interval" yaml:"interval"`
    MaxHops int `json:"maxhops" yaml:"maxhops"`
    MaxHopsUnknown int `json:"maxhops_unknown" yaml:"maxhops_unknown"`
    Ringbuffer int `json:"ringbuffer" yaml:"ringbuffer"`
    PtrLookup bool `json:"ptr_lookup" yaml:"ptr_lookup"`
    SRCAdresss string `json:"src_address" yaml:"src_address"`
} 

// LoadYAML reads a YAML configuration file and returns a Config struct
// with the options parsed from the file.
func LoadYAML(file string) (Config, error) {
    var config Config
    data, err := ioutil.ReadFile(file)
    if err != nil {
        return config, err
    }
    err = yaml.Unmarshal(data, &config)
    return config, err
}

// LoadJSON reads a JSON configuration file and returns a Config struct
// with the options parsed from the file.
func LoadJSON(file string) (Config, error) {
    var config Config
    data, err := ioutil.ReadFile(file)
    if err != nil {
        return config, err
    }
    err = json.Unmarshal(data, &config)
    return config, err
}
