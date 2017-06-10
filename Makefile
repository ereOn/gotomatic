MODULES=conditional configuration executor time trigger

ifeq ($(OS),Windows_NT)
EXT:=.exe
else
EXT:=
endif

all: build lint test

build:
	go build -o bin/gotomate${EXT} ./gotomate

lint:
	for MODULE in $(MODULES); do \
		golint ./$${MODULE}; \
		go vet ./$${MODULE}; \
	done

test:
	for MODULE in $(MODULES); do \
		go test -v --coverprofile coverage.$${MODULE} ./$${MODULE} --trace=trace.$${MODULE}; \
		go tool cover -func=coverage.$${MODULE}; \
	done

coverage: test
	for MODULE in $(MODULES); do \
		go tool cover -html=coverage.$${MODULE}; \
	done

trace: test
	for MODULE in $(MODULES); do \
		go tool trace trace.$${MODULE}; \
	done
