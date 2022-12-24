package export

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Charlie-Root/npv/pkg/db"
	"github.com/Charlie-Root/npv/pkg/logging"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

var logger = logging.NewLogger("export")

type Graph struct {
	Nodes *[]db.Host `json:"nodes"`
	Links *[]db.Link `json:"links"`
}

func ParseHostsTable(database *db.DB) []db.Host {
	rows, err := database.Query("SELECT * FROM hosts")
	if err != nil {
		logger.Error("Something went wrong")
	}
	defer rows.Close()

	// Create a slice to hold the rows
	var hosts []db.Host
	for rows.Next() {
		var h db.Host
		err = rows.Scan(&h.Id, &h.Name, &h.Address, &h.HostPTR, &h.HostASN, &h.HostTTL, &h.HostCount)
		if err != nil {
			fmt.Println(err)
		}
		hosts = append(hosts, h)
	}
	return hosts
}

func ParseGraphData(database *db.DB) ([]byte, error) {
	rows, err := database.Query("SELECT * FROM hosts")
	if err != nil {
		logger.Error("Something went wrong")
	}
	defer rows.Close()

	nodes := make(map[string]*db.Host)
	links := make([]db.Link, 0)

	for rows.Next() {

		var h db.Host
		err = rows.Scan(&h.Id, &h.Name, &h.Address, &h.HostPTR, &h.HostASN, &h.HostTTL, &h.HostCount)
		if err != nil {
			fmt.Println(err)
		}
		//logger.Debug("Adding host " + h.Id + " to tree with TTL: " + strconv.Itoa(h.HostTTL))
		nodes[h.Id+"_"+strconv.Itoa(h.HostTTL)] = &db.Host{Id: h.Id + "_" + strconv.Itoa(h.HostTTL), Name: h.Name, Address: h.Address, HostPTR: h.HostPTR, HostASN: h.HostASN, HostTTL: h.HostTTL, HostCount: h.HostCount}

	}

	// start links
	rows2, err := database.Query("SELECT * FROM links")
	if err != nil {
		logger.Error("Something went wrong")
	}
	defer rows.Close()

	for rows2.Next() {
		var l db.Link
		err = rows2.Scan(&l.Source, &l.Target, &l.TargetTTL, &l.TargetLoss, &l.TargetSNT, &l.TargetLast, &l.TargetAVG, &l.TargetBest, &l.TargetWRST, &l.TargetStDev, &l.LinkCount)
		if err != nil {
			fmt.Println(err)
		}

		links = append(links, db.Link{Source: l.Source + "_" + strconv.Itoa(l.TargetTTL-1), Target: l.Target + "_" + strconv.Itoa(l.TargetTTL), TargetTTL: l.TargetTTL, TargetLoss: l.TargetLoss, TargetSNT: l.TargetSNT, TargetLast: l.TargetLast, TargetAVG: l.TargetAVG, TargetBest: l.TargetBest, TargetWRST: l.TargetWRST, TargetStDev: l.TargetStDev, LinkCount: l.LinkCount})
	}
	data, err := json.MarshalIndent(Graph{Nodes: values(nodes), Links: &links}, "", "	")

	if err != nil {
		return nil, fmt.Errorf("JSON marshaling failed: %s", err)
	}
	return data, nil
}
func values(nodes map[string]*db.Host) *[]db.Host {
	array := []db.Host{}
	for _, n := range nodes {
		array = append(array, *n)
	}
	return &array
}
