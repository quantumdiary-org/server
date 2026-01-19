# Build
build:
	go build -o bin/server ./api/cmd/server

# Run
run: build
	./bin/server --config config/dev.yaml

# Database
setup-db:
	psql -U postgres -c "CREATE DATABASE netschool_proxy;"
	psql -U postgres -c "CREATE USER proxy_user WITH PASSWORD 'dev_password';"
	psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE netschool_proxy TO proxy_user;"
	go run ./api/cmd/migrate up

# Tests
test:
	go test -v ./api/... -count=1

test-integration:
	go test -v -tags=integration ./api/... -count=1

# Linting
lint:
	golangci-lint run ./api/...

# Swagger
generate-swagger:
	swag init -g api/cmd/server/main.go -o api/pkg/docs

# Cleanup
clean:
	rm -rf bin/
	rm -rf coverage.out
	rm -rf logs/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run with race detection
race:
	go run -race ./api/cmd/server/main.go --config config/dev.yaml

# Build for different platforms
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/server-linux ./api/cmd/server

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/server.exe ./api/cmd/server

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/server-mac ./api/cmd/server

.PHONY: build run setup-db test test-integration lint generate-swagger clean deps race build-linux build-windows build-macos