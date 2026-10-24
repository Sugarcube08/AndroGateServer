#!/bin/sh
# docker-entrypoint.sh

# Immediately exit if any command has a non-zero exit status.
set -e

# If PORT is provided by cloud platforms (e.g. Render) and HTTP__LISTEN is not set, map it automatically
if [ -n "${PORT}" ] && [ -z "${HTTP__LISTEN}" ]; then
  export HTTP__LISTEN="0.0.0.0:${PORT}"
fi

# Execute DB migrations only when the main application is about to run
if [ $# -eq 1 ] && [ "${1:-}" = "/app/app" ]; then
  /app/app db:migrate up
fi

# Execute the main application
exec "$@"
