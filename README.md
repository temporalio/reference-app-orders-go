# Temporal Reference Application: Order Management System (Go)

![OMS logo](docs/images/oms-logo.png)

The Order Management System (OMS) is a reference application 
that demonstrates how to design and implement an order processing 
system with Temporal. You can run this application locally 
(directly on a laptop) or in a Kubernetes cluster. In either case, 
the Temporal Service can be local, a remote self-hosted deployment, 
or Temporal Cloud. 

NOTE: This application is under active development and we're
working to expand the documentation and finish up a few features 
before we encourage widespread use. Please check back regularly.

## Quickstart
We recommend that you begin by reading the [documentation](docs/README.md), 
which will explain the features of the application as well as aspects 
of its design. It also provides instructions for deploying and 
running the application in various environments, including in 
Kubernetes and with Temporal Cloud.

If you'd prefer to jump right in and run it locally, follow these steps. 
All commands should be executed from the root directory of this project, 
unless otherwise noted. 

### Start the Temporal Service
Run the following command in your terminal:

```command
temporal server start-dev --ui-port 8080 --db-filename temporal-persistence.db
```

This command uses the `--db-filename` option so that the development 
server will persist its data to a file instead of memory, thus making 
it available during later sessions. The file will be created if it
does not already exist.

### Start the Worker(s)
Run the following command in another terminal:

```command
go run ./cmd/oms worker
```

Although one Worker is sufficient for local development, we recommend 
running multiple Workers in production since this can improve both the 
scalability and availability of an application. You can repeat this 
step to launch as many additional Workers as you like.

### Start the API Servers
Run the following command in another terminal:
```command
go run ./cmd/oms api
```

The API Servers provide REST APIs that the web application uses to 
interact with the OMS. 


### Run the Web Application
You will need to clone the code for the web application, which is 
maintained in the [reference-app-orders-web](https://github.com/temporalio/reference-app-orders-web) 
repository:

```command
cd ..
git clone https://github.com/temporalio/reference-app-orders-web.git
```

You will then need to run the following commands to start it:

```command
cd reference-app-orders-web
pnpm install
pnpm dev
```

You will then be able to access the OMS web application at 
[http://localhost:5173/] and the Temporal Web UI at 
[http://localhost:8080/]. In the OMS web application, select 
the **User** role, submit an order (we recommend order #1 to 
start), and then return to select the **Courier** role and 
process. As you proceed with each of these steps, be sure to 
refresh the Temporal Web UI so that you can see the Workflows 
created and updated as a result. 


## Finding Your Way Around
This repository provides four subdirectories of interest:

| Directory                                             | Description                                                       |
| ----------------------------------------------------- | ----------------------------------------------------------------- |
| <code><a href="app/">app/</a></code>                  | Application code                                                  |
| <code><a href="cmd/">cmd/</a></code>                  | Command-line tools provided by the application                    |
| <code><a href="deployments/">deployments/</a></code>  | Tools and configuration files used to deploy the application      |
| <code><a href="docs/">docs/</a></code>                | Documentation                                                     |

See the [documentation](docs/README.md) for more information.
