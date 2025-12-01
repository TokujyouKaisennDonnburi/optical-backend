#!/bin/bash
# migrate.sh

# Load environment variables from .env file
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

sql-migrate up
sql-migrate status

mc alias set myminio http://${MINIO_ENDPOINT} $MINIO_ACCESS_KEY_ID $MINIO_SECRET_ACCESS_KEY
mc mb myminio/$MINIO_IMAGE_BUCKET_NAME
mc anonymous set download myminio/$MINIO_IMAGE_BUCKET_NAME
