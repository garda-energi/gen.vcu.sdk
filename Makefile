test:
	go test -v -race -cover

cover:
	go test -v -coverprofile=go.cover.out
	go tool cover -html=go.cover.out

build:
	go build ./...

run:
	go run example/main.go 

pull:
	git pull --no-ff
