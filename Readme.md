# go-openapi-playground

OpenAPI 3 example with code generation from `spec.yaml` and integration into an echo http server. Tested with go v1.17.

## Prepare

```bash
go mod tidy
```

## Build

```bash
go generate ./...
# generates *.gen.go files in ./api
go build
```

## Run

```bash
./go-openapi-playground
```
