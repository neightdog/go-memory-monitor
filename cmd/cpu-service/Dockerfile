FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN cd cmd/cpu-service && go build -o /cpu-service

FROM alpine:latest
COPY --from=builder /cpu-service /cpu-service
ENTRYPOINT ["/cpu-service"]