#!/usr/bin/env bash

if [ ! "$(docker ps -q -f name=shared-sqldb)" ]; then
    if [ "$(docker ps -aq -f status=exited -f name=shared-sqldb)" ]; then
        docker start shared-sqldb
    else
        docker-compose -f shared-sqldb.yml up -d
    fi
fi