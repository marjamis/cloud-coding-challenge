package router

var routes = Routes{
	// HealthCheck Route
	Route{
		Name:        "HealthCheck",
		Method:      "GET",
		Pattern:     "/healthcheck",
		HandlerFunc: HealthCheck,
	},
	// Index Route
	Route{
		Name:        "Index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: Index,
	},
}
