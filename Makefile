test:
	go test -v -race

cover:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic


cover-read:
	go tool cover -html=coverage.txt

build:
	go build ./...

run:
	go run example/main.go 

pull:
	git pull --no-ff
