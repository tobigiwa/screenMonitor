build:
	@rm -rf *_templ.go && npm run build && templ generate

fmt:
	@templ fmt .

live:
	@templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

run: build fmt live