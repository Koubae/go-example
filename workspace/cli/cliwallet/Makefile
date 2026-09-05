

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


