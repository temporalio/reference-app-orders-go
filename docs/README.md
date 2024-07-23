# Order Management System (Go) Documentation

![OMS logo](images/oms-logo.png)

The documentation for the Order Management System reference application 
is organized into multiple sections:

## Understanding the OMS
* [Overview](overview.md): 
      Provides a brief high-level overview of the OMS
* [Product Requirements](product-requirements.md):
      Describes what the OMS does
* [Technical Description](technical-description.md):
      Describes the design and implementation of the OMS

## Running the OMS
* [Using a Local Temporal Service](run-local-cli-service.md): 
      Provides step-by-step instructions for running the 
      application locally, with a Temporal Service provided
	  by the `temporal` CLI
* [Using Temporal Cloud](run-temporal-cloud.md): 
      Provides step-by-step instructions for running the 
      application locally, with a Temporal Service provided 
	  by Temporal Cloud
* [Using Encryption](run-codec-server.md): 
      Provides step-by-step instructions for running the 
	  application with a Custom Data Converter that will 
	  encrypt payload data and a Codec Server that enables 
	  you to view decrypted data in the Temporal Web UI 
	  and output of the `temporal` command

## Processing Orders
* [Processing a Basic (Single-Item) Order](process-basic-order.md): 
      Describes how to use the web application to process a basic 
	  order, which consists of a single item in a single shipment
* [Processing a Complex Order](process-complex-order.md): 
      Describes how to use the web application to process a more 
	  complex order, which involves multiple shipments and an 
	  out-of-stock item that requires customer interaction for
	  processing to continue

## Deploying the OMS to Production
* [Deploying to Kubernetes](deploy-on-k8s.md) 
