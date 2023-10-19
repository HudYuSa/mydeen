generate_migration:
	migrate create -ext sql -dir db/migration -seq $(name)