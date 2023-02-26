DB_URL=postgresql://postgres:postgres@localhost:5432/imchat?sslmode=disable&timeZone=Asia/Shanghai
SQL_DIR=/opt/qianxia/imchat

server:
	go run main.go

postgres:
	docker run -d -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root --restart=always --name postgres postgres:14-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root imchat

migrateup:
	migrate -path db/pg/migration -database "${DB_URL}" -verbose up

migrateup1:
	migrate -path db/pg/migration -database "${DB_URL}" -verbose up 1

migratedown:
	migrate -path db/pg/migration -database "${DB_URL}" -verbose down

migratedown1:
	migrate -path db/pg/migration -database "${DB_URL}" -verbose down 1

sqlc:
	docker run --rm -v $(SQL_DIR):/src -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/pg/mock/store.go IMChat/db/pg/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=imchat \
    proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: server postgres createdb migrateup migrateup1 migratedown migratedown1 sqlc test mock proto evans