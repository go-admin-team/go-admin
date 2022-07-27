PROJECT:=go-admin

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -a -installsuffix -o go-admin .
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix -o go-admin .
build-sqlite:
	go build -tags sqlite3 -ldflags="-w -s" -a -installsuffix -o go-admin .
#.PHONY: test
#test:
#	go test -v ./... -cover

#.PHONY: docker
#docker:
#	docker build . -t go-admin:latest
