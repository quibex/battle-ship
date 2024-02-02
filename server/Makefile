POSTGRES_CONTAINER_NAME=battleship-postgres
RABBITMQ_CONTAINER_NAME=battleship-rabbit

POSTGRES_PORT=5432
RABBITMQ_PORT=5672
RABBITMQ_MANAGER_PORT=15672

VOLUME=battleship-data



build:
	go build -o server ./cmd/server.go

run_postgres:
	if [ -z $$(docker ps -a -q -f name=$(POSTGRES_CONTAINER_NAME)) ]; then \
		docker run --name $(POSTGRES_CONTAINER_NAME) -e POSTGRES_PASSWORD=mysecretpassword -p $(POSTGRES_PORT):5432 -v $(VOLUME):/var/lib/postgresql/data -d postgres; \
	else \
		if [ $$(docker inspect -f '{{.State.Running}}' $(POSTGRES_CONTAINER_NAME)) = "false" ]; then \
 			docker start $(POSTGRES_CONTAINER_NAME); \
		fi \
	fi

run_rabbitmq:
	if [ -z $$(docker ps -a -q -f name=$(RABBITMQ_CONTAINER_NAME)) ]; then \
		docker run -d --name $(RABBITMQ_CONTAINER_NAME) -p $(RABBITMQ_PORT):5672 -p $(RABBITMQ_MANAGER_PORT):15672 -v $(VOLUME):/var/lib/rabbitmq rabbitmq:3-management; \
	else \
		if [ $$(docker inspect -f '{{.State.Running}}' $(RABBITMQ_CONTAINER_NAME)) = "false" ]; then \
		docker start $(RABBITMQ_CONTAINER_NAME); \
		fi \
	fi

run:
	CONFIG_PATH=config/local.yaml go run ./cmd/server.go

stop:
	docker stop $(POSTGRES_CONTAINER_NAME)
	docker stop $(RABBITMQ_CONTAINER_NAME)