.PHONY: run
run:
	@wgo run -file .html -file .css ./cmd/main.go

.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v
