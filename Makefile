GOPATH := $(shell go env GOPATH)

all: check install agent browser trayIcon desktop

dev:
	@go run . restart -d

run:
	@go run . restart

install:
	@go install

check:
	@go run -tags install installation.go
	@if [ $$? -ne 0 ]; then \
		echo "Error occured. Exiting."; \
		exit 1; \
	fi


agent:
	@cd agent/internal/frontend && $(MAKE) -f Makefile build

browser:
	@cd browser && $(MAKE) -f Makefile install

trayIcon:
	@cd icon && $(MAKE) -f Makefile install

desktop:
	# we need to start the daemon first before compiling the desktop app.
	@go run . restart & 	 
	@echo "waiting for daemon to setup (4secs)..."  && sleep 4
	@cd desktop && $(MAKE) -f Makefile install
	@go run . stop

.PHONY: agent browser trayIcon desktop  dev install all