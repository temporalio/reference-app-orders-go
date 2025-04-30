# Running the OMS Locally, Using Temporal Cloud

The [Using a Local Temporal Service](run-local-cli-service.md)
instructions described how to run the OMS locally, with 
its Workers and API servers relying on the Temporal Service 
provided by the `temporal server start-dev` command.

The instructions on this page also describe how to run 
the OMS locally, but with its Workers and API servers 
using [Temporal Cloud](https://temporal.io/cloud) instead 
of a local Temporal Service. Successfully completing these 
steps requires a Temporal Cloud account. 


### Set the Environment Variables

By default, the OMS assumes the use of the `default` Namespace 
and a Temporal Service that listens on `localhost` port 7233 
without TLS. However, its design enables you to
[customize those settings](https://github.com/temporalio/reference-app-orders-go/blob/3fa995740d2f9ad31890c0ca093bc40524250a19/app/server/server.go#L26-L69) 
by setting environment variables. Therefore, moving from a local
Temporal Service to one provided by Temporal Cloud requires no 
change to application code.

#### Authentication Options

Temporal Cloud supports two authentication methods: mTLS (mutual TLS) and API Keys. You need to choose one of these methods when connecting to Temporal Cloud.

##### Option 1: Using mTLS Authentication

If using mTLS authentication, you must define the following environment variables. We recommend defining them in a reusable script, since you'll need to set them in multiple terminals:

1. **`TEMPORAL_NAMESPACE`**: Set this to the name of your Namespace 
    in the Temporal Cloud Account
    * Example: `oms-demo.d6rd8`
2. **`TEMPORAL_ADDRESS`**: Set this to the gRPC Endpoint for your 
    Namespace that uses the `tmprl.cloud` domain
    * Example: `oms-demo.d6rd8.tmprl.cloud:7233`
3. **`TEMPORAL_TLS_CERT`**: Set this to the path of a TLS certificate 
    file associated with your Namespace
    * Example: `/Users/tomwheeler/private/tls/oms-demo.pem`
4. **`TEMPORAL_TLS_KEY`**: Set this to the path of the private key
    for your TLS certificate
    * Example: `/Users/tomwheeler/private/tls/oms-demo.key`

##### Option 2: Using API Key Authentication

If using API Key authentication, you must define the following environment variables.  We recommend defining them in a reusable script, since you'll need to set them in multiple terminals

1. **`TEMPORAL_NAMESPACE`**: Set this to the name of your Namespace 
    in the Temporal Cloud Account
    * Example: `oms-demo.d6rd8`
2. **`TEMPORAL_ADDRESS`**: Set this to the API Key gRPC Endpoint for your 
    Namespace that uses the `temporal.io` domain (note the domain difference from mTLS)
    * Example: `oms-demo.d6rd8.temporal.io:7233`
3. **`TEMPORAL_API_KEY`**: Set this to your Temporal Cloud API Key

You can find both the endpoint and API key in the Temporal Cloud UI under your namespace settings.

Because you'll use Temporal Cloud, you won't need to run the 
Temporal Service locally.

### Start the Workers

Make sure the environment variables described earlier are 
properly set in your terminal, and then execute this command 
to start the Workers:

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

Ensure that the environment variables described earlier are 
properly set in your terminal, and then run the command to 
start the API Servers:

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

Since the web application does not interact with the Temporal 
Service, it is unnecessary to set the environment variables in
the terminal where you'll run the web application.


Run the following commands to start it:

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

