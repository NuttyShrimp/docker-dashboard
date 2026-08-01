all: build

build:
	@go build -o main cmd/api/main.go

run:
	@go run cmd/api/main.go

watch:
	@go tool air -c .air.toml

dead:
	@go tool deadcode ./...

.PHONY: all build run watch dead
