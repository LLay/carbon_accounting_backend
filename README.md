# Golang API Skeleton
> A simple API skeleton written in Go with metrics pre-configurated


[![Build Status](https://travis-ci.org/michelaquino/golang_api_skeleton.svg?branch=master)](https://travis-ci.org/michelaquino/golang_api_skeleton)
[![License][license-image]][license-url]


## Includes
  - [Zap - Uber Log library](https://github.com/uber-go/zap)
  - [Echo Framework](https://github.com/labstack/echo)
  - [MongoDB driver](https://github.com/mongodb/mongo-go-driver)
  - [Go-Redis](github.com/go-redis/redis)
  - [Prometheus](https://github.com/prometheus) 
  - [Prometheus PushGateway](https://github.com/prometheus/pushgateway) 
  - [Grafana](https://grafana.com/) 

## Dependencies

- Docker
- Docker Compose

## Configuration
- Docker Compose
    - Nginx with `proxy_pass` preconfigured
    - API
    - MongoDB
    - Redis
    - Prometheus
    - Prometheus Push Gateway
    - Grafana

## Run
`make run`

## Usage
`curl http://localhost/healthcheck`

`curl -i -X POST -H 'Content-Type: application/json' -d '{"name": "user name", "email": "user@email.com"}' http://localhost/user`

### Metrics
Access:
- http://localhost:3000 to view Grafana metrics pre-configurated
- http://localhost:9090 to view Prometheus server

[license-image]: https://img.shields.io/badge/License-GPL3.0-blue.svg
[license-url]: LICENSE
[travis-image]: https://img.shields.io/travis/michelaquinoe/golang_api_skeleton/master.svg
[travis-url]: https://travis-ci.org/michelaquino/golang_api_skeleton

###
InfluxDB
### UI
http://localhost:8086/signin

Docs:
- https://docs.influxdata.com/influxdb/v2.7/get-started/#influxdb-user-interface-ui
- https://docs.influxdata.com/influxdb/v2.7/visualize-data/

### CLI
influx -username admin -password adminalskdjfhalksdjh

https://docs.influxdata.com/influxdb/v2.7/tools/influx-cli/

create auth token
```
influx auth create \
  --all-access \
  --host http://localhost:8086 \
  --org my_organization \
  --token admin

ID			Description	Token												User Name	User ID			Permissions
0bc25448b773a000			D-CqTaAVm3PWGHlyfrSXh9efMZpnfS7_L2EP06DGa17nB8FSetZcaYfR4zqPswN_JVIhTIOtzJ9j7cz6gr3icw==	admin		0bc21576dd22e000	[read:orgs/322813f618e44b4d/authorizations write:orgs/322813f618e44b4d/authorizations read:orgs/322813f618e44b4d/buckets write:orgs/322813f618e44b4d/buckets read:orgs/322813f618e44b4d/dashboards write:orgs/322813f618e44b4d/dashboards read:/orgs/322813f618e44b4d read:orgs/322813f618e44b4d/sources write:orgs/322813f618e44b4d/sources read:orgs/322813f618e44b4d/tasks write:orgs/322813f618e44b4d/tasks read:orgs/322813f618e44b4d/telegrafs write:orgs/322813f618e44b4d/telegrafs read:/users/0bc21576dd22e000 write:/users/0bc21576dd22e000 read:orgs/322813f618e44b4d/variables write:orgs/322813f618e44b4d/variables read:orgs/322813f618e44b4d/scrapers write:orgs/322813f618e44b4d/scrapers read:orgs/322813f618e44b4d/secrets write:orgs/322813f618e44b4d/secrets read:orgs/322813f618e44b4d/labels write:orgs/322813f618e44b4d/labels read:orgs/322813f618e44b4d/views write:orgs/322813f618e44b4d/views read:orgs/322813f618e44b4d/documents write:orgs/322813f618e44b4d/documents read:orgs/322813f618e44b4d/notificationRules write:orgs/322813f618e44b4d/notificationRules read:orgs/322813f618e44b4d/notificationEndpoints write:orgs/322813f618e44b4d/notificationEndpoints read:orgs/322813f618e44b4d/checks write:orgs/322813f618e44b4d/checks read:orgs/322813f618e44b4d/dbrp write:orgs/322813f618e44b4d/dbrp read:orgs/322813f618e44b4d/notebooks write:orgs/322813f618e44b4d/notebooks read:orgs/322813f618e44b4d/annotations write:orgs/322813f618e44b4d/annotations read:orgs/322813f618e44b4d/remotes write:orgs/322813f618e44b4d/remotes read:orgs/322813f618e44b4d/replications write:orgs/322813f618e44b4d/replications]
```

### GraphQL

### Run it
cd src/graphql
go run server.go
http://localhost:8080/


#### Generate
cd src/graphql
go run github.com/99designs/gqlgen generate