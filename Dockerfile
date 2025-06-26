FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o gateway ./cmd/main.go

FROM scratch
COPY --from=builder /app/gateway /gateway
ENTRYPOINT ["/gateway"]
