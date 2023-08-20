#========================== Docker Compose ==========================

local:
	@echo "Starting local environment"
	@docker-compose -f build/local/docker-compose.yml up --build


#========================== Docker support ==========================


FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f

logs-local:
	docker logs -f $(FILES)
