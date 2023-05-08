# code-challenge

## Description
Code-challenge is a simple containerized service written in Golang

## Features
- Auth with JWT token
- Web forms for login and upload images

## Component
- HTTP Server
- Postgres as a storage

## Environment
- Docker

## Development Guide
1. Start the postgres database
    ```shell
    # host: localhost
    # port: 5432
    # username: postgres
    # password: changeme
    # database name: postgres
    
    docker-compose start postgres
    ```
2. App configuration
    ```shell
    server_port: The server port, default 8080
    dsn: The data source name, eg: "host=localhost user=postgres password=changeme dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    jwt_signing_key: The secret key
    jwt_expiration: The jwt token expiration, default: 3600
    ```
3. Start app
    ```shell
    make start
    #or
    go run ./cmd/server/main.go
    ```

## Deployment:
1. Create `.env` file
2. Populate all configurations as ENV variables
   ```shell
   # cat .env
   
   APP_ENVIRONMENT=dev
   APP_DSN="host=postgres user=postgres password=changeme dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
   APP_JWT_SIGNING_KEY=signing_key_dev
   APP_JWT_EXPIRATION=3600
   ```
3. Start docker compose
- DEV: `docker-compose -f ./docker-compose.yml up`
