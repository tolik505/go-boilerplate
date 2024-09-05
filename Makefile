NAME = goboilerplate
#
#Variables for mysql integration tests
#
MYSQL_IMAGE = mysql:8.0
MYSQL_CONTAINER_NAME = $(NAME)-mysql-test
MYSQL_PORT = 3369
MYSQL_HOST = localhost
MYSQL_DB_NAME = test

#
# Function for mysql integration tests
#
define test_mysql_run
	docker run -t --rm --name $(MYSQL_CONTAINER_NAME) \
	-e MYSQL_DATABASE=$(MYSQL_DB_NAME) \
	-e MYSQL_ROOT_PASSWORD=password \
	-p $(MYSQL_PORT):3306 \
	-d $(MYSQL_IMAGE)
endef

# Runs tests locally
test:
	cd pkg && env TEST_DB_HOST=$(MYSQL_HOST) TEST_DB_PORT=$(MYSQL_PORT) TEST_DB_NAME=$(MYSQL_DB_NAME) go test -v ./...

# Builds binary files for both http and grpc servers
build-binary:
	cd cmd/httpserver && env GOOS=linux go build -o ../../bin/httpserver
	cd cmd/grpcserver && env GOOS=linux go build -o ../../bin/grpcserver

# Runs mysql container for test db
test-mysql:
	$(call test_mysql_run)

# Generates GraphQL resolvers boilerplate
gqlgen:
	cd config/graphql && go run github.com/99designs/gqlgen generate

# Builds binaries and runs application in docker containers
docker-compose: build-binary
	docker-compose up --build

# Runs db migration
# migrate should be installed https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
migrate:
	migrate -path /migrations -database "mysql://root:password@tcp(localhost:3306)/goboilerplate" up

# Generates protobuf file from config/pb/example.proto schema
protoc:
	cd config/pb &&\
	protoc --go_out=plugins=grpc,paths=source_relative:../../pkg/grpcapp/pb *.proto
