BINARY_NAME=chicago.wthr.cloud

build:
	GOARCH=amd64 GOOS=darwin go build -o release/${BINARY_NAME}-darwin .
	GOARCH=amd64 GOOS=linux go build -o release/${BINARY_NAME}-linux .
	GOARCH=amd64 GOOS=windows go build -o release/${BINARY_NAME}-windows .

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows

test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all
