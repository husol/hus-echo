#!/usr/bin/env bash

if [ ! "$(docker ps -q -f name=shared-dynamodb)" ]; then
    if [ "$(docker ps -aq -f status=exited -f name=shared-dynamodb)" ]; then
        docker start shared-dynamodb
    else
        docker-compose -f shared-dynamodb.yml up -d
    fi
fi