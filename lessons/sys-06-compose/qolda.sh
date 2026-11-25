#!/usr/bin/env bash
# compose'siz: xuddi shu ikki container'ni qoʻlda koʻtarish (compose.yaml bilan solishtiring)
set -euo pipefail
source .env
docker network create qolda
docker volume create qolda-pg
docker run -d --name qolda-db --network qolda -v qolda-pg:/var/lib/postgresql \
  -e POSTGRES_USER="$DB_USER" -e POSTGRES_PASSWORD="$DB_PASSWORD" -e POSTGRES_DB="$DB_NAME" postgres:18-alpine
until docker exec qolda-db pg_isready -q -h 127.0.0.1 -U "$DB_USER" -d "$DB_NAME"; do sleep 1; done
docker build -q -t qolda-app .
docker run -d --name qolda-app --network qolda -p 127.0.0.1:8077:8077 \
  -e DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@qolda-db:5432/$DB_NAME?sslmode=disable" qolda-app
# tozalash: docker rm -f qolda-app qolda-db; docker network rm qolda; docker volume rm qolda-pg; docker image rm qolda-app
