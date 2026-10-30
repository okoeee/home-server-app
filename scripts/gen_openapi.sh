#!/usr/bin/env sh

set -eu

npm --prefix openapi run bundle
npm --prefix openapi run generate:typescript

cd apps/api
mkdir -p openapi
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.1 \
  -config oapi-codegen.yaml ../../openapi/dist/openapi.yaml
