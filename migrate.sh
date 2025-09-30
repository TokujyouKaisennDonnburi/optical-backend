#!/bin/bash
# migrate.sh

# Load environment variables from .env file
if [ -f .env ]; then
    export $(grep -v '^#' .env | grep -v '^$' | xargs)
fi

sql-migrate up
