package energy

import (
	"testing"
)

func TestFetchEIADataResty(t *testing.T) {
	// FetchEIADataResty()
}

// func TestFetchEIAData(t *testing.T) {
// 	FetchEIAData()
// }

func TestGetDataInDateRange(t *testing.T) {
	start, err := ParseDateTimeToTime("2023-08-31T00")
	end, err := ParseDateTimeToTime("2023-09-03T00")
	if err != nil {
		t.Errorf("Error parsing date: %v", err)
	}
	err = getDataInDateRange(start, &end)
	if err != nil {
		t.Errorf("Error getting data: %v", err)
	}
}

// func TestCreateOrg(t *testing.T) {
// 	createOrg()
// }

// func TestCreateBucket(t *testing.T) {
// 	createBucket()
// }

func TestUnmarshalResponse(t *testing.T) {
	// JSON data to unmarshal
	jsonData := []byte(`{
		"response": {
			"total": 2496,
			"dateFormat": "YYYY-MM-DD\"T\"HH24",
			"frequency": "hourly",
			"data": [
				{
					"period": "2023-09-02T07",
					"respondent": "AVA",
					"respondent-name": "Avista Corporation",
					"fueltype": "WND",
					"type-name": "Wind",
					"value": 0,
					"value-units": "megawatthours"
				}
			],
			"description": "Hourly net generation by balancing authority and energy source."
		},
		"request": {
			"command": "/v2/electricity/rto/fuel-type-data/data/",
			"params": {
				"frequency": "hourly",
				"data": [
					"value"
				],
				"facets": [],
				"start": "2023-09-02T00",
				"end": null,
				"sort": [
					{
						"column": "period",
						"direction": "desc"
					}
				],
				"offset": 0,
				"length": 5000
			}
		},
		"apiVersion": "2.1.4"
	}`)

	parsedData, err := parseResponse(jsonData)
	if err != nil {
		t.Errorf("Error parsing JSON data: %v", err)
	}

	// Verify the values in the parsed struct
	if parsedData.Response.Total != 2496 {
		t.Errorf("Expected Total to be 2496, got %d", parsedData.Response.Total)
	}

	if parsedData.Response.Frequency != "hourly" {
		t.Errorf("Expected Frequency to be 'hourly', got '%s'", parsedData.Response.Frequency)
	}

	// Add more value checks as needed

	// Optionally, you can print the parsed struct for debugging
	t.Logf("Parsed Data: %+v", parsedData)
}
