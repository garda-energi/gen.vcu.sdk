test:
	go test -v -coverprofile=go.coverage.out

coverage:
	go tool cover -html=go.coverage.out

build:
	go build ./...

run:
	go run example/main.go 

pull:
	git pull --no-ff

release:
	make test
	git push
	git checkout master
	git merge development
	git push
	git checkout development