# System Overview

```mermaid
sequenceDiagram
    participant Customer
    participant Order
    participant Inventory
    participant Billing
    participant Shipment
    participant Courier
    
    Customer->>Order: place order
    Order->>Inventory: fulfill order
    Order->>Billing: create invoice
    Order->>Billing: charge customer
    Order->>Shipment: create shipment
    Shipment->>Courier: book shipment
    Courier->>Shipment: shipment booked
    Shipment->>Customer: shipment booked
    Courier->>Shipment: shipment dispatched
    Shipment->>Customer: shipment dispatched
    Courier->>Shipment: shipment delivered
    Shipment->>Customer: shipment delivered
    Order->>Customer: order complete
```
