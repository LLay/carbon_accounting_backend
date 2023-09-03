package energy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	influxClientv2 "github.com/influxdata/influxdb-client-go/v2"
	// iclient "github.com/influxdata/influxdb1-client/v2"
	//   "github.com/InfluxCommunity/influxdb3-go/influx"
)

// Response represents the JSON response structure.
type Response struct {
	Response struct {
		Total      int    `json:"total"`
		DateFormat string `json:"dateFormat"`
		Frequency  string `json:"frequency"`
		Data       []struct {
			Period         string `json:"period"`
			Respondent     string `json:"respondent"`
			RespondentName string `json:"respondent-name"`
			FuelType       string `json:"fueltype"`
			TypeName       string `json:"type-name"`
			Value          int    `json:"value"`
			ValueUnits     string `json:"value-units"`
		} `json:"data"`
		Description string `json:"description"`
	} `json:"response"`
	Request struct {
		Command string `json:"command"`
		Params  struct {
			Frequency string      `json:"frequency"`
			Data      []string    `json:"data"`
			Facets    []string    `json:"facets"`
			Start     string      `json:"start"`
			End       interface{} `json:"end"`
			Sort      []struct {
				Column    string `json:"column"`
				Direction string `json:"direction"`
			} `json:"sort"`
			Offset int `json:"offset"`
			Length int `json:"length"`
		} `json:"params"`
	} `json:"request"`
	APIVersion string `json:"apiVersion"`
}

func main() {
	c := cron.New()

	// Define the cron job to run the main function once an hour
	_, err := c.AddFunc("@hourly", func() {
		fmt.Println("Running FetchEIAData function...")

		// Your main function code here
		FetchEIAData()

		fmt.Println("FetchEIAData function completed.")
	})
	if err != nil {
		fmt.Printf("Error scheduling cron job: %v\n", err)
		return
	}

	// Start the cron scheduler
	c.Start()

	// Keep the program running until interrupted
	select {}
}

func FetchEIAData() {
	cmd := exec.Command("/bin/bash", "request.sh")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running script: %v\n", err)
		return
	}

	// Extract the response file path from the script output
	responseFilePath := "tmp.json"
	// responseFilePath := strings.TrimSpace(string(output))

	// Read the response from the file
	responseBody, err := os.ReadFile(responseFilePath)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Now 'response' contains the response body
	fmt.Printf("Response: %s\n", responseBody)

	data, err := parseResponse(responseBody)
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}

	// Initialize InfluxDB client
	// dBclient := influxClientv2.NewClient("http://localhost:8086", "your-token")
	// InfluxDB configuration
	influxURL := "http://localhost:8086" // URL of your InfluxDB instance

	token := "admin"
	dBclient := influxClientv2.NewClient(influxURL, token)

	if err != nil {
		log.Fatal(err)
	}
	defer dBclient.Close()

	// Create a write API
	writeAPI := dBclient.WriteAPIBlocking("your-org", "your-bucket")

	// Iterate over the data and write it to InfluxDB
	for _, entry := range data.Response.Data {
		// Create a point
		p := influxClientv2.NewPointWithMeasurement("energy_data").
			AddTag("respondent", entry.Respondent).
			AddTag("fuel_type", entry.FuelType).
			AddField("value", entry.Value).
			SetTime(time.Now())

		// Write the point to InfluxDB
		writeAPI.WritePoint(context.Background(), p)
		// p := influxdb2.NewPoint("stat",
		//     map[string]string{"unit": "temperature"},
		//     map[string]interface{}{"avg": 24.5, "max": 45.0},
		//     time.Now())
		// // write point immediately
		// writeAPI.WritePoint(context.Background(), p)
	}

	// Close the write API and dBclient
	// writeAPI.Close()
	dBclient.Close()
	fmt.Println("Data written to InfluxDB successfully!")
}

func FetchEIADataBustedHeaders() {
	url := "https://api.eia.gov/v2/electricity/rto/fuel-type-data/data/?api_key=CZdQsisRJzwOfqUWV3jiMPNEx3ZbHcuJ2VQus04i"

	// Create an HTTP client
	client := &http.Client{}

	// Create an HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	fmt.Println("Fetching data from EIA API...")
	log.Println("Fetching data from EIA API...")

	// Set headers as needed
	req.Header.Set("authority", "api.eia.gov")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("dnt", "1")
	req.Header.Set("origin", "https://www.eia.gov")
	req.Header.Set("referer", "https://www.eia.gov/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v="116"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	req.Header.Set("x-params", "{\"frequency\":\"hourly\",\"data\":[\"value\"],\"facets\":{},\"start\":\"2023-09-02T00\",\"end\":null,\"sort\":[{\"column\":\"period\",\"direction\":\"desc\"}],\"offset\":0,\"length\":5")

	// Convert all header keys to lowercase
	lowerCaseHeader := make(http.Header)
	for key, value := range req.Header {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
		lowerCaseHeader[strings.ToLower(key)] = value
	}

	// Need headers to NOT BE SLICES
	// headers := map[string]string{
	// 	"authority":          "api.eia.gov",
	// 	"accept":             "application/json, text/plain, */*",
	// 	"accept-language":    "en-US,en;q=0.9",
	// 	"content-type":       "application/json",
	// 	"dnt":                "1",
	// 	"origin":             "https://www.eia.gov",
	// 	"referer":            "https://www.eia.gov/",
	// 	"sec-ch-ua":          `"Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v=", Value:16"`,
	// 	"sec-ch-ua-mobile":   "?0",
	// 	"sec-ch-ua-platform": `"macOS"`,
	// 	"sec-fetch-dest":     "empty",
	// 	"sec-fetch-mode":     "cors",
	// 	"sec-fetch-site":     "same-site",
	// 	"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	// 	"x-params":           `{"frequency":"hourly","data":["value"],"facets":{},"start":"2023-09-02T00","end":null,"sort":[{"column":"period","direction":"desc"}],"offset":0,"length":5`,
	// }

	fmt.Print("lowerCaseHeader")
	// req.Header = headers
	fmt.Printf("Request: %v\n", req)
	return
	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}
	// fmt.Printf("Response body: %s\n", responseBody)

	data, err := parseResponse(responseBody)
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}
	// Now 'data' contains the parsed response data
	fmt.Printf("Total: %d\n", data.Response.Total)
	fmt.Printf("Frequency: %s\n", data.Response.Frequency)
	fmt.Printf("Description: %s\n", data.Response.Description)

	// Access other fields in the 'data' struct as needed
	fmt.Printf("Command: %s\n", data.Request.Command)
	fmt.Printf("Start: %s\n", data.Request.Params.Start)

	// https://github.com/influxdata/influxdb-client-go
	// Create a new client using an InfluxDB server base URL and an authentication token
	// client := influxdb2.NewClient("http://localhost:8086", "my-token")
	// // Use blocking write client for writes to desired bucket
	// writeAPI := client.WriteAPIBlocking("my-org", "my-bucket")
	// // Create point using full params constructor
	// p := influxdb2.NewPoint("stat",
	//     map[string]string{"unit": "temperature"},
	//     map[string]interface{}{"avg": 24.5, "max": 45.0},
	//     time.Now())
	// // write point immediately
	// writeAPI.WritePoint(context.Background(), p)
	// // Create point using fluent style
	// p = influxdb2.NewPointWithMeasurement("stat").
	//     AddTag("unit", "temperature").
	//     AddField("avg", 23.2).
	//     AddField("max", 45.0).
	//     SetTime(time.Now())
	// err := writeAPI.WritePoint(context.Background(), p)
	// if err != nil {
	// 	panic(err)
	// }
	// // Or write directly line protocol
	// line := fmt.Sprintf("stat,unit=temperature avg=%f,max=%f", 23.5, 45.0)
	// err = writeAPI.WriteRecord(context.Background(), line)
	// if err != nil {
	// 	panic(err)
	// }

	// Initialize InfluxDB client
	// dBclient := influxClientv2.NewClient("http://localhost:8086", "your-token")
	// InfluxDB configuration
	influxURL := "http://localhost:8086" // URL of your InfluxDB instance
	// influxDB := "mydb"                   // Name of your database
	// username := "admin"                  // Admin username
	// password := "admin_password"         // Admin password

	// Create a new HTTP client
	token := "admin"
	dBclient := influxClientv2.NewClient(influxURL, token)
	// NewHTTPClient(influxClientv2.HTTPConfig{
	// 	Addr:     influxURL,
	// 	Username: username,
	// 	Password: password,
	// })
	if err != nil {
		log.Fatal(err)
	}
	defer dBclient.Close()

	// Create a write API
	writeAPI := dBclient.WriteAPIBlocking("your-org", "your-bucket")

	// Iterate over the data and write it to InfluxDB
	for _, entry := range data.Response.Data {
		// Create a point
		p := influxClientv2.NewPointWithMeasurement("energy_data").
			AddTag("respondent", entry.Respondent).
			AddTag("fuel_type", entry.FuelType).
			AddField("value", entry.Value).
			SetTime(time.Now())

		// Write the point to InfluxDB
		writeAPI.WritePoint(context.Background(), p)
		// p := influxdb2.NewPoint("stat",
		//     map[string]string{"unit": "temperature"},
		//     map[string]interface{}{"avg": 24.5, "max": 45.0},
		//     time.Now())
		// // write point immediately
		// writeAPI.WritePoint(context.Background(), p)
	}

	// Close the write API and dBclient
	// writeAPI.Close()
	dBclient.Close()
	fmt.Println("Data written to InfluxDB successfully!")
}

// func writeToInfluxDB()
func parseResponse(responseBody []byte) (data *Response, err error) {

	// Parse the JSON response into the struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return nil, err
	}

	return data, nil
}
