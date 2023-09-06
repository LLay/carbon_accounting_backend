set time to be correct thing
upsert data

request EIA data with net/http
put in DB

Make an API to query data
Frontend queries data

cron to fetch daily data

finish local version first? 
then host?
ya
what's left?

aggregate endpoint
start, end
- {
    "<FuelType>":[]
  }


FE makes requests to DB
- aggregate hourly data by fuel type. why? speed? ease of use

Host
- InfluxDB hosted
- Domain name
- webserver dockerized and deployed 
- server dockerized and deployed