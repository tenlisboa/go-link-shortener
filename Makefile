
deps:
	go mod tidy

build: deps
	go build -o bin/app cmd/main.go
