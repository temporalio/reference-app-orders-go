# Temporal Reference Application: Order Management System

![Screen Shot 2024-04-29 at 9 22 59 AM](https://github.com/temporalio/orders-reference-app-go/assets/7967403/b1ff7aa2-f3d6-4f47-9113-9dee1015634d)


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
