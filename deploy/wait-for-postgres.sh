#!/bin/bash

until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -c '\q'; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

exec "$@"
