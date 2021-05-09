.PHONY: swagger db-check-local db-create-local db-run-local db-new db-migrate db-unmigrate

swagger-check :
	which swagger || (echo "need to install swagger!" && exit 1)


swagger : swagger-check
	swagger generate spec -o ./swagger.yml --scan-models

run : monolith
	./monolith

monolith : install
	go build

install : go.mod go.sum
	go get

db-check-local :
	which migrate || (echo "need to install migrate!" && exit 1)
	which psql || (echo "need to install postgresql!" && exit 1)
	which docker || (echo "need to install docker!" && exit 1)

db-create-local : db-check-local
	psql -h localhost -U postgres -w -c "create database turnkey;"

db-run-local : db-check-local
	psql -h localhost -U postgres -w turnkey -c "$(CMD)"

db-start-local : db-check-local
	sudo docker run --rm -ti --network host -e POSTGRES_PASSWORD=password postgres

db-new :
	migrate create -ext sql -dir db/migrations -seq $(SEQ_NAME)

db-migrate :
	migrate -database $(DB_CONN) -path db/migrations up

db-unmigrate :
	migrate -database $(DB_CONN) -path db/migrations down
