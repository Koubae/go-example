
# ============================
# 	Init
# ============================
init: install-deps update-env-file

install-deps:
	go mod tidy
	go install gotest.tools/gotestsum@latest

update-env-file:
	@echo 'Updating .env from .env.example 🖋️...'
	@cp .env.example .env
	@echo '.env Updated ✨'


# ============================
# 	Run
# ============================

# //////////////////////
# 	local
# //////////////////////
run-cli-wallet:
	go run ./cmd/cliwallet $(ARGS)

run-wallet-server:
	go run ./cmd/server $(ARGS)




# //////////////////////
# 	docker
# //////////////////////

# ······················
#		DB -- PostGreSQL
# ······················
db-up:
	@docker compose up -d db-postgres
db-down:
	@docker compose down db-postgres
db-down-clean-up:
	@docker compose down -v db-postgres db-postgres

db-shell:
	docker compose exec \
		-e PGPASSWORD='admin' \
		db-postgres \
		psql -U admin -d cli_wallet


# //////////////////////
# 	Tests
# //////////////////////
COVERAGE_THRESHOLD ?= 80

TestsPaths := ./internal/...
CoverPaths := ./internal/...

comma := ,
empty :=
space := $(empty) $(empty)
CoverPkg := $(subst $(space),$(comma),$(strip $(CoverPaths)))
CoverageProfile := coverage.out
CoverageHTML := coverage.html

.PHONY: test-unit
test-unit:
	gotestsum --format pkgname -- $(TestsPaths)

.PHONY: test-coverage
test-coverage:
	gotestsum --format pkgname -- -coverprofile=$(CoverageProfile) -covermode=atomic -coverpkg=$(CoverPkg) $(TestsPaths)
	@go tool cover -html=$(CoverageProfile) -o $(CoverageHTML)
	@echo "Coverage report: $(CoverageHTML)"
	@coverage=$$(go tool cover -func=$(CoverageProfile) | awk '/^total:/ { gsub(/%/, "", $$NF); print $$NF }'); \
	echo "Total coverage: $$coverage% (threshold: $(COVERAGE_THRESHOLD)%)"; \
	awk -v c="$$coverage" -v t="$(COVERAGE_THRESHOLD)" 'BEGIN { if ((c+0) < (t+0)) { printf "Coverage %.1f%% is below threshold %s%%\n", c, t; exit 1 } }'
