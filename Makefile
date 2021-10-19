test:
	go test -v -race

coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic

cover:
	make coverage
	go tool cover -html=coverage.txt

build:
	go build ./...

run command:
	go run example/command/main.go 

run listener:
	go run example/listener/main.go 

pull:
	git pull --no-ff
