run:
	go run cmd/mini-jira/main.go

seed:
	go run cmd/seed/main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

lint:
	golangci-lint run

dev:
	air
