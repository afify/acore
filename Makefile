include .env

COMMIT := $(shell git rev-parse --short HEAD)
export COMMIT

.PHONY: deploy ensure-infra migrate migrate-new restart-infra clean-all

all: deploy

ensure-infra:
	@printf "\033[35m*** Ensuring infra services are running…\033[0m\n"
	@for svc in acore-postgres acore-redis acore-traefik acore-grafana; do \
		status=$$(docker inspect -f '{{.State.Running}}' $$svc 2>/dev/null || echo false); \
		if [ "$$status" != "true" ]; then \
			printf "\033[33m*** Starting $$svc…\033[0m\n"; \
			docker compose up -d $$svc; \
		fi; \
	done

reload-infra:
	@printf "\033[35m*** Reloading infra services with new settings…\033[0m\n"
	docker compose pull acore-postgres acore-redis acore-traefik acore-grafana
	docker compose up -d --build --force-recreate acore-postgres acore-redis acore-traefik acore-grafana

# Blue/Green deploy: ensure infra → start green → stop blue → rebuild blue → stop green → bring green back
deploy: ensure-infra
	@printf "\033[35m*** Stopping blue (traffic → green)…\033[0m\n"
	docker compose stop acore-blue

	@printf "\033[35m*** Rebuilding & starting blue…\033[0m\n"
	docker compose build acore-blue
	docker compose up -d --build --no-deps acore-blue

	@printf "\033[35m*** Waiting for blue to be healthy…\033[0m\n"
	until [ "$$(docker inspect -f '{{.State.Health.Status}}' acore-blue)" = "healthy" ]; do sleep 1; done

	@printf "\033[35m*** Stopping green (traffic → blue)…\033[0m\n"
	docker compose stop acore-green

	@printf "\033[35m*** Restarting green (both up)…\033[0m\n"
	docker compose up -d --build --no-deps acore-green

	@printf "\033[35m*** Done — Blue & Green both healthy. \033[0m\n"

	@printf "\033[35m*** Cleaning up old acore images…\033[0m\n"
	@docker images "${APP_NAME}" --format "{{.ID}}" | xargs -r docker rmi 2>/dev/null || true

migrate:
	@printf "\033[35m*** Running migrations…\033[0m\n"
	docker compose run --rm acore-migrate \
		-path /migrations -database "${PG_URL}" up

migrate-new:
	@printf "\033[35m*** Creating a new migration file…\033[0m\n"
	@read -p "Migration name: " name; \
		timestamp=$$(date -u +"%Y%m%d%H%M%S"); \
		snake=$$(echo "$$name" \
		| tr '[:upper:]' '[:lower:]' \
		| tr ' ' '_' \
		| tr -cd 'a-z0-9_-'); \
		file="database/migrations/$${timestamp}_$${snake}.up.sql"; \
		touch "$$file"; \
	printf "\033[35m*** Created %s\033[0m\n" "$$file"

clean-all:
	@printf "\033[35m*** Stopping all containers…\033[0m\n"
	@docker stop $$(docker ps -q) 2>/dev/null || true

	@printf "\033[35m*** Pruning all containers, images, networks, volumes, and caches…\033[0m\n"
	@docker system prune -af --volumes
