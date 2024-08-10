GOPATH := $(shell go env GOPATH)

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


.PHONY: browser trayIcon desktop
browser:
	@cd browser && $(MAKE) -f makefile install

trayIcon:
	@cd icon && $(MAKE) -f makefile install

desktop:
	@cd desktop && $(MAKE) -f makefile install




all: check install browser trayIcon desktop