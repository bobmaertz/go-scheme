clean: 
	rm /bin/go

build:
	go build -o bin/scheme ./...

test: 
	go test ./... -coverprofile cover.out