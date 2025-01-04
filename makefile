server:
	go run cmd/main.go

mock:
	mockgen -destination ./db/mock/${file_name} -package ${pkg_name}  /home/ccat/Repos/backend_masterclass/db/sqlc ${interfaces}

proto:
	rm -f ./rpc/*.go ./views/openapi/*.json
	protoc \
	--proto_path=proto \
	--proto_path=/home/ccat/go/pkg/mod/\
	github.com/grpc-ecosystem/grpc-gateway/v2@v2.25.1 \
	--go_out=controllers/protoc --go_opt=paths=source_relative \
	--go-grpc_out=controllers/protoc --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=rpc --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=views/openapi --openapiv2_opt=logtostderr=true \
	--openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto
	# statik -src=views/openapi -dest=views

evans:
	evans --host=localhost --port=9090 --reflection repl --package pb --service SimpleBank

.PHONY: server mock proto evans