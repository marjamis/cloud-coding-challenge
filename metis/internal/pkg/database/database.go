package database

import (
	"database/sql"
	"os"

	// Used to connect to mysql databases
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// DBConnection stores all the details to connect to the datbase backend
var DBConnection *sql.DB

const (
	// DBMaxIdleConnections sets how many idle connections metis can have to the database backend
	DBMaxIdleConnections int = 0
	// DBMaxOpenConnections sets how many open connections metis can have to the database backend
	DBMaxOpenConnections int = 0
)

// CreateDBConnection creates a connection to the database and sets any required configurations
func CreateDBConnection() {
	DBConnection, err := sql.Open("mysql",
		os.Getenv("UserName")+":"+
			os.Getenv("Password")+"@tcp("+
			os.Getenv("Endpoint")+":"+
			os.Getenv("Port")+")/"+
			os.Getenv("Name")+"?parseTime=true")
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "DatabaseConnection",
		}).Fatal(err)
	}

	// Ping the DB just to make sure it's contactable
	err = DBConnection.Ping()
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "DatabaseConnection",
		}).Fatal(err)
	}

	DBConnection.SetMaxIdleConns(DBMaxIdleConnections)
	DBConnection.SetMaxOpenConns(DBMaxOpenConnections)
}
