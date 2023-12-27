ATLAS_URL := "postgres://rankit:rankit@localhost:5432/rankit?sslmode=disable"
ATLAS_DEV_URL := "docker://postgres/16"
ATLAS_SCHEMA := "file://schema.sql"
GOBIN := $(shell go env GOPATH)/bin

setup:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.24.0
	brew install ariga/tap/atlas

generate: generate-sql

generate-sql:
	$(GOBIN)/sqlc generate -f .sqlc.yml
	$(GOBIN)/sqlc vet -f .sqlc.yml

migrate:
	atlas schema apply --url $(ATLAS_URL) --to $(ATLAS_SCHEMA) --dev-url $(ATLAS_DEV_URL)

shell:
	psql postgres://rankit:rankit@localhost:5432/rankit

watch:
	watch -n 3 make generate