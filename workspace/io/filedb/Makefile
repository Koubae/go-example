.PHONY: run build stop tests

run:
	@go run ./cmd/filedb/


# ============================
# 	Tests
# ============================

### Using locally installed Go
test-unit-local:
	go test -v ./... -coverprofile=coverage.out && go tool cover -o coverage.html -html=coverage.out

test-specific:
ifndef TEST
	@echo "Please provide a test pattern using TEST=<pattern>"
	@echo "Example: make test-specific TEST=TestGetEnv/string_tests"
	@echo "make test-specific TEST=TestGetEnv"
	@echo "\nAvailable test patterns:"
	@go test ./... -v -list=. | grep "^Test"
else
	go test ./... -v -run $(TEST)
endif
