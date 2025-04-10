#!/bin/bash
CONTAINER_NAME="wadsworth-db"
SSH_PORT=2222
DB_PORT=8000

sudo docker run -d --name $CONTAINER_NAME -p $SSH_PORT:22 -p $DB_PORT:27017 -v $(pwd)/mongod.conf:/etc/mongod.conf:ro -v wadsworth-db-data:/data/db wadsworth-db-image
