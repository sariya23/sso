run:
	go run cmd/sso/main.go --config=./config/local.yaml

migrate:
	go run ./cmd/migrator --migrations-path=./migrations

test_migrate:
	 go run ./cmd/migrator/main.go --migrations-path=./tests/migrations --migrations-table=migrations_test
