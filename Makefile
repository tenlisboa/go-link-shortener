
deps:
	go mod tidy

build: deps
	go build -gcflags="all=-N -l" -o bin/app cmd/main.go
