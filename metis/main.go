package main

import (
	"net/http"

	"github.com/marjamis/cloud-coding-challenge/metis/internal/pkg/database"
	"github.com/marjamis/cloud-coding-challenge/metis/internal/pkg/router"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := router.NewRouter()
	defer database.DBConnection.Close()

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
