.PHONY: run test unit bdd cover fmt vet

run: ## Levanta la aplicación en http://localhost:8080
	go run ./cmd/server

test: vet ## Tests unitarios + escenarios BDD
	go test -race ./...

unit: ## Solo tests unitarios
	go test -race $$(go list ./... | grep -v /features)

bdd: ## Solo escenarios BDD (godog)
	go test ./features/ -v

cover: ## Cobertura total y reporte HTML en coverage.html
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -n 1
	go tool cover -html=coverage.out -o coverage.html

fmt:
	gofmt -w .

vet:
	go vet ./...
