all: build lint test

build:
	go build ./conditional
	go build ./time

lint:
	golint ./conditional
	go vet ./conditional
	golint ./time
	go vet ./time

test:
	go test -v --coverprofile coverage.conditional ./conditional --trace=trace.conditional
	go test -v --coverprofile coverage.time ./time --trace=trace.time
	go tool cover -func=coverage.conditional
	go tool cover -func=coverage.time

coverage: test
	go tool cover -html=coverage.conditional
	go tool cover -html=coverage.time

trace: test
	go tool trace trace.conditional
	go tool trace trace.time
