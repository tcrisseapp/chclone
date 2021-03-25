up:
	@docker-compose -f ./.docker/docker-compose.yaml build
	@docker-compose -f ./.docker/docker-compose.yaml up

down: 
	@docker-compose -f ./.docker/docker-compose.yaml down
