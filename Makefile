BINARY=imap-mailbox-exporter

build:
	go build -o $(BINARY) ./...

run: build
	./$(BINARY)

test:
	go test ./...

test-coverage:
	go test ./... -cover

test-coverage-html:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out