build:
	@cd ../frontend/ && rm -f *_templ.go && npm run build && templ generate
compile:
	@go run .

run: build compile

