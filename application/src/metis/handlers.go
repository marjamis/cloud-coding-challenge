package main

import (
	"html/template"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

/*
Index Handler - Controls any call to / and will call functions to read/write to the database and use a template for the html response.
*/
func Index(w http.ResponseWriter, r *http.Request) {
	filename := "templates/index.html"

	t, err := template.ParseFiles(filename)
	if err != nil {
	  log.Error(err)
	}

	dbData, err := GetDBData()
	if err != nil {
		log.Error(err)
	}

	WriteData(r)

	data := &IndexResponse{InstanceId: INSTANCE_ID, Data: dbData}
	t.Execute(w, data)
}

/*
NotFound Handler - will return a 404 HTTP Code and a custom HTML page when the route is matched.
*/
func NotFound(w http.ResponseWriter, r *http.Request) {
	filename := "templates/http_404.html"

	w.WriteHeader(http.StatusNotFound)

	t, err := template.ParseFiles(filename)
	if err != nil {
		log.Error(err)
	}
	t.Execute(w, nil)
}

/*
HealthCheck Handler - Simple healthcheck which checks that the DB is still contactable and if not returns a 503 for the TG HealthCheck.
Also returns a light html page for easy review of the health.
*/
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	filename := "templates/healthcheck.html"
	state := false

	err := DB_CONNECTION.Ping()
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "HealthCheckDatabasePing",
		}).Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		state = true
	}

  data := &HealthCheckResponse{DBHealthy: state}
	t, err := template.ParseFiles(filename)
	if err != nil {
		log.Error(err)
	}
	t.Execute(w, data)
}
