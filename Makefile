build:
	go build -o bin/gateway ./cmd/main.go

run:
	go run ./cmd/main.go

lint:
	golangci-lint run

docker-build:
	docker build -t api-gateway .

test:
	go test ./...
