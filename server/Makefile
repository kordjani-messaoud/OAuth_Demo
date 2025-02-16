build:
	@echo "Compiling..."
	@GOOS=linux GOARCH=amd64 go build -o bin/server server.go

run: build
	@./bin/server