#!/bin/bash

# Define a temporary file to store the response
response_file=tmp.json

curl 'https://api.eia.gov/v2/electricity/rto/fuel-type-data/data/?api_key=CZdQsisRJzwOfqUWV3jiMPNEx3ZbHcuJ2VQus04i' \
  -H 'authority: api.eia.gov' \
  -H 'accept: application/json, text/plain, */*' \
  -H 'accept-language: en-US,en;q=0.9' \
  -H 'content-type: application/json' \
  -H 'dnt: 1' \
  -H 'origin: https://www.eia.gov' \
  -H 'referer: https://www.eia.gov/' \
  -H 'sec-ch-ua: "Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v="116"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "macOS"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-site' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36' \
  -H 'x-params: {"frequency":"hourly","data":["value"],"facets":{},"start":"2023-09-02T00","end":null,"sort":[{"column":"period","direction":"desc"}],"offset":0,"length":5000}' \
  --compressed > "$response_file"

# Output the response file path
echo "$response_file"
