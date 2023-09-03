package energy

import (
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

// func writeToInfluxDB(ctx context.Context, influxClient client.Client, bucket string, measurement string, tags map[string]string, fields map[string]interface{}) error {
// 	// Now time for metrics
// 	now := time.Now()

// 	logAction := fmt.Sprintf("writing to influxdb")
// 	logger.Info(ctx, logAction, "", nil)

// 	writeAPI := influxClient.WriteAPIBlocking("", bucket)
// 	p := client.NewPoint(measurement, tags, fields, time.Now())
// 	if err := writeAPI.WritePoint(ctx, p); err != nil {
// 		logger.Error(ctx, logAction, err.Error(), nil)
// 		return err
// 	}

// 	// Send metrics to Prometheus
// 	metrics.InfluxDBDurationsSumary.WithLabelValues("Write").Observe(time.Since(now).Seconds())
// 	metrics.InfluxDBDurationsHistogram.WithLabelValues("Write").Observe(time.Since(now).Seconds())

// 	logger.Info(ctx, logAction, "object inserted with success", nil)
// 	return nil
// }

func writeToInfluxDB2() {
	// InfluxDB configuration
	influxURL := "http://localhost:8086" // URL of your InfluxDB instance
	influxDB := "mydb"                   // Name of your database
	username := "admin"                  // Admin username
	password := "admin_password"         // Admin password

	// Create a new HTTP client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxURL,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Create a new batch of points
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxDB,
		Precision: "s", // Use seconds precision
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a data point
	tags := map[string]string{"location": "office"}
	fields := map[string]interface{}{"temperature": 25.5}
	pt, err := client.NewPoint("sensor_data", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	// Add the data point to the batch
	bp.AddPoint(pt)

	// Write the batch to InfluxDB
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data written to InfluxDB successfully!")
}
