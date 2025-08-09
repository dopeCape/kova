.PHONY: test clean all build 

dev:
	air --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"

generate:
	sqlc generate

migrate-up: 
	@echo "${YELLOW}Running migrations up...${NC}"
	migrate -path internal/store/migrations -database "postgresql://admin:password123@localhost:5432/mydb?sslmode=disable" up
	@echo "${GREEN}Migrations completed!${NC}"

migrate-down: 
	@echo "${YELLOW}Running migrations down...${NC}"
	migrate -path internal/store/migrations -database "postgresql://admin:password123@localhost:5432/mydb?sslmode=disable" down 1
	@echo "${GREEN}Migration rolled back!${NC}"
