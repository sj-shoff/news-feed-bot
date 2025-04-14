#!/bin/sh

host_port=$1
host=${host_port%:*}
port=${host_port#*:}

shift

until PGPASSWORD=$DB_PASSWORD psql -h "$host" -p "$port" -U "$DB_USER" -d "$DB_NAME" -c '\q'; do
  >&2 echo "PostgreSQL $host:$port is unavailable - sleeping"
  sleep 2
done

>&2 echo "PostgreSQL is ready - executing command"
exec "$@"