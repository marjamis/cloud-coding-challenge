package main

import (
	"net/http"
	"time"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type DBData struct {
	Id        int
	InstanceId      string
	RequesterIp     string
	AddedDate time.Time
}

type IndexResponse struct {
	InstanceId string
	Data       DBData
}

type HealthCheckResponse struct {
	DBHealthy bool
}
