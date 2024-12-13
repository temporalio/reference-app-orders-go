# Deploying on Kubernetes

In order to deploy the application to Kubernetes we have provided some manifests and helper scripts in the `./deployments` directory.

The scripts will help you:
- Create a Kubernetes cluster on EKS, if you don't already have one you'd like to use
- Deploy a basic Temporal cluster to the Kubernetes cluster, if you don't already have a Temporal install you'd like to use
- Deploy the OMS application into the Kubernetes cluster
- Access the OMS Web UI and Temporal Web UI from your local machine

## Creating a Kubernetes Cluster

We have provided a script to create a new EKS cluster on AWS, but using EKS is not a requirement. If you already have a Kubernetes cluster you would like to use, from AWS or any other provider, please feel free to skip this step. Any Kubernetes cluster will work.

In order to use the script please ensure you have valid credentials set up to use the AWS APIs. Setting up authentication is beyond the scope of this document, but you can find some details on what IAM policies you may need in the [`eksctl` documentation](https://eksctl.io/usage/minimum-iam-policies/).

To create the cluster (called "temporal-oms"):

```sh
./deployments/create-eks.sh
```

This will create a new EKS cluster on AWS and configure `kubectl` to talk to it.

You should be able to run:

```sh
kubectl get pods
```

And see the message `No resources found in default namespace.`

## Deploy Temporal to Kubernetes

You can use the Temporal Helm Charts to install a basic Temporal 
Cluster into Kubernetes.

To install Temporal:

```sh
./deployments/install-temporal.sh
```

This may take a while, as the Temporal Cluster relies on Cassandra
and Elasticsearch, which can both take a while to boot up.

Once the Temporal Cluster is up, the script creates the Temporal 
Namespace `default` for you. This is the Temporal Namespace the OMS 
application expects to use, unless configured otherwise.

## Deploy the OMS application to Kubernetes

Once you have Temporal installed and running in your Kubernetes 
cluster, it's time to install the OMS application. We have created 
some manifests that will take care of installing the application for 
you. As is best practice, the OMS application will live in its own 
Kubernetes namespace `oms`. To ensure this namespace is present before 
installing anything into it, use the `kubectl` command to create
the namespace:

```sh
kubectl apply -f ./deployments/k8s/oms-namespace.yaml
```

Once that has completed, you can then install all of the application:

```sh
kubectl apply -f ./deployments/k8s
```

You can check pod status to verify that the application is running:

```sh
kubectl get pods -n oms
```

You should see output similar to the following:

```
NAME                              READY   STATUS    RESTARTS   AGE
billing-api-5fdc8b9d8f-tm2tl      1/1     Running   0          24h
billing-worker-5498c6ffd4-klfwk   1/1     Running   0          24h
codec-server-847b59b84f-9ct54     1/1     Running   0          24h
main-api-777cd87b47-bsw6h         1/1     Running   0          24h
main-worker-777c59ccdb-jx46w      1/1     Running   0          24h
web-557c8dc97-p5fkq               1/1     Running   0          24h
mongo-0                           1/1     Running   0          24h
```

You'll know that the application is up when the READY column shows
"1/1" for every pod.

### Application cache

There is a noteworthy configuration difference when deploying the 
OMS to Kubernetes. As detailed in the 
[technical description](../technical-description.md#application-cache), 
the OMS application's API servers maintain a cache of all orders 
and shipments. By default, this is backed by a file-based SQLite 
database. While this choice of database inherently limits scalability 
to a single node, we felt that adding a dependency on an external 
database created an obstacle for developers who want to get the 
OMS running with minimal effort. We have since updated the OMS to 
add the option to use MongoDB for the application cache. Not only
does this represent a more typical production deployment, given
MongoDB's reputation for scalability, it also avoids the complexity 
of using Statefulsets for the API servers. For those reasons, we 
default to MongoDB for the API servers' cache when deploying the 
OMS to Kubernetes. 

## Using the Application

Now that the application is successfully installed in Kubernetes, how 
can you use it? The standard practice for accessing a service inside a 
Kubernetes cluster would be to create a Service with type LoadBalancer. 
On AWS, for example, this would cause AWS to automatically provision an 
ELB or ALB for you and connect it to the Kubernetes service, allowing 
external access. As we are only providing a demonstration system and 
don't wish to make things available over the public internet, we do not 
use a LoadBalancer in our deployment. Instead we provide some scripts 
which use Kubernetes port-forwarding so that you can reach the Services 
in Kubernetes from your local machine.

To see the application's Web UI you can run the script:

```sh
./deployments/port-forward-web.sh
```

That command will continue to run until you terminate it (for example, 
by pressing Ctrl+C in the terminal). While the command is running, it 
forwards any traffic from your local machine on port 3000 to the 
application's Web UI in the cluster. You should now be able to see the 
application at: http://localhost:3000/

We also provide two other port forwarding scripts. The first gives you 
access to the web UI for the Temporal Cluster running inside of 
Kubernetes:

```sh
./deployments/port-forward-temporal-ui.sh
```

While this command is running, you will be able to see the Temporal 
Web UI at: http://localhost:8080. Try placing an order with the 
application's Web UI and then looking at the Temporal Web UI to see 
which Workflows were created.

You may notice on the Temporal Web UI that the inputs and results for 
the Workflows and Activities are encrypted, due to the application's 
use of an encrypting data converter. In order to decrypt those on the 
Web UI, which will enable you to see the input and result values, you 
can use a Codec Server that is deployed alongside the app.

The last of our port-forwarding commands is to allow you access to the 
Codec Server:

```sh
./deployments/port-forward-codec-server.sh
```

Leave this command running and browse to the Temporal Web UI at 
http://localhost:8080. Now you can configure the Temporal Web UI to 
use the application's Codec Server to decrypt inputs and results by 
clicking on the sunglasses icon top right of the page. In the pop-up 
that displays, enter "http://localhost:8089" into the **Codec Server 
browser endpoint** field, and click the **Apply** button. You should 
now be able to navigate to one of the Workflow Executions and see the 
inputs and results in plaintext. 

Remember: the Codec Server will only work if the port forwarding 
command is running, otherwise your browser will not be able to access 
the Codec Server. In this case, you would see encrypted content again, 
and the sunglasses icon on the Temporal UI will turn red to indicate a 
communication problem with the Codec Server. Just start the port 
forward command again, and the Temporal UI will then be able to 
reconnect next time you navigate to a Workflow Execution detail page 
(or refresh the one you are on).
