build:
	go build -o bin/run main/main.go
	go build -o bin/db_init db/db_init.go