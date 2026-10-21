#!/bin/sh
# Exit immediately if a command exits with a non-zero status
set -e

DB_ENGINE=sqlite3

echo "migrating the DBs to the most recent version available..."
goose -dir services/catalog/internal/datasource/sqlite/migrations $DB_ENGINE ./data/catalog.db up
echo "migration complete \n"

echo "compiling all packages and dependencies..."
go build ./...
echo "compiled \n"


echo "testing all packages"
go build ./...
echo "tested \n"

echo "creating binaries..."
go build -o out/catalog services/catalog/cmd/server/main.go
echo "created binaries: output=out/* \n"

echo "executing all the binaries..."
./out/catalog &

echo "enter any key to top the processes _\n"
read _
pkill -15 catalog
echo "stopped all the processes gracefully"
