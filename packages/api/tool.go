package api

//go:generate sh bundle.sh
//go:generate go tool oapi-codegen -config cfg.yaml api.bundled.yaml
//go:generate echo "Generating Echo server stubs with OpenAPI Generator..."
//go:generate sh -c "if [ ! -f ../../tools/openapi-generator-cli.jar ]; then ../../tools/openapi-generator.sh version; fi"
//go:generate ../../tools/openapi-generator.sh generate -i api.bundled.yaml -g go-echo-server -o ./gen --package-name api
