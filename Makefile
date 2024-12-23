.PHONY: up down rebuild bash start

up:
	@docker compose up -d

down:
	@docker compose down

rebuild:
	@docker compose build --no-cache

bash:
	@docker exec -it pow bash

start:
	@docker compose up -d && docker exec -it pow bash