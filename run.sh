#!/bin/bash

# Create the postgres-data directory if it does not exist
if [ ! -d "./postgres-data" ]; then
  mkdir -p ./postgres-data
  echo "postgres-data directory created."
else
  echo "postgres-data directory already exists."
fi

# Set read and write permissions for the postgres-data directory
chmod -R 775 ./postgres-data
echo "Permissions set for postgres-data directory."

# Function to find and kill the process using a given port
kill_process_on_port() {
  local PORT=$1
  local PROCESS_ID=$(lsof -t -i:$PORT)
  if [ ! -z "$PROCESS_ID" ]; then
    echo "Killing process $PROCESS_ID using port $PORT"
    kill -9 $PROCESS_ID
  fi
}

# Function to stop and remove containers using a given port
stop_containers_on_port() {
  local PORT=$1
  local RUNNING_CONTAINERS=$(docker ps -q --filter expose=$PORT)
  if [ ! -z "$RUNNING_CONTAINERS" ]; then
    echo "Stopping containers using port $PORT"
    docker stop $RUNNING_CONTAINERS
    docker rm $RUNNING_CONTAINERS
  fi
}

# Kill processes using ports 5432 and 8087
kill_process_on_port 5432
kill_process_on_port 8087

# Stop and remove containers using ports 5432 and 8087
stop_containers_on_port 5432
stop_containers_on_port 8087

# Stop and remove any existing containers
docker stop spy_cat_db_1 spy_cat_app 2>/dev/null || true
docker rm spy_cat_db_1 spy_cat_app 2>/dev/null || true

# Stop and remove any existing containers, networks, and volumes
docker-compose down -v --remove-orphans

# Build and start containers
docker-compose up --build -d

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
DB_CONTAINER=$(docker-compose ps -q db)
while ! docker exec -i $DB_CONTAINER pg_isready -U spycat -d spycatdb -h db; do
  sleep 1
done

# Run database migrations
echo "Running database migrations..."
docker exec -i $DB_CONTAINER psql -U spycat -d spycatdb < ./sql/create_tables.sql

# Wait additional time to ensure database is fully ready
sleep 10

echo "Application and database are up and running. You can now send API requests."
