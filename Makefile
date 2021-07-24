test:
	go test -v
build:
	go build ./...
run:
	go run example/main.go 
pull:
	git pull --no-ff
release:
	git checkout master
	git merge development
	git push
	git checkout development