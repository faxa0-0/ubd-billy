DB_NAME=tempi
DB_USER=postgres
DB_PASSWORD=faxa
DB_HOST=localhost
DB_PORT=5432

run: build
	@export JWT_SECRET=secret && ./bin/billy.exe

tidy:
	@go mod tidy

build:
	@go build -o bin/billy.exe main.go

rebuild:
	@echo "Make: Checking dependencies..."
	@go mod tidy
	@echo "Make: Forming database..."
	@export PGPASSWORD=$(DB_PASSWORD) && psql -h $(DB_HOST) -U $(DB_USER) -p $(DB_PORT) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@export PGPASSWORD=$(DB_PASSWORD) && psql -h $(DB_HOST) -U $(DB_USER) -p $(DB_PORT) -d postgres -c "CREATE DATABASE $(DB_NAME);"
	@echo "Make: Building executable..."
	@go build -o bin/billy.exe main.go
	@echo "Make: Executing...."
	@export JWT_SECRET=secret && ./bin/billy.exe

clean:
	@rm -rf bin