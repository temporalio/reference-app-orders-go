# Temporal Reference Application: Order Management System

## Finding your way around the repository

* `app/` Application code
* `cmd/` Command line tools for the application
* `deployments/` Tools to deploy the application
* `docs/` Documentation
* `web/` Web interface and assets


### To run all Worker and API services

`go run ./cmd/oms worker`
`go run ./cmd/oms api`

### To run web

`cd web && pnpm i && pnpm dev`
