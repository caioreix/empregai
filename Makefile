#============================ Migrations ============================

force:
	@migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations force 1

version:
	@echo "Migration version:"| tr "\n" " "
	@migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations version

create:
	@migrate create -ext sql -dir migrations -seq -digits 3 $(name)

migrate-up:
	@migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations up 1

migrate-down:
	@migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations down 1

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
