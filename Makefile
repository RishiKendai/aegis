K6_IMAGE = grafana/k6
APP_NAME = aegis
TEST_DIR = load-tests
PWD := $(shell pwd)

## run: Run the aegis server
run:
	go run ./cmd

help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the aegis server
build:
	go build -o build/${APP_NAME} ./cmd

## load: Run k6 load test in Docker (mounts load-tests folder)
load:
	docker run --network host --rm -i -v $(PWD):/app grafana/k6 run /app/load-tests/load.js

## prom: Run prometheus in Docker
prom:
	docker run -d --name prometheus -p 9090:9090 -v $(PWD)/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

## grafana: Run grafana in Docker
grafana:
	docker compose -f metrics.yml up -d