.PHONY: up
up:
	@docker-compose up --build

.PHONY: stop
stop:
	@docker compose stop

.PHONY: down
down:
	@docker-compose down -v

.PHONY: logs
logs:
	@docker compose logs -f

.PHONY: db-dev
db-dev:
	@docker exec -it hack_mysql bash

.PHONY: app-dev
app-dev:	
	@docker exec -it backend_app_1 bash