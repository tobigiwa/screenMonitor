build:
	@rm -rf *_templ.go && npm run build && templ generate

fmt:
	@templ fmt .

live:
	@go generate && templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

run: build fmt live