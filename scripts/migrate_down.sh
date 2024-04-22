#!/bin/bash

if [ -f .env ]; then
    source .env
fi

service=$1
service_upper=${service^^}
url_var="${service_upper}_DB_URL"
cmd="echo \$${url_var}"
db_url=$(eval $cmd)

cd services/$service/sql/schema
goose turso $db_url down
