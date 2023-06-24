#!/bin/bash

docker pull postgres

docker run --rm \
    --name postgres_test \
    --network host \
    -e "POSTGRES_PASSWORD=postgres" \
    -e "POSTGRES_DB=praktikum" \
    -e GITHUB_ACTIONS=true \
    -e CI=true \
    postgres:latest