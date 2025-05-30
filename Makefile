include .env

build-n-run:
	swag init -d cmd/dengovie,internal/app/dengovie,internal/web,internal/domain
	go run cmd/dengovie/main.go

test-cov:
	go test -coverprofile cover.out  `go list ./internal/... | grep -v ./internal/mocks`

test-env-up:
	touch .env
	docker run -d -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=dengovie -p 5432:5432 --name=dengovie postgres

test-env-down:
	-docker kill dengovie
	docker rm dengovie

db-create:
	@[ "$(NAME)" ] || ( echo '💥 Please use:  make NAME="create_pages" db-create'; exit 1 )
	 goose -dir migrations create "$(NAME)" sql

db-create-test-migration:
	@[ "$(NAME)" ] || ( echo '💥 Please use:  make NAME="create_pages" db-create'; exit 1 )
	 goose -dir test/test_migrations create "$(NAME)" sql

db-up-test-data:
	goose -dir test/test_migrations postgres $(POSTGRES_CONN_STRING) up

db-up:
	goose -dir migrations postgres $(POSTGRES_CONN_STRING) up

binaries:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/vektra/mockery/v2@latest

.PHONY: mocks
mocks:
	mockery