# orders-reference-app
Order processing reference application

## Finding your way around the repository

* `app/` Application code
* `cmd/` Command line tools for the application
* `deployments/` Tools to deploy the application
* `docs/` Documentation
* `web/` Web interface and assets


### To run all Worker and API services

`go run ./cmd/dev-server`

### To run web

`cd web && pnpm i && pnpm dev`
