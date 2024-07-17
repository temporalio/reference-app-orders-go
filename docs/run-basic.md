# Running the OMS Locally

Follow these instructions to run the OMS locally, using the 
Temporal Service provided by the `temporal server start-dev` 
command. This is an expanded version of the instructions found 
in the _Quickstart_ section of the top-level README file.


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


### Start the Worker

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

<!--
	TODO: expand this section to cover more advanced cases
	      and then move it to a separate document (with a 
		  demo video) that can be referenced by the other
		  instructions for running it.
-->

You will then be able to access the OMS web application at 
<http://localhost:5173/> and the Temporal Web UI at 
<http://localhost:8080/>. In the OMS web application, select 
the **User** role, and then submit an order (we recommend 
choosing order #1 to start). Next, return to the main page of  
the web application, select the **Courier** role, locate
the shipments corresponding to your order, and then click 
the **Dispatch** and **Deliver** buttons to complete the 
process. As you proceed with each of these steps, be sure 
to refresh the Temporal Web UI so that you can see the 
Workflows created and updated as a result. 

