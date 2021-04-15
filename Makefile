up:
	@docker compose -f docker-compose.yaml up -d --build --force-recreate

update:
	@docker compose -f docker-compose.yaml up --build --force-recreate -d room
	@docker compose -f docker-compose.yaml up --build --force-recreate -d postgres

down: 
	@docker compose -f docker-compose.yaml down
	@docker volume rm docker_db

