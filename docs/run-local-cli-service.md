# Run Everything Locally

Follow these instructions to run the OMS locally, using the 
Temporal Service provided by the `temporal server start-dev` 
command. This is an expanded version of the instructions found 
in the _Quickstart_ section of the top-level README file.

### Required Software
You will need [Go](https://go.dev/) to run the core OMS application,
the [Temporal CLI](https://docs.temporal.io/cli#install) to run the
Temporal Service locally, plus [Node.js](https://nodejs.org/) and
the [pnpm](https://pnpm.io/) package manager to run the OMS web
application.


### Start the Temporal Service

Run the following command in your terminal (from the root
directory of your repo clone) to start the Temporal Service:

```command
temporal server start-dev \
    --ui-port 8080 \
    --db-filename temporal-persistence.db
```

The Temporal Service manages application state by assigning tasks
related to each Workflow Execution and tracking the completion of 
those tasks. The detailed history it maintains for each execution 
enables the application to recover from a crash by reconstructing 
its pre-crash state and resuming the execution.

The `temporal` CLI provides a convenient way of running the Temporal 
Service locally for development purposes. By default, it provides a 
Web UI on port 8233 and persists data to an ephemeral in-memory 
database. The options in this command change the Web UI port to 8080 
and instructs the Temporal Service to persist its data to a file so 
that it will be available in subsequent sessions. This file will be 
created if it does not exist.

You can verify that this is running by using your browser to 
access the Temporal Web UI at <http://localhost:8080/>.


### Start the Workers

Run the following command in another terminal to start the Workers:

```command
go run ./cmd/oms worker
```
This command starts both Workflow and Activity Workers in a single 
process. The Workers run Workflow and Activity code, which carry out 
the various aspects of order processing in the OMS.


Although one Worker Process is sufficient for local development, you 
will want to run multiple Workers in production since this can improve 
both the scalability and availability of an application. You can 
repeat this step to launch as many additional Workers as you like. 
Temporal's SDK automatically distributes processing load among all 
running Workers.


### Start the API Servers

Run the following command in another terminal to start the API Servers:

```command
go run ./cmd/oms api
```

The API Servers provide REST APIs that the web application uses to 
interact with the OMS. This design decouples the web application from 
the Temporal Service and the order management system's back-end 
processing, which increases the flexibility and security of the entire 
system.


### Run the Web Application
You will need to clone the code for the web application, which is 
maintained separately in the [reference-app-orders-web](https://github.com/temporalio/reference-app-orders-web) repository:

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

### Next Steps

Setup is now complete. You have started the Temporal Service, 
the OMS Workers, the OMS API Servers, and the web application 
you'll use to interact with the OMS.

Continue by following the instructions for [processing a basic 
order](process-basic-order.md) or [processing a more complex 
order](process-complex-order.md). 

