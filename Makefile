db_up:
	migrate -path ./migrate -database 'postgres://postgres:123@localhost:5436/postgres?sslmode=disable' up

db_down:
	migrate -path ./migrate -database 'postgres://postgres:123@localhost:5436/postgres?sslmode=disable' down

build:
	go build -o wbL0 -v ./cmd/main.go

run: build
	./wbL0

pub:
	go run ./nats-streaming/script/script.go

.DEFAULT_GOAL := run