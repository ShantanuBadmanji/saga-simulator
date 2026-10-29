#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

DB_ENGINE=sqlite3

echo "migrating the DBs to the most recent version available..."
goose -dir services/catalog/internal/datasource/sqlite/migrations $DB_ENGINE ./data/catalog.db up
echo "migration complete."
echo ""

echo "compiling all packages and dependencies..."
go build ./...
echo "compiled."
echo ""


echo "testing all packages..."
go test ./...
echo "tested."
echo ""

echo "creating binaries..."
go build -o out/catalog services/catalog/cmd/server/main.go
echo "created binaries: output=out/*"
echo ""

echo "executing all the binaries..."
./out/catalog > >(sed 's/^/[catalog] /') 2>&1 &
catalog_pid=$!



cleanup() {
    trap '' INT TERM

	echo ""
    echo "Stopping services..."

    kill -TERM "$catalog_pid"	2>/dev/null || true

    wait "$catalog_pid" 2>/dev/null || true

    echo "All services stopped."
	echo ""

}

trap cleanup INT TERM

# Basic startup check, not a readiness/health check.
sleep 1

if ! err=$(kill -0 "$catalog_pid" 2>&1); then
    printf 'Catalog crashed on startup. Error: "%s"\n' "$err" >&2
    cleanup
    exit 1
fi

echo "Catalog are running."
echo "Press Ctrl+C to stop the services."
echo ""

# keep this shell alive while waiting for both services to exit.
wait "$catalog_pid" || true