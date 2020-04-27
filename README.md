# Simple JSON-RPC 2.0 API w/ PostgreSQL

* using Docker Compose

## Prerequisites

* docker & docker-compose configured for current user
* golang 1.13 and above

## How to start?

1. Run in the terminal: `sh scripts/run-db.sh`
2. Make **.env** file using `cp .env_example .env`
3. You should configure **.env**, check is all variables are correct
4. Start program using: `sh scripts/start.sh`
