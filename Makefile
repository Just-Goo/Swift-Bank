postgres16:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=swiftsecret -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root Swift_Bank_DB

dropdb:
	docker exec -it postgres16 dropdb Swift_Bank_DB

createmigration:
	migrate create -ext sql -dir database/migrations -seq init_schema

migrateup:
	migrate -path database/migrations -database "postgresql://root:swiftsecret@localhost:5432/Swift_Bank_DB?sslmode=disable" -verbose up

migrateup1:
	migrate -path database/migrations -database "postgresql://root:swiftsecret@localhost:5432/Swift_Bank_DB?sslmode=disable" -verbose up 1

migratedown:
	migrate -path database/migrations -database "postgresql://root:swiftsecret@localhost:5432/Swift_Bank_DB?sslmode=disable" -verbose down

migratedown1:
	migrate -path database/migrations -database "postgresql://root:swiftsecret@localhost:5432/Swift_Bank_DB?sslmode=disable" -verbose down 1

test:
	go clean -testcache && go test -v -cover ./...

run:
	go run main.go

mockrepo:
	mockgen -package mockedproviders -destination mock/repository_provider.go github.com/Just-Goo/Swift_Bank/repository RepositoryProvider

mockservice:
	mockgen -package mockedproviders -destination mock/service_provider.go github.com/Just-Goo/Swift_Bank/service ServiceProvider

.PHONY: postgres16 createdb dropdb migrateup migrateup1 migratedown migratedown1 test run createmigration mockrepo mockservice