local_run:
	go run cmd/sso/main.go --config=./config/local.yaml

migrate:
	go run ./cmd/migrator --migrations-path=./migrations