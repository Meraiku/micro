proto:
	@mkdir -p pkg/user_v1
	@protoc --proto_path api/user_v1 \
		--go_out=pkg/user_v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
		api/user_v1/user.proto

build:
	@go build -o ./.bin/app ./cmd/app

run:build
	@./.bin/app

