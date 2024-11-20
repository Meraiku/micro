build_user:
	@go build -o ./user/.bin/app ./user/cmd/app

build_ws:
	@go build -o ./websocket/.bin/app ./websocket/cmd/app

build_notif:
	@go build -o ./notification/.bin/app ./notification/cmd/app


docker:
	@docker compose -f docker-compose.local.yaml up --build -d

stop:
	@docker compose -f docker-compose.local.yaml stop
