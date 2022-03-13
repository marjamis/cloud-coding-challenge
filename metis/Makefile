.DEFAULT_GOAL := helper
GIT_COMMIT ?= $(shell git rev-parse --short=12 HEAD || echo "NoGit")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
TEXT_RED = \033[0;31m
TEXT_BLUE = \033[0;34;1m
TEXT_GREEN = \033[0;32;1m
TEXT_NOCOLOR = \033[0m

DOCKER_REGISTRY=042836644938.dkr.ecr.us-west-2.amazonaws.com
DOCKER_REPOSITORY=metis
DOCKER_IMAGE_TAG_APP=metis-v1
DOCKER_IMAGE_TAG_REDIRECTOR=metis-redirecter

helper: # Adapted from: https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@echo "Available targets..." # @ will not output shell command part to stdout that Makefiles normally do but will execute and display the output.
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deploy: build build_image push_image ## Build and push the docker image to ECR

build_image: ## Build the docker image
	docker build -t $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):$(DOCKER_IMAGE_TAG_APP) .

run_image: ## Run the docker image locally for test
	docker run -p 8081:8081 --name metis --link test_db:db -e UserName=root -e Password=testing -e Endpoint=db -e Port=3306 -e Name=testdb $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):$(DOCKER_IMAGE_TAG_APP)

push_image: ## Push the built image to ECR
	docker push $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):$(DOCKER_IMAGE_TAG_APP)

build: ## Build the staticly compiled binary
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/metis.go metis

prod: ## Runs the prod version of the application
	$(MAKE) command OPTIONS="-p 8080"

dev: ## Runs a dev version of the application
	# Setup a fresh database docker container to be used for local testing
	docker stop test_db && docker rm test_db || echo "Image not running"
	docker run -p 3306:3306 -dit -e MYSQL_DATABASE=testdb -e MYSQL_ROOT_PASSWORD=testing --name test_db mariadb:10.1.22
	docker cp ./test/tempdb_data.sql test_db:/schema.sql
	sleep 20
	docker exec -i test_db /usr/bin/mysql -uroot -ptesting -e "source schema.sql; \. schema.sql"

clean: ## Cleans up any old/unneeded items
