# Order Management System (Go) - Product Requirements

![OMS logo](images/oms-logo.png)


## Vision
The Order Management System (OMS) is a mission-critical system for 
processing product orders.

Reliability is essential, since failure to process individual orders
correctly or system outages that affect the processing of all orders
have a substantial negative impact on both revenue and customer
satisfaction. The OMS must be flexible, allowing the development team 
to add new features and modify existing ones as business needs dictate.
Finally, it must enable support staff to quickly locate an order and
understand its current status, as well as audit the steps that led to
that status.

Note that while the OMS is responsible for _processing_ orders, 
it is not involved with _composing_ those orders. Presumably this is 
done by something upstream from the OMS, such as a shopping cart that 
is part of a typical e-commerce application.



## Personas: Who Uses It?
There are three primary roles who will interact with the OMS:

* **Customer**: Someone who orders products from the store 
* **Courier**: Someone who delivers products to the customer 
* **Manager**: Someone who performs administrative functions

Additionally, developers and support staff use the OMS to locate 
details about both open and recently closed orders, enabling them 
to respond to customer inquiries.


## Features: What Does It Do?


### Customer Interaction: Typical Flow 
When a customer submits an order, the OMS accepts the order details as 
input. It then requests the items from inventory, potentially splitting 
them into multiple shipments based on whether they are all available 
from the same warehouse.

The OMS contacts the billing system to calculate the total cost of each
shipment, including tax and shipping, and then generates an invoice and
charges the customer. Since a damaged or lost package will only affect a
single shipment, the customer is billed on a per-shipment basis instead
of a per-order basis.

After the customer is successfully billed for the shipment, the OMS
contacts a courier service to request that they deliver the package.
After this is booked, the OMS waits for a driver to be dispatched to the
warehouse, pick up the shipment, and deliver it to the customer. Once
all shipments have been delivered to the customer, the order is closed.

At any time after placing an order, the customer may view its status
from the Orders page of the web application. The detail page for an
order lists each of its shipments, including its invoice and current
status.


#### Customer Interaction: Item(s) Unavailable
In some cases, at least one item in the order is not available in any
warehouse. When this occurs, the OMS amends the order to exclude the
unavailable item(s) and waits for the customer to accept or cancel the
order. If the customer accepts the amended order, the processing
continues as described in the preceding section. If the customer doesn't
respond within a set period, the order times out and all shipments are
canceled.


#### Customer Interaction: Charge Rejected 
In order to prevent fraud, the store manager will have the option to set 
a cumulative spending limit that is uniformly applied to each customer.
For example, if the limit is set to $1000 per customer, then an order
for $1200 from a given customer would fail (and thus, their order would
fail). However, the customer could successfully place an order for $700,
but a subsequent order for $500 from that customer would fail because
the sum of all their orders exceeds the limit. This spending limit is
not set by default and can be increased or removed by a manager at any
time.


### Courier Interaction
When Couriers accesses the web application, they are presented with a 
list of shipments and the current status of each. When the courier 
arrives at the warehouse to pick up a shipment, they will click the 
detail page for that shipment, and then click the "Dispatch" button on 
that page. After they successfully deliver it to the customer, they will 
click the "Deliver" button. These interactions update the current status 
of each shipment and are immediately visible to customers and support 
staff.

### Manager Interaction
As previously described, the store manager has the ability to combat 
fraud by setting a global limit on the total charges (expressed in 
cents) that each customer is allowed. This is set to 0 by default, 
meaning that there is no limit. The manager can increase, decrease, 
or reset this limit at any time.
