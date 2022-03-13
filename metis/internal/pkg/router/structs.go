package router

import (
	"net/http"
	"time"
)

// Route contains the details for each route configured for metis
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is a list of the Routes of metis
type Routes []Route

// DBData contains all the information to create a database connection
type DBData struct {
	ID          int
	InstanceID  string
	RequesterIP string
	AddedDate   time.Time
}

// IndexResponse contains the details required for the index template
type IndexResponse struct {
	InstanceID string
	Data       DBData
}

// HealthCheckResponse contains the details of a database healthcheck
type HealthCheckResponse struct {
	DBHealthy bool
}
