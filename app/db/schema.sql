CREATE TABLE IF NOT EXISTS orders (
    id TEXT PRIMARY KEY,
    customer_id TEXT NOT NULL,
    received_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    status TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS orders_received_at ON orders(received_at DESC);
CREATE INDEX IF NOT EXISTS orders_completed_at ON orders(completed_at DESC);

CREATE TABLE IF NOT EXISTS shipments (
    id TEXT PRIMARY KEY,
    status TEXT NOT NULL,
    booked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS shipments_booked_at ON shipments (booked_at DESC);
CREATE INDEX IF NOT EXISTS shipments_pending ON shipments (status != 'delivered');
