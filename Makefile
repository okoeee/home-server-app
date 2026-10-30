.DEFAULT_GOAL := dev

.PHONY: dev generate-sqlc generate-openapi

# 並列実行を指定しても、生成が完了してからアプリをビルドする。
.NOTPARALLEL:

dev: generate-sqlc generate-openapi
	docker compose up -d --build web

generate-sqlc:
	docker compose run --rm sqlc generate

generate-openapi:
	npm --prefix openapi ci
	npm --prefix openapi run bundle
	npm --prefix openapi run generate:typescript
	cd apps/api && go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.1 -config oapi-codegen.yaml ../../openapi/dist/openapi.yaml
