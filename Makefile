/sudoSWAG := go run github.com/swaggo/swag/cmd/swag@v1.16.6
POD_NAME := pismo-stack
JAEGER_NAME := pismo-jaeger
API_NAME := pismo-api-test
IMAGE_NAME := localhost/pismo-api:latest

.PHONY: tidy fmt test run swagger build compose-up compose-down podman-build podman-up podman-down podman-logs

tidy:
	go mod tidy

fmt:
	gofmt -w ./cmd ./internal ./docs

test:
	go test ./...

run:
	go run ./cmd/api

swagger:
	$(SWAG) init -g cmd/api/main.go -o docs --outputTypes go,json,yaml

build:
	podman build -t $(IMAGE_NAME) .

compose-up:
	docker compose up --build

compose-down:
	docker compose down

podman-build:
	podman build -t $(IMAGE_NAME) .

podman-up:
	podman pod create --replace --name $(POD_NAME) -p 8080:8080 -p 16686:16686 -p 4317:4317 -p 4318:4318
	podman run -d --pod $(POD_NAME) --name $(JAEGER_NAME) jaegertracing/all-in-one:1.67.0
	podman run -d --pod $(POD_NAME) --name $(API_NAME) \
		-e PORT=8080 \
		-e DATABASE_URL='file:/data/pismo.db?_pragma=foreign_keys(1)' \
		-e OTEL_SERVICE_NAME=pismo-api \
		-e OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318 \
		-e OTEL_EXPORTER_OTLP_INSECURE=true \
		-v pismo_data:/data \
		$(IMAGE_NAME)

podman-down:
	-podman stop $(API_NAME) $(JAEGER_NAME)
	-podman rm $(API_NAME) $(JAEGER_NAME)
	-podman pod rm -f $(POD_NAME)

podman-logs:
	podman logs -f $(API_NAME)
