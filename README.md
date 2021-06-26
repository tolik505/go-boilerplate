# Go-boilerplate
A boilerplate with examples for starting up Golang projects.

## Goals
1. To help developers to start up their first Web/gRPC application in Go.
2. To show how to structure app and implement common use-cases.
3. To show how to test those common use-cases.
4. To introduce how modern technologies like GraphQL and gRPC can be used in Golang application.

## Features
1. HTTP server (supports GraphQL and regular HTTP endpoints).
2. GRPC server.
3. Using [echo](https://echo.labstack.com), [GORM](https://gorm.io/index.html) and [gqlgen](https://gqlgen.com/getting-started/).
4. Unit and integration tests.
5. Multilayered structure.
6. Many examples including file upload and gRPC endpoint with tests.

### Running with docker
For running the service in docker container use ```make run-docker-compose```.

### Regenerating of GraphQL boilerplate
Use ```make gqlgen```

### Regenerating of gRPC protobuf file
Use ```make protoc```

### DB migrations
If you run docker-compose then DB migration will be done automatically.

If you run the app locally use ```make migrate```.

When tests are being run DB migrations will be applied on test DB automatically.

### Dependency injection
The project uses [wire](https://github.com/google/wire/blob/main/docs/guide.md) for building a dependency tree without pain.

### Tests
The service includes mysql integration tests, therefore test database should be running.

Use ```make test-mysql``` to run test DB.

Then ```test``` to run all tests. Also tests can be run directly in modern IDE like Goland if Go installed locally.

In order to generate mocks, use [mockery](https://github.com/vektra/mockery)

### Suggestions, contributions
Are appreciated.