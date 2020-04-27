#!/bin/bash
docker-compose -f deployments/docker-compose.yml down
docker-compose -f deployments/docker-compose.yml up
