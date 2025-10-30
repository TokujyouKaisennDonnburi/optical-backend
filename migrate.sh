#!/bin/bash
# migrate.sh

# Load environment variables from .env file
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

sql-migrate up
