package db

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/temporalio/reference-app-orders-go/app/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "modernc.org/sqlite" // SQLite driver
)

// OrdersCollection is the name of the MongoDB collection to use for Orders.
const OrdersCollection = "orders"

// OrderStatus is a struct that represents the status of an Order
type OrderStatus struct {
	ID         string `db:"id" bson:"id"`
	CustomerID string `db:"customer_id" bson:"customer_id"`
	Status     string `db:"status" bson:"status"`

	ReceivedAt time.Time `db:"received_at" bson:"received_at"`
}

// ShipmentStatus is a struct that represents the status of a Shipment
type ShipmentStatus struct {
	ID     string `db:"id" bson:"id"`
	Status string `db:"status" bson:"status"`
}

// ShipmentCollection is the name of the MongoDB collection to use for Shipment data.
const ShipmentCollection = "shipments"

// DB is an interface that defines the methods that a database driver must implement
type DB interface {
	Connect(ctx context.Context) error
	Setup() error
	Close() error
	InsertOrder(context.Context, *OrderStatus) error
	UpdateOrderStatus(context.Context, string, string) error
	GetOrders(context.Context, *[]OrderStatus) error
	CountCompletedOrdersInRange(context.Context, time.Time, time.Time) (int, error)
	UpdateShipmentStatus(context.Context, string, string) error
	GetShipments(context.Context, *[]ShipmentStatus) error
	GetPendingShipments(context.Context, *[]ShipmentStatus) error
}

// CreateDB creates a new DB instance based on the configuration
func CreateDB(config config.AppConfig) DB {
	if config.MongoURL != "" {
		return &MongoDB{uri: config.MongoURL}
	}

	return &SQLiteDB{path: "./api-store.db"}
}

// MongoDB is a struct that implements the DB interface for MongoDB
type MongoDB struct {
	uri    string
	client *mongo.Client
	db     *mongo.Database
}

// Connect connects to a MongoDB instance
func (m *MongoDB) Connect(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.uri))
	if err != nil {
		return err
	}
	m.client = client
	m.db = client.Database("orders")
	return nil
}

// Setup sets up the MongoDB instance
func (m *MongoDB) Setup() error {
	orders := m.db.Collection(OrdersCollection)
	_, err := orders.Indexes().CreateOne(context.TODO(), mongodb.IndexModel{
		Keys: map[string]interface{}{"received_at": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create orders received_at index: %w", err)
	}

	_, err = orders.Indexes().CreateOne(context.TODO(), mongodb.IndexModel{
		Keys: map[string]interface{}{"completed_at": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create orders completed_at index: %w", err)
	}

	shipments := m.db.Collection(ShipmentCollection)
	_, err = shipments.Indexes().CreateOne(context.TODO(), mongodb.IndexModel{
		Keys: map[string]interface{}{"booked_at": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create shipment booked_at index: %w", err)
	}

	_, err = shipments.Indexes().CreateOne(context.TODO(), mongodb.IndexModel{
		Keys: map[string]interface{}{"status": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create shipment status index: %w", err)
	}

	return nil
}

// InsertOrder inserts an Order into the MongoDB instance
func (m *MongoDB) InsertOrder(ctx context.Context, order *OrderStatus) error {
	_, err := m.db.Collection(OrdersCollection).InsertOne(ctx, order)
	return err
}

// UpdateOrderStatus updates an Order in the MongoDB instance
func (m *MongoDB) UpdateOrderStatus(ctx context.Context, id string, status string) error {
	var fields bson.M
	if status == "completed" {
		fields = bson.M{"status": status, "completed_at": time.Now().UTC()}
	} else {
		fields = bson.M{"status": status}
	}
	_, err := m.db.Collection(OrdersCollection).UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": fields})
	return err
}

// GetOrders returns a list of Orders from the MongoDB instance
func (m *MongoDB) GetOrders(ctx context.Context, result *[]OrderStatus) error {
	res, err := m.db.Collection(OrdersCollection).Find(ctx, bson.M{}, &options.FindOptions{
		Sort: bson.M{"received_at": 1},
	})
	if err != nil {
		return err
	}

	return res.All(ctx, result)
}

// CountCompletedOrdersInRange counts completed Orders in a given time range
func (m *MongoDB) CountCompletedOrdersInRange(ctx context.Context, start time.Time, end time.Time) (int, error) {
	count, err := m.db.Collection(OrdersCollection).CountDocuments(ctx, bson.M{
		"completed_at": bson.M{
			"$gte": start,
			"$lte": end,
		},
	})
	return int(count), err
}

// UpdateShipmentStatus updates a Shipment in the MongoDB instance
func (m *MongoDB) UpdateShipmentStatus(ctx context.Context, id string, status string) error {
	_, err := m.db.Collection(ShipmentCollection).UpdateOne(
		ctx,
		bson.M{"id": id},
		bson.M{
			"$set":         bson.M{"status": status},
			"$setOnInsert": bson.M{"booked_at": time.Now().UTC()},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

// GetShipments returns a list of Shipments from the MongoDB instance
func (m *MongoDB) GetShipments(ctx context.Context, result *[]ShipmentStatus) error {
	res, err := m.db.Collection(ShipmentCollection).Find(ctx, bson.M{}, &options.FindOptions{
		Sort: bson.M{"booked_at": 1},
	})
	if err != nil {
		return err
	}

	return res.All(ctx, result)
}

// GetPendingShipments returns a list of pending Shipments from the MongoDB instance
func (m *MongoDB) GetPendingShipments(ctx context.Context, result *[]ShipmentStatus) error {
	res, err := m.db.Collection(ShipmentCollection).Find(ctx, bson.M{
		"status": bson.M{"$ne": "delivered"},
	}, &options.FindOptions{})
	if err != nil {
		return err
	}

	return res.All(ctx, result)
}

// Close closes the connection to the MongoDB instance
func (m *MongoDB) Close() error {
	return m.client.Disconnect(context.Background())
}

// SQLiteDB is a struct that implements the DB interface for SQLite
type SQLiteDB struct {
	path string
	db   *sqlx.DB
}

//go:embed schema.sql
var sqliteSchema string

// Connect connects to a SQLite instance
func (s *SQLiteDB) Connect(_ context.Context) error {
	db, err := sqlx.Connect("sqlite", s.path)
	if err != nil {
		return err
	}
	s.db = db
	db.SetMaxOpenConns(1) // SQLite does not support concurrent writes
	return nil
}

// Setup sets up the SQLite instance
func (s *SQLiteDB) Setup() error {
	_, err := s.db.Exec(sqliteSchema)
	return err
}

// Close closes the connection to the SQLite instance
func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

// InsertOrder inserts an Order into the SQLite instance
func (s *SQLiteDB) InsertOrder(ctx context.Context, order *OrderStatus) error {
	_, err := s.db.NamedExecContext(ctx, "INSERT OR IGNORE INTO orders (id, customer_id, received_at, status) VALUES (:id, :customer_id, :received_at, :status)", order)
	return err
}

// UpdateOrderStatus updates an Order in the SQLite instance
func (s *SQLiteDB) UpdateOrderStatus(ctx context.Context, id string, status string) error {
	var err error
	if status == "completed" {
		_, err = s.db.ExecContext(ctx, "UPDATE orders SET status = ?, completed_at = ? WHERE id = ?", status, time.Now().UTC(), id)
	} else {
		_, err = s.db.ExecContext(ctx, "UPDATE orders SET status = ? WHERE id = ?", status, id)
	}
	return err
}

// GetOrders returns a list of Orders from the SQLite instance
func (s *SQLiteDB) GetOrders(ctx context.Context, result *[]OrderStatus) error {
	return s.db.SelectContext(ctx, result, "SELECT id, status, received_at FROM orders ORDER BY received_at DESC")
}

// CountCompletedOrdersInRange counts completed Orders in a given time range
func (s *SQLiteDB) CountCompletedOrdersInRange(ctx context.Context, start time.Time, end time.Time) (int, error) {
	var count int
	err := s.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM orders WHERE completed_at >= ? AND completed_at <= ?", start, end)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateShipmentStatus updates a Shipment in the SQLite instance
func (s *SQLiteDB) UpdateShipmentStatus(ctx context.Context, id string, status string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO shipments (id, booked_at, status) VALUES (?, ?, ?) ON CONFLICT(id) DO UPDATE SET status = ?", id, time.Now().UTC(), status, status)
	return err
}

// GetShipments returns a list of Shipments from the SQLite instance
func (s *SQLiteDB) GetShipments(ctx context.Context, result *[]ShipmentStatus) error {
	return s.db.SelectContext(ctx, result, "SELECT id, status FROM shipments ORDER BY booked_at DESC")
}

// GetPendingShipments returns a list of pending Shipments from the SQLite instance
func (s *SQLiteDB) GetPendingShipments(ctx context.Context, result *[]ShipmentStatus) error {
	return s.db.SelectContext(ctx, result, "SELECT id, status FROM shipments WHERE status != 'delivered'")
}
