connect:
	docker-compose exec db psql -U ${DB_USER} -d ${DB_NAME}

# dbfile:
# 	docker cp ~/Repos/backend_masterclass/sql/bank_schema.sql backend_masterclass_db_1:/; docker-compose exec db psql -U ${DB_USER} -d ${DB_NAME} -f bank_schema.sql

migrateup:
	migrate -path db/migrations -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

migraterestart:
	make migratedown; make migrateup

# Command aliasing is considered a "phony target" so it's possible to run it repeatedly.
.PHONY: connect migrateup migratedown sqlc test migraterestart #dbfile

