#!/bin/sh

set -e

echo "run db migraton"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up
exec "$@"