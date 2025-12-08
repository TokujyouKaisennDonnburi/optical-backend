#!/bin/bash

docker exec -it optical_postgres /bin/bash -c "psql -U optical_user optical"
