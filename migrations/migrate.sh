#!/bin/bash

docker exec scylla cqlsh -f ./migrations/V1_initial_database.sql
