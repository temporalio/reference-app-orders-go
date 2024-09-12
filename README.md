# Temporal Reference Application: Order Management System (Go)

![OMS logo](docs/images/oms-logo.png)

The Order Management System (OMS) is a reference application that 
demonstrates one way to approach the design and implementation of 
an order processing system based on Temporal Workflows. You can run 
this application locally (directly on a laptop) or in a Kubernetes 
cluster. In addition, the required Temporal Service can be run locally, 
or be provided by a remote self-hosted deployment, or be provided by 
Temporal Cloud. 

## Quickstart
We recommend that you begin by reading the [documentation](docs/README.md), 
which will explain the features of the application and aspects 
of its design. It also provides instructions for deploying and 
running the application in various environments.

If you'd like to jump right in and run the OMS locally, clone this 
repository to your machine and follow the steps below. Unless otherwise 
noted, you should execute the commands from the root directory of your 
clone.

### Required Software
You will need [Go](https://go.dev/) to run the core OMS application, 
the [Temporal CLI](https://docs.temporal.io/cli#install) to run the 
Temporal Service locally, plus [Node.js](https://nodejs.org/) and 
the [pnpm](https://pnpm.io/) package manager to run the OMS web 
application. 


### Start the Temporal Service
Run the following command in your terminal:

```command
temporal server start-dev --ui-port 8080 --db-filename temporal-persistence.db
```

The Temporal Service manages application state by assigning tasks
related to each Workflow Execution and tracking the completion of 
those tasks. The detailed history it maintains for each execution 
enables the application to recover from a crash by reconstructing 
its pre-crash state and resuming the execution.

### Start the Workers
Run the following command in another terminal:

```command
go run ./cmd/oms worker
```

This command starts both Workflow and Activity Workers in a single
process. The Workers run Workflow and Activity functions, which 
carry out the various aspects of order processing.

### Start the API Servers
Run the following command in another terminal:
```command
go run ./cmd/oms api
```

The API Servers provide REST APIs that the web application uses to 
interact with the OMS. 


### Run the Web Application
You will need to clone the code for the web application, which is 
maintained separately in the 
[reference-app-orders-web](https://github.com/temporalio/reference-app-orders-web) 
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
<http://localhost:5173/> and the Temporal Web UI at 
<http://localhost:8080/>. In the OMS web application, select 
the **User** role, and then submit an order (we recommend 
choosing order #1 to start). Next, return to the main page
of the web application, select the **Courier** role, locate
the shipments corresponding to your order, and then click 
the **Dispatch** and **Deliver** buttons to complete the 
process. As you proceed with each of these steps, be sure 
to refresh the Temporal Web UI so that you can see the 
Workflows created and updated as a result. 


## Find Your Way Around
This repository provides four subdirectories of interest:

| Directory                                             | Description                                                       |
| ----------------------------------------------------- | ----------------------------------------------------------------- |
| <code><a href="app/">app/</a></code>                  | Application code                                                  |
| <code><a href="cmd/">cmd/</a></code>                  | Command-line tools provided by the application                    |
| <code><a href="deployments/">deployments/</a></code>  | Tools and configuration files used to deploy the application      |
| <code><a href="docs/">docs/</a></code>                | Documentation                                                     |

See the [documentation](docs/README.md) for more information.
