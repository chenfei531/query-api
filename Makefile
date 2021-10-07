build:
	go build -o bin/rql_user main/rql_user.go
	go build -o bin/rql_agent main/rql_agent.go
	go build -o bin/graphql_user main/graphql_user.go
	go build -o bin/graphql_agent main/graphql_agent.go
	go build -o bin/db_init db/db_init.go