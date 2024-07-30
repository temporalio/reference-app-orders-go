# Process a More Complex Order

Follow these instructions to process a more complex order, which 
consists of multiple items. One of these items is unavailable in
inventory, so the OMS will amend the order, and you'll be prompted 
to accept or cancel it. The OMS also sets a Timer in this case, so 
if you don't respond quickly enough, the order will time out and
its processing will end.

These instructions assume that you have already started the OMS 
and the corresponding web application. 

We also recommend opening the Temporal Web UI to view the details of 
each Workflow Execution started by the OMS to process the order. The 
URL used to access the Temporal Web UI will vary according to where 
the Temporal Service is running and how it was deployed.

## Submit an Order in the Web Application

1. Access the OMS web application, which should be available 
   at <http://localhost:5173/>.
2. Click the link for the **Customer** role
3. All but the first order is randomly generated. Starting with 
   Order #2, click through each order until you locate one with
   at least three items, one of which is the Adidas UltraBoost. 
   For demonstration purposes, this item is always out of stock.
4. Click the **SUBMIT** button

## Accept the Amended Order

1. The order status page should now show an "Action Required" 
   message, listing the Adidas UltraBoost as unavailable. 
   It lists two shipments below this, one containing an item
   available from Warehouse A and the other containing items
   available from Warehouse B). These show a "Pending" status,
   since the next processing step depends on what you do now.
2. Click the **AMEND** button to accept the amended order. The 
   order will time out if you fail to do this within 30 seconds.

These instructions describe accepting the amended order so that 
processing will continue. We recommend following these steps the 
first time, but afterwards, you might wish to experiment with 
explicitly canceling the amended order or not responding and 
allowing it time out.

## Process the Shipments in the Web Application

1. Open a new browser tab or window and use it to access 
   the front page of the web application. You'll use this 
   to manage the shipment corresponding to the order, as 
   described in the next few steps.
2. Click the link for the **Courier** role. You will then 
   see a page showing all available shipments. 
3. Click the link for the first shipment created for your 
   order in the previous section.
4. Click the **DISPATCH** button at the bottom of the page, 
   which indicates that the courier has picked up the shipment 
   from the warehouse.
5. Click the **DELIVER** button at the bottom of the page, 
   which indicates that the courier has delivered the shipment 
   to the customer.
6. Repeat the previous three steps for the second shipment
   in your order.

All shipments in this order have been delivered to the customer, 
so processing for both the order and these shipments is now 
complete.
