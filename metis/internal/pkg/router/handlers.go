package router

import (
	"html/template"
	"net/http"

	"github.com/marjamis/cloud-coding-challenge/metis/internal/pkg/database"
	"github.com/marjamis/cloud-coding-challenge/metis/internal/pkg/instance"

	log "github.com/sirupsen/logrus"
)

/*
Index Handler - Controls any call to / and will call functions to read/write to the database and use a template for the html response.
*/
func Index(w http.ResponseWriter, r *http.Request) {
	// TODO fix this pathing as it's a bit too crazy
	filename := "../../../configs/templates/index.html"

	t, err := template.ParseFiles(filename)
	if err != nil {
		log.Error(err)
	}

	dbData, err := GetDBData()
	if err != nil {
		log.Error(err)
	}

	WriteData(r)

	data := &IndexResponse{InstanceID: instance.InstanceID, Data: dbData}
	t.Execute(w, data)
}

/*
NotFound Handler - will return a 404 HTTP Code and a custom HTML page when the route is matched.
*/
func NotFound(w http.ResponseWriter, r *http.Request) {
	// TODO fix this pathing as it's a bit too crazy
	filename := "../../../configs/templates/http_404.html"

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
	// TODO fix this pathing as it's a bit too crazy
	filename := "../../../configs/templates/healthcheck.html"
	state := false

	err := database.DBConnection.Ping()
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
