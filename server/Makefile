# Makefile for battleship
# Platform: Linux

POSTGRES_CONTAINER_NAME=battleship-postgres
RABBITMQ_CONTAINER_NAME=battleship-rabbit

POSTGRES_PORT=5432
RABBITMQ_PORT=5672
RABBITMQ_MANAGER_PORT=15672

VOLUME=battleship-data

# Time to start container (in seconds)
TIME_TO_START_CONTAINER=10

.PHONY: build run run_postgres run_rabbitmq stop clean

# TODO create a docker-compose file to run the application

build:
	go build -o battleship. ./cmd/server.go

run: run_postgres run_rabbitmq build
	CONFIG_PATH=config/local.yaml ./battleship

run_postgres:
	if [ -z $$(docker ps -a -q -f name=$(POSTGRES_CONTAINER_NAME)) ]; then \
		docker run --name $(POSTGRES_CONTAINER_NAME) -e POSTGRES_PASSWORD=mysecretpassword -p $(POSTGRES_PORT):5432 -v $(VOLUME):/var/lib/postgresql/data -d postgres; \
		sleep $(TIME_TO_START_CONTAINER); \
	else \
		if [ $$(docker inspect -f '{{.State.Running}}' $(POSTGRES_CONTAINER_NAME)) = "false" ]; then \
 			docker start $(POSTGRES_CONTAINER_NAME); \
 			sleep $(TIME_TO_START_CONTAINER); \
		fi \
	fi

run_rabbitmq:
	if [ -z $$(docker ps -a -q -f name=$(RABBITMQ_CONTAINER_NAME)) ]; then \
		docker run -d --name $(RABBITMQ_CONTAINER_NAME) -p $(RABBITMQ_PORT):5672 -p $(RABBITMQ_MANAGER_PORT):15672 -v $(VOLUME):/var/lib/rabbitmq rabbitmq:3-management; \
		sleep $(TIME_TO_START_CONTAINER); \
	else \
		if [ $$(docker inspect -f '{{.State.Running}}' $(RABBITMQ_CONTAINER_NAME)) = "false" ]; then \
			docker start $(RABBITMQ_CONTAINER_NAME); \
			sleep $(TIME_TO_START_CONTAINER); \
		fi \
	fi

stop:
	docker stop $(POSTGRES_CONTAINER_NAME)
	docker stop $(RABBITMQ_CONTAINER_NAME)

clean: stop
	docker rm $(POSTGRES_CONTAINER_NAME)
	docker rm $(RABBITMQ_CONTAINER_NAME)
	rm -rf $(VOLUME)
	rm -f battleship