run: build
	@go run cmd\web\

build:
	@go build -o bin/web cmd\web\