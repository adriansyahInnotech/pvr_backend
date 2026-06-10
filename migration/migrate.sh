#!/bin/bash

set -e

# lokasi folder migration
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ENV_FILE="$SCRIPT_DIR/.env"
MIGRATION_PATH="$SCRIPT_DIR/sql"

echo "📂 Script dir: $SCRIPT_DIR"
echo "📦 Migration path: $MIGRATION_PATH"

# load env
if [ -f "$ENV_FILE" ]; then
  echo "📄 Loading env from $ENV_FILE"
  export $(grep -v '^#' "$ENV_FILE" | xargs)
else
  echo "❌ .env file not found in migration folder"
  exit 1
fi

# cek DATABASE_URL
if [ -z "$DATABASE_URL" ]; then
  echo "❌ DATABASE_URL not set in migration/.env"
  exit 1
fi

echo "🔗 DATABASE_URL loaded"

COMMAND=$1
NAME=$2

case "$COMMAND" in

  up)
    echo "🚀 Running migrations..."
    migrate -path "$MIGRATION_PATH" -database "$DATABASE_URL" up
    ;;

  down)
    echo "⚠️ Rolling back last migration..."
    migrate -path "$MIGRATION_PATH" -database "$DATABASE_URL" down 1
    ;;

  drop)
    echo "🧨 Dropping all database tables..."
    migrate -path "$MIGRATION_PATH" -database "$DATABASE_URL" drop -f
    ;;

  reset)
    echo "🧨 Reset database..."
    migrate -path "$MIGRATION_PATH" -database "$DATABASE_URL" drop -f
    migrate -path "$MIGRATION_PATH" -database "$DATABASE_URL" up
    ;;

  version)
    echo "📊 Migration version:"
    migrate -path "$MIGRATION_PATH" -database "$DATABASE_URL" version
    ;;

  create)
    if [ -z "$NAME" ]; then
      echo "❌ Migration name required"
      echo "Usage: ./migrate.sh create migration_name"
      exit 1
    fi

    migrate create -ext sql -dir "$MIGRATION_PATH" -seq "$NAME"
    echo "✅ Migration created: $NAME"
    ;;

  force)
    if [ -z "$NAME" ]; then
      echo "❌ Version required"
      echo "Usage: ./migrate.sh force VERSION"
      exit 1
    fi

    echo "⚠️ Forcing migration version to $NAME"
    migrate -path "$MIGRATION_PATH" -database "$DATABASE_URL" force "$NAME"
    ;;

  *)
    echo ""
    echo "Usage:"
    echo "./migrate.sh up"
    echo "./migrate.sh down"
    echo "./migrate.sh drop"
    echo "./migrate.sh reset"
    echo "./migrate.sh version"
    echo "./migrate.sh create migration_name"
    echo "./migrate.sh force VERSION"
    echo ""
    ;;
esac