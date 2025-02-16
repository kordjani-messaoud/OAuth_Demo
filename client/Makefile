build:
	@echo "Compiling..."
	@GOOS=linux GOARCH=amd64 go build -o bin/client client.go

run: build
	@./bin/client