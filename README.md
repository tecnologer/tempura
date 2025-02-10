# Tempura

Tempura is a tool to manage records from the sensors to store temperature, humidity and liquid levels. It is a simple tool that allows you to store the data in a database and retrieve it later.

## Environment Configuration

To run the project you need to have a `.env` file in the root of the project with the following variables:

```env
TEMPURA_DB_HOST=database // If you are using docker-compose
TEMPURA_DB_PORT=5432
TEMPURA_DB_USERNAME=username
TEMPURA_DB_PASSWORD=password
TEMPURA_DB_NAME=database_name
TEMPURA_DB_SSL_MODE=disable
TEMPURA_DB_CONTAINER=container_db
TEMPURA_API_PORT=80
TEMPURA_LOCAL_API_PORT=8080
TEMPURA_TELEGRAM_BOT_TOKEN=telegram_bot_token
```

## Docker 

To run the project you need to have docker installed in your machine. You can install it from the official website [Docker](https://www.docker.com/).

```bash
docker compose -f docker-compose.yml up --build
```

If you run on ARM architecture you can use the following command:

```bash
docker compose -f docker-compose-rasp.yml up --build
```
