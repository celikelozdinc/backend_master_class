DB_NAME := automation
USR := admin
PASSWD := nimda
IMAGE := postgres:12-alpine
CONTAINER := postgres12
SQLC := sqlc
PKGS := ./...

.PHONY: start-database
start-database:
	docker run --name $(CONTAINER) \
		-e POSTGRES_USER=$(USR) -e POSTGRES_PASSWORD=$(PASSWD) -e POSTGRES_DB=$(DB_NAME) \
		-p 5432:5432 -d $(IMAGE)

.PHONY : stop-database
stop-database:
	@docker stop $(CONTAINER)
	@docker rm $(CONTAINER)

.PHONY : sqlc
sqlc:
	$(info Generating/Regenerating Golang code for SQL queries...)
	@which $(SQLC) || (echo $(ERR) ; exit 1;)
	@echo "sqlc : $(shell which $(SQLC))"
	@sqlc generate

.PHONY : coverage
coverage:
	@go test -v -cover -count 1 $(PKGS)

.PHONY: createdb
createdb:
	@docker exec -it postgres12 createdb --username=$(USR) --owner=$(USR) $(DB_NAME)


.PHONY : dropdb
dropdb:
	@docker exec -it postgres12 dropdb $(DB_NAME)


.PHONY: create-schema
create-schema:
	$(info Creating/Recreating database schema...)
	@migrate -path sql/migration -database "postgresql://admin:nimda@localhost:5432/automation?sslmode=disable" -verbose up


.PHONY: delete-schema
delete-schema:
	$(info Deleting database schema...)
	@migrate -path sql/migration -database "postgresql://admin:nimda@localhost:5432/automation?sslmode=disable" -verbose down

.PHONY : server
server:
	$(info Starting server...)
	@go run ./main.go

.PHONY : mock
mock: 
	$(info Generating mock database artifacts...)
	@mockgen -package mockdb -destination internal/db/mock/transaction.go code.siemens.com/ozdinc.celikel/backend_master_vlass/internal/db Tx

.PHONY : start-database stop-database sqlc coverage createdb dropdb create-schema delete-schema server mock