GOPATH := $(shell go env GOPATH)

build:
	@go generate && wails build --clean
	
dev:
	@go generate && wails build --clean && ./build/bin/smDesktop -dev

debug:
	@go generate && wails build --clean --devtools && echo "\n\nNow press SHFT+CTRL+F12 to open the devtool panel\n\n" && ./build/bin/smDesktop -dev
	
run:
	@go generate && wails build --clean && ./build/bin/smDesktop

install:
	@go generate && wails build --clean && cp ./build/bin/smDesktop $(GOPATH)/bin/


.PHONY: build dev debug run install
