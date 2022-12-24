package db

import (
	"database/sql"
	"strconv"

	"github.com/Charlie-Root/npv/pkg/logging"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/mattn/go-sqlite3"    // SQLite driver
)

var logger = logging.NewLogger("database")

// DB is a wrapper around sql.DB that provides a common interface for interacting
// with different SQL databases.
type DB struct {
	*sql.DB
}

type Host struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Address   string
	HostPTR   string
	HostASN   string
	HostTTL   int
	HostCount int
}

type Link struct {
	Source      string `json:"source"`
	Target      string `json:"target"`
	TargetTTL   int
	TargetLoss  int
	TargetSNT   int
	TargetLast  int
	TargetAVG   int
	TargetBest  int
	TargetWRST  int
	TargetStDev int
	LinkCount   int
}

// Open creates a new connection to the database.
//
// The driver parameter specifies the SQL database driver to use.
// The connectionString parameter specifies the database connection string.
func Open(driver, connectionString string) (*DB, error) {
	db, err := sql.Open(driver, connectionString)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) CreateTables() {
	createTableHosts := `CREATE TABLE IF NOT EXISTS hosts (
		"id" TEXT,
		"host_name" TEXT,
		"host_address" TEXT,
		"host_rdns" TEXT,
		"host_asn" TEXT,
		"host_ttl" INTEGER,
		"count" INTEGER
	  );`

	_, err := db.DB.Exec(createTableHosts)
	if err != nil {
		print(err.Error())
	}

	createTableLinks := `CREATE TABLE IF NOT EXISTS links (
		"hostid_from" TEXT,
		"hostid_to" TEXT,
		"target_ttl" INTEGER,
		"target_loss" INTEGER,
		"target_snt" INTEGER,
		"target_last" INTEGER,
		"target_avg" INTEGER,
		"target_best" INTEGER,
		"target_wrst" INTEGER,
		"target_stdev" INTEGER,
		"count" INTEGER
	  );`

	_, err = db.DB.Exec(createTableLinks)
	if err != nil {
		print(err.Error())
	}

}

func (db *DB) InsertHost(host Host) error {

	query, err := db.DB.Prepare("INSERT INTO hosts (id, host_name, host_address, host_rdns, host_asn, host_ttl, count) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		logger.Error(err.Error())
	}

	res, err := query.Exec(host.Address, host.Name, host.Address, host.HostPTR, host.HostASN, host.HostTTL, host.HostCount)
	if err != nil {
		logger.Error(err.Error())
	}

	_, err = res.LastInsertId()
	if err != nil {
		logger.Error(err.Error())
	}

	return nil
}

// Query executes the given query and returns the result.
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.DB.Query(query, args...)
}

func (db *DB) UpdateHostCount(host_address string, count int) {
	query, err := db.DB.Prepare("UPDATE hosts SET count = ? WHERE host_address = ?")

	if err != nil {
		logger.Error(err.Error())
	}
	var newCount = (count + 1)
	_, err = query.Exec(newCount, host_address)

	if err != nil {
		logger.Error(err.Error())
	}

	//RowsAffected, err := res.RowsAffected()

	//fmt.Println(RowsAffected)

	if err != nil {
		logger.Error(err.Error())
	}

}

func (db *DB) QueryLink(hostid_from int, hostid_to int, target_ttl int) int {

	rows, err := db.Query("SELECT count FROM links WHERE hostid_from = ? AND hostid_to = ? AND target_ttl = ?")
	if err != nil {
		logger.Error(err.Error())
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		count++
	}

	// Check for errors
	if err := rows.Err(); err != nil {
		return 0
	}
	return count

}

func (db *DB) QueryHost(host_address string) int {
	rows, err := db.Query("SELECT count FROM HOSTS where host_address = '" + host_address + "'")
	if err != nil {
		logger.Error(err.Error())
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		count++
	}

	// Check for errors
	if err := rows.Err(); err != nil {
		return 0
	}
	return count

}

// QueryRow executes a SQL query and returns the first row of the result.
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRow(query, args...)
}

// Exec executes a SQL query that does not return any rows.
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.DB.Exec(query, args...)
}

// Exec executes a SQL query that does not return any rows.
func (db *DB) GetHostByAddress(host_address string, host_ttl int) (string, string, int) {

	selectStatement := `
		SELECT *
		FROM hosts
		WHERE host_address = ? AND host_ttl = ?
		`

	var host Host

	err := db.QueryRow(selectStatement, host_address, host_ttl).Scan(&host.Id, &host.Name, &host.Address, &host.HostPTR, &host.HostASN, &host.HostTTL, &host.HostCount)
	if err != nil {
		return "", "", 0
	}

	return host.Id, host.Address, host.HostCount

}

func (db *DB) GetLink(hostid_from string, hostid_to string, target_ttl int) (string, string, int, int) {

	selectStatement := `
		SELECT *
		FROM links
		WHERE hostid_from = ? AND hostid_to = ? AND target_ttl = ?
		`

	var link Link

	//logger.Debug(hostid_from)
	//logger.Debug(hostid_to)
	//logger.Debug(strconv.Itoa(target_ttl))

	err := db.QueryRow(selectStatement, hostid_from, hostid_to, strconv.Itoa(target_ttl)).Scan(&link.Source, &link.Target, &link.TargetTTL, &link.TargetLoss, &link.TargetSNT, &link.TargetLast, &link.TargetAVG, &link.TargetBest, &link.TargetWRST, &link.TargetStDev, &link.LinkCount)
	if err != nil {
		//logger.Debug(err.Error())
		return "", "", 0, 0
	}

	return link.Source, link.Target, link.TargetTTL, link.LinkCount

}

func (db *DB) InsertLink(link Link) error {
	query, err := db.DB.Prepare("INSERT INTO links (hostid_from, hostid_to, target_ttl, target_loss, target_snt, target_last, target_avg, target_best, target_wrst, target_stdev, count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		logger.Error(err.Error())
	}

	_, err = query.Exec(link.Source, link.Target, link.TargetTTL, link.TargetLoss, link.TargetSNT, link.TargetLast, link.TargetAVG, link.TargetBest, link.TargetWRST, link.TargetStDev, link.LinkCount)
	if err != nil {
		logger.Error(err.Error())
	}

	return nil
}

func (db *DB) UpdateLinkCount(hostid_from string, hostid_to string, target_ttl int, count int) {
	query, err := db.DB.Prepare("UPDATE links SET count = ? WHERE hostid_from = ? AND hostid_to = ? AND target_ttl = ? ")

	if err != nil {
		logger.Error(err.Error())
	}
	var newCount = (count + 1)
	res, err := query.Exec(newCount, hostid_from, hostid_to, target_ttl)

	if err != nil {
		logger.Error(err.Error())
	}

	_, err = res.RowsAffected()

	if err != nil {
		logger.Error(err.Error())
	}

}
