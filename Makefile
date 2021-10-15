build:
	go build -o bin/n2m cmd/n2m/main.go

run:
	go run cmd/n2m/main.go

install: build
	cp bin/n2m /usr/local/bin
