FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN cd cmd/disk-service && go build -o /disk-service

FROM alpine:latest
COPY --from=builder /disk-service /disk-service
ENTRYPOINT ["/disk-service"]