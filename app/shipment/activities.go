package shipment

import (
	"context"
	"fmt"
	"net/smtp"
)

type Activities struct {
	SMTPStub bool
	SMTPHost string
	SMTPPort int
}

var a Activities

// RegisterShipmentInput is the input for the RegisterShipment operation.
// All fields are required.
type RegisterShipmentInput struct {
	OrderID string
	Items   []Item
}

// RegisterShipmentResult is the result for the RegisterShipment operation.
// CourierReference is recorded where available, to allow tracking enquiries.
type RegisterShipmentResult struct {
	CourierReference string
}

// RegisterShipment registers a shipment with a courier.
func (a *Activities) RegisterShipment(ctx context.Context, input RegisterShipmentInput) (RegisterShipmentResult, error) {
	return RegisterShipmentResult{}, nil
}

const from = "orders@reference-app.example"
const to = "customer@reference-app.example"

// ShipmentCreatedNotificationInput is the input for a ShipmentCreated notification.
type ShipmentCreatedNotificationInput struct {
	OrderID string
}

// ShipmentCreatedNotification sends a ShipmentCreated notification to a user.
func (a *Activities) ShipmentCreatedNotification(ctx context.Context, input ShipmentCreatedNotificationInput) error {
	err := a.sendMail(from, to,
		fmt.Sprintf("Shipment for order: %s", input.OrderID),
		"Your order has been processed and shipping has been arranged with the courier. We'll be in touch once its dispatched.",
	)

	return err
}

// ShipmentDispatchedNotificationInput is the input for a ShipmentDispatched notification.
type ShipmentDispatchedNotificationInput struct {
	OrderID string
}

// ShipmentDispatchedNotification sends a ShipmentDispatched notification to a user.
func (a *Activities) ShipmentDispatchedNotification(ctx context.Context, input ShipmentDispatchedNotificationInput) error {
	err := a.sendMail(from, to,
		fmt.Sprintf("Shipment dispatched for order: %s", input.OrderID),
		"Your order has been dispatched.",
	)

	return err
}

// ShipmentDeliveredNotificationInput is the input for a ShipmentDelivered notification.
type ShipmentDeliveredNotificationInput struct {
	OrderID string
}

// ShipmentDeliveredNotification sends a ShipmentDelivered notification to a user.
func (a *Activities) ShipmentDeliveredNotification(ctx context.Context, input ShipmentDeliveredNotificationInput) error {
	err := a.sendMail(from, to,
		fmt.Sprintf("Shipment delivered for order: %s", input.OrderID),
		"Your order has been delivered.",
	)

	return err
}

func (a *Activities) sendMail(from string, to string, subject string, body string) error {
	if a.SMTPStub {
		return nil
	}

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", a.SMTPHost, a.SMTPPort),
		nil,
		from,
		[]string{to},
		[]byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)),
	)
}
