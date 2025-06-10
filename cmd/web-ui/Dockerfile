FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN cd cmd/web-ui && go build -o /web-ui

FROM alpine:latest
COPY --from=builder /web-ui /web-ui
ENTRYPOINT ["/web-ui"]