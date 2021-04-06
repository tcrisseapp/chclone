up:
	@docker-compose -f docker-compose.yaml up -d --build --force-recreate

down: 
	@docker-compose -f docker-compose.yaml down
	@docker volume rm docker_db

