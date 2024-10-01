vet:
	go vet ./...

fmt:
	go fmt ./...

bin:
	go build -o bin/webhook-adapter

build: vet fmt bin

run:
	./bin/webhook-adapter

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

clean:
	rm -rf bin coverage.out
