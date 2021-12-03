.PHONY: up
up:
	@docker-compose up --build

.PHONY: down
down:
	@docker-compose down -v

.PHONY: db-dev
db-dev:
	@docker exec -it backend_db_1 bash

.PHONY: app-dev
app-dev:	
	@docker exec -it backend_app_1 bash