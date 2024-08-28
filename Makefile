.PHONY: test format vet staticcheck magneticod magneticow

all: test magnetico

magnetico:
	go install --tags fts5 .

vet:
	go vet ./...

test:
	CGO_ENABLED=1 go test --tags fts5 -v -race ./...
	CGO_ENABLED=0 go test -v ./...

format:
	gofmt -w ./dht/
	gofmt -w ./metadata/
	gofmt -w ./persistence/
	gofmt -w ./web/
	gci write -s standard -s default -s "prefix(github.com/tgragnato/magnetico)" -s blank -s dot .
