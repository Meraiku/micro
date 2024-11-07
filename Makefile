build_user:
	@go build -o ./user/.bin/app ./user/cmd/app

build_ws:
	@go build -o ./websocket/.bin/app ./websocket/cmd/app


docker:
	@docker compose -f docker-compose.local.yaml up

