# Deploying on Kubernetes

In order to deploy the application to Kubernetes we have provided some manifests and helper scripts in the `./deployments` directory.

The scripts will help you:
- Create a Kubernetes cluster on EKS, if you don't already have one you'd like to use
- Deploy a basic Temporal cluster to the Kubernetes cluster, if you don't already have a Temporal install you'd like to use
- Deploy the OMS application into the Kubernetes cluster
- Access the OMS Web UI and Temporal Web UI from your local machine

## Creating a Kubernetes cluster

We have provided a script to create a new EKS cluster on AWS, but using EKS is not a requirement. If you already have a Kubernetes cluster you would like to use, from AWS or any other provider, please feel free to skip this step. Any Kubernetes cluster will work.

In order to use the script please ensure you have valid credentials setup to use the AWS APIs. Setting up authentication is beyond the scope of this document, but you can find some details on what IAM policies you may need in the [`eksctl` documentation](https://eksctl.io/usage/minimum-iam-policies/).

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

## Installing Temporal into Kubernetes

You can use the Temporal Helm Charts to install a basic Temporal cluster into Kubernetes.

To install Temporal:

```sh
./deployments/install-temporal.sh
```

This may take a while, as the Temporal cluster relies on Cassandra and Elasticsearch, which can both take a while to boot up.

Once the Temporal cluster is up, the script creates the Temporal namespace `default` for you. This is the Temporal namespace the OMS application expects to use, unless configured otherwise.

## Installing the application into Kubernetes

Once you have the Kubernetes cluster and Temporal is installed, it's time to install the application. We have created some manifests that will take care of installing the application for you. As is best practice, the OMS application will live in its own Kubernetes namespace `oms`. To ensure this namespace is present before we try and install anything into it, create the namespace using `kubectl`:

```sh
kubectl apply -f ./deployments/oms-namespace.yaml
```

Once that has completed, you can then install all of the application:

```sh
kubectl apply -f ./deployments/
```

To ensure the application is running you can check on the status of pods:

```sh
kubectl get pods -n oms
```

You should see something like the following:

```
NAME                              READY   STATUS    RESTARTS   AGE
billing-api-0                     1/1     Running   0          24h
billing-worker-5498c6ffd4-klfwk   1/1     Running   0          24h
codec-server-847b59b84f-9ct54     1/1     Running   0          24h
main-api-0                        1/1     Running   0          24h
main-worker-777c59ccdb-jx46w      1/1     Running   0          24h
web-557c8dc97-p5fkq               1/1     Running   0          24h
```

Once the READY column shows "1/1" for all the pods, the application is up.

Please note that while in most regards the deployment is a standard Kubernetes setup, there is one thing which is a little unusual. Our API services make use of a cache to ensure that they can serve the listings for orders and shipments quickly. This cache lives in an SQLite file on disk, which makes the API services stateful. If an API pod needs to be restarted (due to upgrade, node failure etc) then it must be brought back up with the same disk as before to ensure it does not lose its cache. For this reason, the API services use Statefulset instead of Deployment. The cache is mounted via a Kubernetes volume to ensure that the pod will maintain it's disk between runs.

As this cache uses the disk, which is not easily shared amongst pods, this also means the API Statefulsets cannot be safely scaled. If there were to be more than one main-api pod, for example, each pod would see different order creations, and therefore have a different cache of orders. This would result in unpredictable behaviour.

These are all limitations of our deliberately naive caching mechanism, and would not occur should the APIs use a shared database such as MySQL or PostgreSQL.

## Using the application

So now you have the application installed in Kubernetes, how can you use it? The standard practise for accessing a service inside a Kubernetes cluster would be to create a Service with type LoadBalancer. On AWS, for example, this would cause AWS to automatically provision an ELB or ALB for you and connect it to the Kubernetes service, allowing external access. As we are only providing a demonstration system and don't wish to make things available over the public internet, we are not using LoadBalancer in our deployment. Instead we provide some scripts which use Kubernetes port-forwarding so that you can reach the Services in Kubernetes from your local machine.

To see the application's Web UI you can run the script:

```sh
./deployments/port-forward-web.sh
```

That command will continue to run until you terminate it (with Ctrl+C) for example. While the command is running it will forward any traffic from your local machine on port 3000 to the application's Web UI in the cluster. You should now be able to see the application at: http://localhost:3000/

We also provide two other port forwarding scripts. The first is to give you access to the web UI for the Temporal cluster running in your cluster:

```sh
./deployments/port-forward-temporal-ui.sh
```

While this command is running you will be able to see the Temporal Web UI at: http://localhost:8080. Try placing an order with the application's Web UI and then looking at the Temporal Web UI to see which workflows were created.

You may notice on the Temporal Web UI that the inputs and results for the Workflows and Activities are encrypted, due to the application's use of an encrypting data converter. In order to decrypt those on the Web UI so that you can see what the inptus and results were, you can use a codec server which we deployed alongside the app.

The last of our port-forwarding commands is to allow you access to the codec server:

```sh
./deployments/port-forward-codec-server.sh
```

Leave this command running and browse to the Temporal Web UI at http://localhost:8080. Now you can configure the Temporal Web UI to use the application's codec server to decrypt inputs and results by clicking on the sunglasses icon top right of the page. In the pop-up that displays, enter "http://localhost:8089" into the "Codec Server browser endpoint" field, and hit "Apply". You should now be able to navigate to one of the workflows and see the inputs and results in plaintext. Remember: the codec server will only work if the port forwarding command is running, otherwise your browser will not be able to access the codec server. In this case you would see encrypted content again, and the sunglasses icon on the Temporal UI will turn red to indicate a problem talking to the codec server. Just start the port forward command again, and the Temporal UI will then be able to reconnect next time you navigate to a workflow page (or refresh the one you are on).