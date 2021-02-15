package drivermgr

import (
	"database/sql"
	"log"

	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
)

//DataSource definition
type DataSource struct {
	URL        string
	DbType     string
	DriverName string
}

//GetConnection for given datasource
func GetConnection(ds DataSource) (*sql.DB, error) {
	oracleDriverName := "godror"
	postgresDriverName := "postgres"
	sqlite3DriverName := "godror"
	mssqlDriverName := "godror"
	db2DriverName := "godror"
	var err error
	var db *sql.DB
	switch ds.DbType {
	case "oracle":
		log.Println("Oracle connection for url: " + ds.URL)
		db, err = sql.Open(oracleDriverName, ds.URL)
		break
	case "postgres":
		db, err = sql.Open(postgresDriverName, ds.URL)
	case "sqlite":
		db, err = sql.Open(sqlite3DriverName, ds.URL)
	case "mysql":
		db, err = sql.Open(sqlite3DriverName, ds.URL)
	case "mssql":
		db, err = sql.Open(mssqlDriverName, ds.URL)
	case "db2":
		db, err = sql.Open(db2DriverName, ds.URL)

	}
	db.Ping()
	return db, err
}
