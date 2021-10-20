build:
	go build -o bin/query_worker main/query_worker.go
	go build -o bin/graphql_query main/graphql_query.go
	go build -o bin/db_init db/db_init.go
