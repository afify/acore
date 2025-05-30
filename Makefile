include .env

COMMIT := $(shell git rev-parse --short HEAD)
export COMMIT

.PHONY: all check infra-ensure infra-reload deploy migrate migrate-schema migrate-func migrate-new migrate-dropall clean-all frontend b migrate-clean

all: deploy

check:
	@printf "\033[36m*** go fmt…\033[0m\n"
	@go fmt ./...
	@printf "\033[36m*** go update…\033[0m\n"
	@go get -u ./...
	@printf "\033[36m*** go vet…\033[0m\n"
	@go vet ./...
	@printf "\033[36m*** staticcheck…\033[0m\n"
	@staticcheck ./...
	@printf "\033[36m*** govulncheck…\033[0m\n"
	@govulncheck ./...
	@printf "\033[36m*** gosec…\033[0m\n"
	@gosec ./...
	@printf "\033[32m*** Lint & security checks passed! 🎉\033[0m\n"

infra-ensure:
	@printf "\033[35m*** Ensuring infra services are running…\033[0m\n"
	@for svc in ${APP_NAME}-postgres ${APP_NAME}-redis ${APP_NAME}-migrate ${APP_NAME}-traefik ${APP_NAME}-grafana; do \
		status=$$(docker inspect -f '{{.State.Running}}' $$svc 2>/dev/null || echo false); \
		if [ "$$status" != "true" ]; then \
			printf "\033[33m*** Starting %s\033[0m\n" "$$svc"; \
			docker compose up -d $$svc; \
		fi; \
	done

infra-reload:
	@printf "\033[35m*** Reloading infra services with new settings…\033[0m\n"
	docker compose pull ${APP_NAME}-postgres ${APP_NAME}-redis ${APP_NAME}-traefik ${APP_NAME}-grafana
	docker compose up -d --build --force-recreate ${APP_NAME}-postgres ${APP_NAME}-redis ${APP_NAME}-traefik ${APP_NAME}-grafana

# Blue/Green deploy: ensure infra → start green → stop blue → rebuild blue → start blue → stop green
deploy: infra-ensure
	@printf "\033[35m*** Starting green (new version)…\033[0m\n"
	docker compose up -d --build --no-deps ${APP_NAME}-green

	@printf "\033[35m*** Waiting for green to be healthy…\033[0m\n"
	until [ "$$(docker inspect -f '{{.State.Health.Status}}' ${APP_NAME}-green)" = "healthy" ]; do sleep 1; done

	@printf "\033[35m*** Stopping blue…\033[0m\n"
	docker compose stop ${APP_NAME}-blue

	@printf "\033[35m*** Building blue (new version)…\033[0m\n"
	docker compose up -d --build --no-deps ${APP_NAME}-blue

	@printf "\033[35m*** Starting blue…\033[0m\n"
	docker compose start ${APP_NAME}-blue

	@printf "\033[35m*** Waiting for blue to be healthy…\033[0m\n"
	until [ "$$(docker inspect -f '{{.State.Health.Status}}' ${APP_NAME}-blue)" = "healthy" ]; do sleep 1; done

	@printf "\033[35m*** Stopping green…\033[0m\n"
	docker compose stop ${APP_NAME}-green

	@printf "\033[35m*** Building green for next cycle…\033[0m\n"
	docker compose build ${APP_NAME}-green

	@printf "\033[35m*** Clean docker cache.\033[0m\n"
	@printf "\033[35m*** Clean 💙/💚 builder cache…\033[0m\n"
	docker builder prune --all --force

	@printf "\033[35m*** Prune dangling images (keep acore-blue & acore-green)…\033[0m\n"
	docker image prune --filter "dangling=true" --force

	@printf "\033[35m*** Prune unused volumes…\033[0m\n"
	docker volume prune --force

	@printf "\033[35m*** Prune unused networks…\033[0m\n"
	docker network prune --force
	docker system df


migrate: infra-ensure init-db migrate-func
	@printf "\033[32m*** All migrations complete! 🎉\033[0m\n"

init-db:
	@printf "\033[35m*** Loading init.sql into Postgres…\033[0m\n"
	@cat database/init/init.sql | docker compose exec -T ${APP_NAME}-postgres \
		psql -U ${PG_USER} -d ${PG_NAME} -v ON_ERROR_STOP=1

migrate-schema:
	@printf "\033[35m*** Running migrations…\033[0m\n"
	docker compose run --rm ${APP_NAME}-migrate \
		-path /migrations -database "${PG_URL}" up

migrate-func:
	@printf "\033[35m*** Applying SQL functions…\033[0m\n"
	@for f in database/migrations/functions/*.sql; do \
		printf "\033[33m*** Applying %s\033[0m\n" "$$f"; \
		cat "$$f" | docker compose exec -T ${APP_NAME}-postgres \
			psql -U ${PG_USER} -d ${PG_NAME} \
			-v ON_ERROR_STOP=1; \
	done

migrate-new:
	@printf "\033[35m*** Creating a new migration file…\033[0m\n"
	@read -p "Migration name: " name; \
	timestamp=$$(date -u +"%Y%m%d%H%M%S"); \
	snake=$$(echo "$$name" \
	  | tr '[:upper:]' '[:lower:]' \
	  | tr ' ' '_' \
	  | tr -cd 'a-z0-9_-'); \
	dir="database/migrations"; \
	read -p "Is this a function? (y/N): " isfunc; \
	[ "$$isfunc" = "y" ] && dir="$$dir/functions"; \
	mkdir -p "$$dir"; \
	file="$$dir/$${timestamp}_$${snake}.up.sql"; \
	touch "$$file"; \
	printf "\033[32m*** Created %s\033[0m\n" "$$file"

migrate-dropall: infra-ensure
	@printf "\033[35m*** Resetting public schema…\033[0m\n"
	docker compose exec -T ${APP_NAME}-postgres \
		psql -U ${PG_USER} -d ${PG_NAME} \
		-v ON_ERROR_STOP=1 \
		-c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

clean-all:
	@printf "\033[35m*** Stopping all containers…\033[0m\n"
	@docker stop $$(docker ps -q) 2>/dev/null || true

	@printf "\033[35m*** Pruning all containers, images, networks, volumes, and caches…\033[0m\n"
	@docker system prune -af --volumes

b: tailwindcss minify-js
	@printf "\033[35m*** Rebuilding blue only…\033[0m\n"
	@docker compose up -d --build --no-deps ${APP_NAME}-blue

tailwindcss:
	@printf "\033[36m*** Building Tailwind CSS…\033[0m\n"
	@cd views && tailwindcss -i ./input.css -o ./static/css/main.css --minify

docker-clean:
	@printf "\033[36m*** Cleaning Docker cache…\033[0m\n"
	@docker builder prune --all --force

minify-js:
	@printf "\033[36m*** Minify js…\033[0m\n"
	@for src in views/static/js/*.src.js; do \
		out="$${src%.src.js}.js"; \
		printf " → $${src##*/} → $${out##*/}\n"; \
		closure-compiler \
		  --compilation_level ADVANCED \
		  --js="$$src" \
		  --js_output_file="$$out"; \
	done
