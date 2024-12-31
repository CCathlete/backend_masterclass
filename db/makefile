connect:
	docker-compose exec db psql -U ${DB_USER} -d ${DB_NAME}

# dbfile:
# 	docker cp ~/Repos/backend_masterclass/sql/bank_schema.sql backend_masterclass_db_1:/; docker-compose exec db psql -U ${DB_USER} -d ${DB_NAME} -f bank_schema.sql

migrateup:
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose up

migrateup_1: # Up one migration.
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose up 1

migratedown:
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose down

migratedown_1: # Down one migration back.
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ../...

migraterestart:
	make migratedown; make migrateup

server:
	go run ../cmd/main.go

mock:
	mockgen -destination ./mock/${file_name} -package ${pkg_name}  /home/ccat/Repos/backend_masterclass/db/sqlc ${interfaces}

proto:
	rm -f ../rpc/*.go
	protoc --proto_path=../proto --go_out=../rpc --go_opt=paths=source_relative \
	--go-grpc_out=../rpc --go-grpc_opt=paths=source_relative \
	../proto/*.proto

evans:
	evans --host=localhost --port=9090 --reflection repl --package pb --service SimpleBank

# Command aliasing is considered a "phony target" so it's possible to run it repeatedly.
.PHONY: connect migrateup migratedown sqlc test migraterestart server mock migratedown_1 migrateup_1 proto evans #dbfile
