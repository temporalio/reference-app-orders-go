# Process a Basic (Single-Item) Order

Follow these instructions to process a basic order, which will consist
of a single item currently in stock. These instructions assume that 
you have already started the OMS and the corresponding web application.

We also recommend opening the Temporal Web UI to view the details of 
each Workflow Execution started by the OMS to process the order. The 
URL used to access the Temporal Web UI will vary according to where 
the Temporal Service is running and how it was deployed.

## Submit an Order in the Web Application

1. Access the OMS web application, which should be available 
   at <http://localhost:5173/>.
2. Click the link for the **Customer** role
3. Select **Order #1** (this contains a single item, which is always 
   available in inventory at Warehouse A)
4. Click the **SUBMIT** button

## Process the Shipment in the Web Application

1. Open a new browser tab or window and use it to access 
   the front page of the web application. You'll use this 
   to manage the shipment corresponding to the order, as 
   described in the next few steps.
2. Click the link for the **Courier** role
3. The page now shows all available shipments. Click the 
   link for the shipment created for your order in the 
   previous section.
4. Click the **DISPATCH** button at the bottom of the page, 
   which indicates that the courier has picked up the shipment 
   from the warehouse.
5. Click the **DELIVER** button at the bottom of the page, 
   which indicates that the courier has delivered the shipment 
   to the customer.

Since this order contained a single item, there is only one 
shipment. Now that the shipment has been delivered to the 
customer, processing for both the order and the shipment is 
complete. We recommend you continue exploring the OMS by 
following the instructions for [processing a more complex 
order](process-complex-order.md).
