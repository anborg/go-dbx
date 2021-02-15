package drivermgr

import (
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb" //MSSQL
	_ "github.com/godror/godror"         //Oracle
	"github.com/jmoiron/sqlx"            //Extend sql
	_ "github.com/lib/pq"                //Postgres
	_ "github.com/mattn/go-sqlite3"      //Sqlite3
)

//DataSource definition
type Datasource struct {
	URL    string //TODO convert to host, port, user, pwd - do not log user/pwd
	DbType string
}

//GetConnection for given datasource
func GetConnection(ds Datasource) (*sqlx.DB, error) {
	defaultMaxOpenConn := 5
	defaultMaxIdleConn := 1
	defaultMaxConnLifetime := time.Minute * 10

	var err error
	var db *sqlx.DB
	log.Println("Using driver : " + ds.DbType)
	log.Println("Atteption DB connection for url: " + ds.URL)

	switch ds.DbType {
	case "oracle":
		db, err = sqlx.Open(Dbdriver.Oracle, ds.URL)
	case "postgres":
		db, err = sqlx.Open(Dbdriver.Postges, ds.URL)
	case "sqlite":
		db, err = sqlx.Open(Dbdriver.Sqlite, ds.URL) // user ":memory:" for inmem
	case "mysql":
		db, err = sqlx.Open(Dbdriver.Mysql, ds.URL)
	case "sqlserver":
		db, err = sqlx.Open(Dbdriver.Mssql, ds.URL)
	case "db2":
		db, err = sqlx.Open(Dbdriver.Db2, ds.URL)

	}
	//Add default option
	db.SetMaxOpenConns(defaultMaxOpenConn)
	db.SetMaxIdleConns(defaultMaxIdleConn)
	db.SetConnMaxLifetime(defaultMaxConnLifetime)
	db.Ping()
	log.Println("DB connection got for url: " + ds.URL)
	return db, err
}

var Dbdriver = newDbDriverRegistry()

func newDbDriverRegistry() *dbDriverRegistry {
	return &dbDriverRegistry{
		Oracle:  "godror",
		Postges: "postgres",
		Sqlite:  "sqlite3",
		Mssql:   "mssql",
		Mysql:   "mssql",
		Db2:     "db2",
	}
}

type dbDriverRegistry struct {
	Oracle  string
	Postges string
	Sqlite  string
	Mssql   string
	Mysql   string
	Db2     string
}
