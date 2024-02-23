# `/api`

OpenAPI/Swagger specs, JSON schema files, protocol definition files.

## Generate openapi.json file
swagger-cli bundle -o ./api/openapi.json ./api/openapi.yml

## Generate models.gen.go code structs
oapi-codegen --config ./api/oapiconfig.yml ./api/openapi.json
