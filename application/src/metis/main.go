package main

import (
	"database/sql"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

/*
Some global Variables
*/
var (
	DB_CONNECTION *sql.DB
	INSTANCE_ID   string
)

const (
	DB_MAX_IDLE_CONNS int = 0
	DB_MAX_OPEN_CONNS int = 0
)

/*
Starts the requirements for the application, such as get the instanceId and open a DB connection.
*/
func init() {
	INSTANCE_ID = GetInstanceId()

	// Opens up the DB connection
	var err error
	DB_CONNECTION, err = sql.Open("mysql",
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
	err = DB_CONNECTION.Ping()
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "DatabaseConnection",
		}).Fatal(err)
	}

	DB_CONNECTION.SetMaxIdleConns(DB_MAX_IDLE_CONNS)
	DB_CONNECTION.SetMaxOpenConns(DB_MAX_OPEN_CONNS)
}

func main() {
	router := NewRouter()
	defer DB_CONNECTION.Close()

	//Start the HTTP webserver listening
	log.WithFields(log.Fields{
		"Function": "main",
	}).Info("HTTP Listening on port 8081...")
	err := http.ListenAndServe(":8081", router)
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "main",
		}).Error(err)
		return
	}
}
