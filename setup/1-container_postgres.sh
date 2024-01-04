#!/bin/bash

# Permissão antes de rodar:
# chmod +x container_postgres.sh

docker run -d \
  --name skills \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=postgres \
  -e POSTGRES_HOST=172.17.0.2 \
  -p 5432:5432 \
  postgres:latest
