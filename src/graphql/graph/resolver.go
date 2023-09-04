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

// func (r *Resolver) EnergyMeasurements(ctx context.Context) ([]*model.EnergyMeasurement, error) {
// 	// Implement the resolver to fetch data from InfluxDB
// 	// You'll need to set up your InfluxDB client and query the database.
// 	// Return the fetched data as an array of EnergyData objects.
// 	return []*model.EnergyMeasurement{
// 		{
// 			Value: 1234,
// 		},
// 	}, nil
// }

// func (r *Resolver) GetAllMeasurements(ctx context.Context) ([]*model.EnergyMeasurement, error) {

//     measurements := []*model.EnergyMeasurement{
//         {
// 			Value: 1234,
//             // Period:        "2018-07-02T06",
//             // Respondent:    "AEC",
//             // RespondentName: "PowerSouth Energy Cooperative",
//             // FuelType:      "COL",
//             // TypeName:      "Coal",
//         },
//         // Add more Measurement objects as needed
//     }
//     return measurements, nil
// }
