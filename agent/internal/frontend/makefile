build:
	@rm -rf *_templ.go && templ generate

fmt:
	@templ fmt .

dev: build fmt live

live:
	@go generate && templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v

tailwind:
	@tailwindcss -i assets/styles/input.css  -o assets/styles/output.css --watch
