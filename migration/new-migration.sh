#!/bin/sh
TIMESTAMP=$(date +%s)
MIGRATION_NAME=$1
[ -z "$MIGRATION_NAME" ] && echo "Migration name is required" && exit 1
cp ./0-migration-template ./$TIMESTAMP-$MIGRATION_NAME.sql
