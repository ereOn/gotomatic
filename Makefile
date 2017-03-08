all: build lint test

build:
	go build ./conditional

lint:
	golint ./conditional
	go vet ./conditional

test:
	go test --race --coverprofile coverage ./conditional
	go tool cover -func=coverage

coverage: test
	go tool cover -html=coverage
