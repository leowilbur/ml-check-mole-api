build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/molepatrol

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: update-jwk
update-jwk:
	echo "package auth\n\nvar jwkKey = \`$$(curl https://cognito-idp.ap-southeast-2.amazonaws.com/ap-southeast-2_gfSuuHw6e/.well-known/jwks.json)\`" > pkg/auth/jwk_key.go

.PHONY: update-docs
update-docs:
	swagger-codegen generate -i swagger.yaml -l swagger -o docs/api
	echo "package api\n\nvar JSON = \`$$(cat docs/api/swagger.json)\`" > docs/api/swagger.go
	sed -i "s|https://petstore.swagger.io/v2/swagger.json|./swagger.json|g" ./node_modules/swagger-ui-dist/index.html
	statik -src=./node_modules/swagger-ui-dist -p swagger -dest ./docs -f

.PHONY: migrate-dev
migrate-dev:
	migrate -source file://./migrations -database postgres://molepatrol:Q637XUy1oUNdCgX1@molepatrol.czaaedbzmswz.ap-southeast-2.rds.amazonaws.com/molepatrol-dev up

.PHONY: migrate-staging
migrate-staging:
	migrate -source file://./migrations -database postgres://molepatrol:Q637XUy1oUNdCgX1@molepatrol.czaaedbzmswz.ap-southeast-2.rds.amazonaws.com/molepatrol-staging up

.PHONY: migrate-prod
migrate-prod:
	migrate -source file://./migrations -database postgres://molepatrol:Q637XUy1oUNdCgX1@molepatrol.czaaedbzmswz.ap-southeast-2.rds.amazonaws.com/molepatrol-prod up

.PHONY: deploy-dev
deploy-dev: clean build
	sls deploy --verbose --stage=dev

.PHONY: deploy-staging
deploy-staging: clean build
	sls deploy --verbose --stage=staging

.PHONY: deploy-prod
deploy-prod: clean build
	sls deploy --verbose --stage=prod
