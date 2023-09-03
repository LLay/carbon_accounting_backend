package energy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	resty "github.com/go-resty/resty/v2"

	"github.com/robfig/cron/v3"

	influxClientv2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	// iclient "github.com/influxdata/influxdb1-client/v2"
	//   "github.com/InfluxCommunity/influxdb3-go/influx"
)

// Total rows in EIA dataset: 17518447

// First ever row. Note lack of "value" field.
//  {
//     "period":"2018-07-02T06",
//     "respondent":"AEC",
//     "respondent-name":"PowerSouth Energy Cooperative",
//     "fueltype":"COL",
//     "type-name":"Coal"
//  },

func FetchEIAData() {
	c := cron.New()

	// Define the cron job to run the main function once an hour
	_, err := c.AddFunc("@hourly", func() {
		fmt.Println("Running FetchEIAData function...")

		// Your main function code here
		// getDataInDateRange(start, end)

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

// Response represents the JSON response structure.
type EIAPowerResponse struct {
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

// XParams represents the content of the X-Params header.
type XParams struct {
	Frequency string                 `json:"frequency"`
	Data      []string               `json:"data"`
	Facets    map[string]interface{} `json:"facets"`
	Start     string                 `json:"start"`
	End       interface{}            `json:"end"`
	Sort      []SortItem             `json:"sort"`
	Offset    int                    `json:"offset"`
	Length    int                    `json:"length"`
}

// SortItem represents an item in the Sort array.
type SortItem struct {
	Column    string `json:"column"`
	Direction string `json:"direction"`
}

// SerializeXParams converts an XParams struct to a JSON string.
func SerializeXParams(x XParams) (string, error) {
	jsonData, err := json.Marshal(x)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func ConvertXParamsToMap(x XParams) map[string]interface{} {
	xParamsMap := map[string]interface{}{
		"frequency": x.Frequency,
		"data":      x.Data,
		"facets":    x.Facets,
		"start":     x.Start,
		"end":       x.End,
		"sort":      x.Sort,
		"offset":    x.Offset,
		"length":    x.Length,
	}
	return xParamsMap
}

// function to parse a datetime into a string of the form "2023-09-02T00" ISO 8601
func parseDateTime(dt time.Time) string {
	return dt.Format("2006-01-02T15")
}

func parseDateTimeToTime(dt string) (time.Time, error) {
	return time.Parse("2006-01-02T15", dt)
}

func getDataInDateRange(start time.Time, end *time.Time) (err error) {
	xParams := XParams{
		Frequency: "hourly",
		Data:      []string{"value"},
		Facets:    map[string]interface{}{},
		Start:     parseDateTime(start),
		End:       nil,
		Sort: []SortItem{
			{"period", "desc"},
		},
		Offset: 0,
		Length: 5000,
	}

	if end != nil {
		xParams.End = parseDateTime(*end)
	}

	err = writePaginatedData(xParams)
	if err != nil {
		fmt.Printf("Error fetching EIA data: %v\n", err)
		return err
	}
	return nil
}

func writePaginatedData(xParams XParams) error {
	resp, err := FetchEIADataResty(xParams)
	if err != nil {
		return err
	}
	// write data to influx
	err = WriteDataToInfluxDB(resp)
	if err != nil {
		return err
	}

	if resp.Response.Total > resp.Request.Params.Offset+resp.Request.Params.Length {
		xParams.Offset = resp.Request.Params.Offset + resp.Request.Params.Length
		return writePaginatedData(xParams)
	}
	return nil
}

// https://www.eia.gov/opendata/browser/electricity/rto/fuel-type-data
// https://api.eia.gov/v2/electricity/rto/fuel-type-data/data/?api_key=CZdQsisRJzwOfqUWV3jiMPNEx3ZbHcuJ2VQus04i
func FetchEIADataResty(xParams XParams) (responseBody *EIAPowerResponse, err error) {
	fmt.Printf("Requesting EIA Power data with X-Params: %+v\n", xParams)

	client := resty.New()
	baseURL := "https://api.eia.gov/v2/electricity/rto/fuel-type-data/data/"
	payload := ConvertXParamsToMap(xParams)
	headers := map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Language":    "en-US,en;q=0.9",
		"Content-Type":       "application/json",
		"DNT":                "1",
		"Origin":             "https://www.eia.gov",
		"Referer":            "https://www.eia.gov/",
		"Sec-Ch-UA":          `"Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v="116"`,
		"Sec-Ch-UA-Mobile":   "?0",
		"Sec-Ch-UA-Platform": "macOS",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-site",
		"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	}

	resp, err := client.R().
		SetHeaders(headers).
		SetBody(payload).
		SetQueryParam("api_key", "CZdQsisRJzwOfqUWV3jiMPNEx3ZbHcuJ2VQus04i").
		Post(baseURL)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	responseBody, err = parseResponse(resp.Body())
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return nil, err
	}

	return responseBody, nil
}

func WriteDataToInfluxDB(data *EIAPowerResponse) (err error) {
	influxURL := "http://localhost:8086" // URL of your InfluxDB instance

	token := "admin"
	dBclient := influxClientv2.NewClient(influxURL, token)
	defer dBclient.Close()
	writeAPI := dBclient.WriteAPIBlocking("my_organization", "my_bucket")

	// iterate over data.response.data and create influxdb points
	var points []*write.Point
	for _, entry := range data.Response.Data {
		timestamp, err := time.Parse("2006-01-02T15", entry.Period)
		if err != nil {
			fmt.Printf("Error parsing time: %v\n", err)
			return err
		}
		p := influxClientv2.NewPointWithMeasurement("energy_data").
			AddTag("respondent_code", entry.Respondent).
			AddTag("respondent_name", entry.RespondentName).
			AddTag("fuel_type_code", entry.FuelType).
			AddTag("fuel_type_name", entry.TypeName).
			AddTag("value_units", entry.ValueUnits).
			AddField("value", entry.Value).
			SetTime(timestamp)

		points = append(points, p)
	}
	// bulk write points to influxdb
	err = writeAPI.WritePoint(context.Background(), points...)
	if err != nil {
		fmt.Printf("Error writing point: %v\n", err)
		return err
	}

	fmt.Println("Data written to InfluxDB successfully!")
	return nil
}

func parseResponse(responseBody []byte) (data *EIAPowerResponse, err error) {

	// Parse the JSON response into the struct
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return nil, err
	}

	return data, nil
}
