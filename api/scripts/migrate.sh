#!/usr/bin/env bash

source ./.env

migrate -source file://db/migrations -database "${DATABASE_URL}" $@
