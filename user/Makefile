user_proto:
	@echo "Compiling user proto"
	@mkdir -p pkg/user_v1
	@protoc --proto_path api/user_v1 \
		--go_out=pkg/user_v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
		api/user_v1/user.proto

auth_proto:
	@echo "Compiling auth proto"
	@mkdir -p pkg/auth_v1
	@protoc --proto_path api/auth_v1 \
		--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
		api/auth_v1/auth.proto

proto:
	@make user_proto
	@make auth_proto

build:
	@go build -o ./.bin/app ./cmd/app

run:build
	@./.bin/app

