# Running the OMS Locally (Basic)

Follow these instructions to run the OMS locally, using the 
Temporal Service provided that the `temporal` command. This 
This is an expanded version of the instructions found in the 
_Quickstart_ section of the top-level README file.


### Start the Temporal Service

Run the following command in your terminal to start the Temporal 
Service:

```command
temporal server start-dev \
    --ui-port 8080 \
    --db-filename temporal-persistence.db
```

The `temporal` CLI provides a convenient way of running a 
Temporal Service locally for development purposes. By default,
it provides a Web UI on port 8233 and persists data to an 
ephemeral in-memory database. The options in this command 
change the Web UI port to 8080 and instructs the Temporal 
Service to persist its data to a file so that it will be 
available in subsequent sessions. This file will be created 
if it does not exist.


### Start the Worker

Run the following command in another terminal to start the Workers:

```command
go run ./cmd/oms worker
```

The Workers run the Workflow and Activity code used to implement the
OMS. Although one Worker is sufficient for local development, Temporal
recommends running multiple Workers in production since this can improve
both the scalability and availability of an application. You can repeat
this step to launch as many additional Workers as you like.


### Start the API Servers

Run the following command in another terminal to start the API Servers:

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

