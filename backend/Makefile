APP_BINARY_NAME=ecochitas_backend

.PHONY: run
run:
	go run ./cmd/api

.PHONY: build
build:
	go build -o bin/$(APP_BINARY_NAME) ./cmd/api

.PHONY: test
test:
	go test ./...

.PHONY: infra-up
infra-up:
	docker compose up -d

.PHONY: infra-down
infra-down:
	docker compose down

