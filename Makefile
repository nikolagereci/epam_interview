develop:
	sudo docker-compose up --remove-orphans cassandra cassandra-load-keyspace zookeeper broker
build:
	CGO_ENABLED=0 go build -C cmd -o ../bin/app
dockerbuild: build
	sudo docker build -t companies:$(cat VERSION) .
lint:
	golangci-lint run -c golangci.yaml
run: build
	sudo docker-compose up --remove-orphans
unit-test:
	go test ./... -count=1
integration-test: build
	sudo docker-compose up -d --remove-orphans cassandra cassandra-load-keyspace zookeeper broker
	go test ./... -count=1 -tags=integration
	sudo docker-compose down
genmocks:
	sh mocks/generate.sh