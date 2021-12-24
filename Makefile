POSTGRES_PASSWORD ?= hheld-pwd
POSTGRES_USER = hheld
DB_PORT ?= 5432

start-db:
	docker run -d \
		--restart always \
		--name weatherdb \
		-p $(DB_PORT):5432 \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_DB=weather \
		-e TIMESCALEDB_TELEMETRY=off \
		-v weatherdb:/var/lib/postgresql/data \
		timescale/timescaledb:2.5.1-pg14

stop-db:
	docker stop weatherdb

rm-db: stop-db
	docker rm weatherdb