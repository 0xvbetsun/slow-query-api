# Slow Query API for PostgreSQL

![Build Status](https://github.com/vbetsun/slow-query-api/workflows/CI/badge.svg)
[![coverage](https://codecov.io/gh/vbetsun/slow-query-api/branch/master/graph/badge.svg)](https://codecov.io/gh/vbetsun/slow-query-api)
[![GoReport](https://goreportcard.com/badge/github.com/vbetsun/slow-query-api)](https://goreportcard.com/report/github.com/vbetsun/slow-query-api)
![license](https://img.shields.io/github/license/vbetsun/slow-query-api)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/vbetsun/slow-query-api.svg)](https://github.com/vbetsun/slow-query-api)
[![GoDoc](https://pkg.go.dev/badge/github.com/vbetsun/slow-query-api)](https://pkg.go.dev/github.com/vbetsun/slow-query-api)

## Install

```sh
# HTTPS
git clone https://github.com/vbetsun/slow-query-api.git

# SSH
git clone git@github.com:vbetsun/slow-query-api.git

cd slow-query-api
```

## Deploy

```sh
cp ./deployments/.env.example ./deployments/.env
```

### Env variables

| Name               | Value     | Description                       |
|--------------------|-----------|-----------------------------------|
|PORT                | 8080      | port for running app              |
|PG_HOST             | postgres  | host for running DB               |
|PG_PORT             | 5432      | port for running DB               |
|PG_USER             | user      | username for DB auth              |
|PG_PASSWORD         | str0ngPass| pass for DB auth                  |
|PG_DATABASE         | postgres  | DB name                           |
|SLOW_QUERY_DURATION | 1000      | max query duration in milliseconds|

```sh
docker compose -f ./deployments/docker-compose.yml up -d
```

## Test

```sh
make test

make cover

```