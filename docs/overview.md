# System Overview

```mermaid
sequenceDiagram
    participant Customer
    participant Order
    participant Shipment
    participant Courier
    
    Customer->>Order: place order
    Order->>Shipment: create shipment
    Shipment->>Courier: register shipment
    Courier->>Shipment: shipment registered
    Shipment->>Customer: shipment created
    Courier->>Shipment: shipment dispatched
    Shipment->>Customer: shipment dispatched
    Courier->>Shipment: shipment delivered
    Shipment->>Customer: shipment delivered
    Order->>Customer: order complete
```