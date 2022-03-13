DOCKER_REGISTRY=042836644938.dkr.ecr.us-west-2.amazonaws.com
DOCKER_REPOSITORY=metis
DOCKER_IMAGE_TAG_APP=metis-v1
DOCKER_IMAGE_TAG_REDIRECTOR=metis-redirecter
GOPATH=/media/sharea/data/go/

default:

deploy: build build_image push_image

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/metis.go metis

setup:
	go get "github.com/go-sql-driver/mysql" "github.com/gorilla/mux" "github.com/satori/go.uuid" "github.com/stretchr/testify/assert"

run:
	UserName=root Password=testing Endpoint=localhost Port=3306 Name=testdb go run main.go routes.go functions.go handlers.go structs.go

build_image:
	docker build -t $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):$(DOCKER_IMAGE_TAG_APP) .
	cd ../../../redirecter && docker build -t $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):$(DOCKER_IMAGE_TAG_REDIRECTOR) .

run_image:
	docker run -p 8081:8081 --name metis --link test_db:db -e UserName=root -e Password=testing -e Endpoint=db -e Port=3306 -e Name=testdb $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):$(DOCKER_IMAGE_TAG_APP)

push_image:
	 docker push $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):$(DOCKER_IMAGE_TAG_APP)
	 docker push $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):$(DOCKER_IMAGE_TAG_REDIRECTOR)

local_test:
	UserName=root Password=testing Endpoint=localhost Port=3306 Name=testdb  go test -v -cover -coverprofile=./metis.out metis
	go tool cover -html=metis.out -o ./coverage.html

local_testd:
	# Setup a fresh database docker container to be used for local testing
	docker stop test_db && docker rm test_db || echo "Image not running"
	docker run -p 3306:3306 -dit -e MYSQL_DATABASE=testdb -e MYSQL_ROOT_PASSWORD=testing --name test_db mariadb:10.1.22
	docker cp ./misc/tempdb_data.sql test_db:/schema.sql
	sleep 20
	docker exec -i test_db /usr/bin/mysql -uroot -ptesting -e "source schema.sql; \. schema.sql"
