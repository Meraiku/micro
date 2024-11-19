FROM golang:1.23-alpine AS builder

WORKDIR /usr/src/app

RUN apk add --no-cache make

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./websocket/ ./websocket
COPY ./pkg/ ./pkg
COPY ./user/pkg/ ./user/pkg
COPY ./Makefile ./
RUN make build_ws



FROM alpine AS runner

COPY --from=builder /usr/src/app/websocket/.bin/app /bin/app

CMD ["/bin/app"]
