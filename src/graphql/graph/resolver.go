package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewResolver() *Resolver {
	r := Resolver{}
	influxURL := "http://localhost:8086" // URL of your InfluxDB instance

	token := "admin"
	r.influxClient = influxdb2.NewClient(influxURL, token)
	return &r
}

type Resolver struct {
	influxClient influxdb2.Client
}
