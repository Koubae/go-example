sync:
	git pull; git status; git add .; git commit -a -m "Sync changes"; git push; printf "\n\n🔎 Checking Sync Status.... 🪄\n"; git status;


# --------------------------
# Init
# --------------------------
init: .install-deps

.install-deps:
	@go mod tidy
	@go work sync

run-playground:
	go run scripts/playground/playground.go
