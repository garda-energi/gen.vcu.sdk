test:
	go test -v -race

coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic

cover:
	make coverage
	go tool cover -html=coverage.txt

build:
	go build ./...

run:
	go run example/main.go 

pull:
	git pull --no-ff
