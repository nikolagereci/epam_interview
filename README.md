# Documentation

This is an implementation of a company registry.  
Run it with `make run`.  
Test it using the available postman collection (Companies.postman_collection.json)   
The project enables you to perform various usefull operations using the `make` command. 

## Commands
To use the Makefile, navigate to the root directory of the project and enter the following commands:

```shell
make develop
```

This command will start the necessary dependencies for local development using Docker Compose. Specifically, it will start the following services:

- Cassandra
- Cassandra load keyspace
- Zookeeper
- Kafka Broker

```shell
make build
```

This command will compile the application and output a binary file named `app` in the `bin` directory.

```shell
make dockerbuild
```

This command will build a Docker image of the application and tag it with the current version specified in the `VERSION` file. The command also depends on the `build` command, which means it will first compile the application before building the Docker image.

```shell
make lint
```

This command will run the linter tool, `golangci-lint`, on the entire codebase. It is configured to use the `golangci.yaml` configuration file located in the root directory of the project.

```shell
make run
```

This command will start the application and its dependencies using Docker Compose.

```shell
make test
```

This command will run all the tests in the `./...` directory and output the results.
```shell
make genmocks
```

This command will generate mocks using `gomock` for any interfaces located in the `./mocks` directory. It executes the `generate.sh` script in the `./mocks` directory.


## Usage in a CI/CD Pipeline

Here's an example of how these commands could be used in a CI/CD pipeline:

1. **Linting**: The `lint` command can be used to check the code for any syntax errors or style violations.

```shell
make lint
```
2. **Run tests**: The next step is to run the unit tests to ensure that the application is working as expected. This can be done using the `test` command.

```shell
make test
```

3. **Build the application**: The first step is to build the application using the `build` command. This command compiles the code and generates an executable file that can be run on any machine.

```shell
make build
```

4. **Build Docker image**: The `dockerbuild` command builds a Docker image of the application using the Dockerfile in the project root.

```shell
make dockerbuild
```

5. **Deploy**: Once the Docker image is built, it can be deployed to a production environment using a tool like Kubernetes, Docker Swarm, or AWS ECS.

These commands can be automated and run as part of a CI/CD pipeline, allowing for continuous integration, testing, and deployment of the application.