help: ## You are here! showing all command documenentation.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#===================#
#== Env Variables ==#
#===================#
DOCKER_COMPOSE_FILE ?= docker-compose.yml

environment: ## Init environment setup
environment:
	docker compose -f ${DOCKER_COMPOSE_FILE} up --build -d
	docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate up

environment-clean: ## Clear environment setup (container & volume)
environment-clean:
	docker compose -f ${DOCKER_COMPOSE_FILE} down -v

shell-db: ## Enter to database console
shell-db:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec db psql -U postgres -d postgres

server: ## Run http server
server:
	go run main.go
