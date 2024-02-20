APPNAME := orders-api
VERSION := 1.0
SUMMARY := Service that gets, creates and updates limit orders.

export REPORTS_DIR=./reports

build:
	mkdir -p build
	GOOS=$(GOOS) GOARCH=$(GOARCH) APPNAME=$(APPNAME) ./scripts/build

run:
	./build/$(APPNAME)

clean:
	rm -rf build