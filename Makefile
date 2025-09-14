# --------------------------
# Init
# --------------------------
init: .install-deps

.install-deps:
	@go mod tidy
	@go work sync

run-playground:
	go run scripts/playground/playground.go
