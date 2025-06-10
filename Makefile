build:
    go build -o bin/memory-service ./cmd/memory-service
    go build -o bin/disk-service ./cmd/disk-service
    go build -o bin/web-ui ./cmd/web-ui

run-memory:
    go run ./cmd/memory-service

run-disk:
    go run ./cmd/disk-service

run-ui:
    go run ./cmd/web-ui
