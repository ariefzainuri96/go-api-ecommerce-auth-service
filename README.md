## air-verse for hot reload

1. <https://github.com/air-verse/air>
2. run air with `air -c .air.toml`

## Migration

1. Head to migrate github `https://github.com/golang-migrate/migrate/tree/master/cmd/migrate` -> this link contains cli installation
2. Create migration file:
   migrate create -seq -ext sql -dir ././cmd/migrate/migrations create_users
3. Perform migration:
   migrate -path ./cmd/migrate/migrations -database="postgres://postgres:postgres@localhost:5432/go_api_ecommerce_auth_service?sslmode=disable&connect_timeout=15" up

## Docker

1. running docker container that we specify in docker-compose.yml:
   docker compose up -d

2. stop docker container:
   docker compose down

3. remove docker container with its volumes:
   docker compose down -v

4. if you encounter failed to connect to db, try to stop docker compose and run again, its because the database is not yet created

## .air.toml

1. current working is for linux because we are using docker for running this apps
   bin = "./bin/api"
   cmd = "go build -o ./bin/api ./cmd/api/"

2. if you running locally, change .air.toml line 7-8 to:
   bin = "./bin/api.exe"
   cmd = "go build -o ./bin/ ./cmd/api/"
