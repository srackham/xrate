test:
	go test ./...

install:
	CGO_ENABLED=0 go install
