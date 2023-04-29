lint:
	golangci-lint run -c golangci.yaml
build:
	CGO_ENABLED=0 go build -o bin/app
run: build
	sudo docker-compose up --remove-orphans
test:
	go test ./... -count=1
mocks:
	sh /mocks/generate.sh