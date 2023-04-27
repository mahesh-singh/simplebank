#!/bin/sh

set -e

echo "run db migraton"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up
exec "$@"