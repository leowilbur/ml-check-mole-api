# ml-check-mole-api

This repository contains an implementation of the CheckMyMole API. All the
endpoints have been described in the [Swagger definition](./swagger.yaml).

## Deploying

In order to deploy this project to AWS Lambda, simply commit the changes to this
repo and the CI process will take care of everything. 

In order to deploy the project manually, run the following:

```bash
ENV="dev"         # or any stage of your choice
make migrate-$ENV # run the database migrations
make deploy-$ENV  # deploy to AWS Lambda
```

## Development

This project uses [`dep`](https://github.com/golang/dep) for dependency management
and [`migrate`](https://github.com/golang-migrate/migrate) for running database
migrations. Additionally [`swagger-codegen`](https://swagger.io/tools/swagger-codegen/)
is required for `swagger.yaml` compilation,
[`statik`](https://github.com/rakyll/statik) is required for resource embedding
and `node`, `npm` and [`serverless`](https://serverless.com/) are required for
deployment to AWS Lambda.

In order to run the tests, you will also need a PostgreSQL database. In order to
make setting it up easier, there's a `docker-compose.yml` file, which lets you
start a fully configured database using a single command:

```bash
docker-compose up -d
```

After the database is running, you can either run the tests:

```bash
# Using ginkgo (more features)
ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --trace --progress

# Or using the standard go test
go test -v ./...
```

Or start the API using the local database:

```bash
export POSTGRES="user=molepatrol password=molepatrol host=localhost dbname=molepatrol sslmode=disable"
migrate -source file://./migrations -database postgres://molepatrol:molepatrol@localhost/molepatrol up
go build -v && ./ml-check-mole-api
```

## Project structure

 - `.circleci` contains all files related to the CI process,
 - `docs` contains all the documentation packages used for resource embedding,
 - `migrations` contains all the SQL migrations,
 - `pkg`:
   - `auth` has all the authentication helpers,
   - `models` contains all the data model declarations and a basic CRUD layer,
   - `rest` contains an implementation of the REST API,
   - `types` contains custom SQL/JSON types,
 - `vendor` is the Go vendor directory,
 - the root directory contains mostly configuration files

## Configuring CORS

All the CORS configuration is available in `serverless.yaml`. The OPTIONS route
is configured by the Serverless software during the deployment.

## Updating documentation

0. Ensure that you have dependencies:
  - [`swagger-codegen`](https://github.com/swagger-api/swagger-codegen)
  - [`statik`](https://github.com/pzduniak/statik)
  - `sed`
  - fetched node modules (`npm install`).
1. Update `swaggger.yaml`.
2. Run `make update-docs`.

## Authenticating against the API

The API expects an access token from the `ap-southeast-2_gfSuuHw6e` AWS Cognito
User Pool to be passed with every single request that requires authentication
in form of an `Authorization: Bearer <token>` header. This token must be acquired
as a result of an OAuth2 flow (ie. it must be an OIDC token). It seems that the
serverside username + password configuration does NOT work with the OIDC endpoints.

For JavaScript, AWS recommends using [Amplify.js](https://github.com/aws-amplify/amplify-js).
In order to test the access token, query the `GET /users/me` endpoint.

Please note that there are no "image upload" endpoints, as an AWS Identity Pool
has been configured. The credentials are as following:

```
// Initialize the Amazon Cognito credentials provider
CognitoCachingCredentialsProvider credentialsProvider = new CognitoCachingCredentialsProvider(
    getApplicationContext(),
    "ap-southeast-2:cd3aaf88-b529-4260-9408-5597eeeea034", // Identity pool ID
    Regions.AP_SOUTHEAST_2 // Region
);
```

### Permissions

The permissions system is set up as following:
 - **Users** who aren't in any group, can only access the public endpoints
   (read `/questions`, `/body-parts` etc.) and read and write their own data
   (`/users/me/*`).
 - **Doctors** can access the "unrestricted" endpoints on top of that (`/reports`,
   `/requests`, `/lesions` and so on).
 - **Administrators** can access the write endpoints of `/body-parts` and `/questions`.

## Example interaction flows

### User interacting with the service

1. User authenticates with AWS Cognito, sets up their acount over there etc. Ends
   up with an access token.
2. They load up body parts and questions from /body-parts and /questions respectively.
3. User loads up all their requests - GET /users/me/requests
4. User creates a new request and marks it as a draft - POST /users/me/requests.
5. User loads up all their lesions - GET /users/me/lesions
4. User creates a new lesion, selects the location on the body part - POST /users/me/lesions.
5. User loads up all their lesions again - GET /users/me/lesions
6. User selects a lesion, creates a new report for it - POST /users/me/lesions/:id/reports
7. User updates the request to submit it - PUT /users/me/requests/:id
8. *User waits for the doctor.*

### Doctor interacting with the service

1. Doctor authenticates using AWS Cognito.
2. They access all the requests - GET /requests?status=submitted&order_by=updated_date DESC
3. They browse all the information about the request - GET /reports, GET /lesions etc.
4. They write up an answer and update the request - PUT /requests/:id
