all: build lint test

build:
	go build ./conditional

lint:
	golint ./conditional
	go vet ./conditional

test:
	go test -v --coverprofile coverage ./conditional --trace=trace
	go tool cover -func=coverage

coverage: test
	go tool cover -html=coverage

trace: test
	go tool trace trace
