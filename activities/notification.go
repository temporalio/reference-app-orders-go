package activities

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

const from = "orders@reference-app.example"
const to = "customer@reference-app.example"

// ShipmentCreatedNotificationInput is the input for a ShipmentCreated notification.
type ShipmentCreatedNotificationInput struct {
	OrderID ordersapi.OrderID
}

// ShipmentCreatedNotification sends a ShipmentCreated notification to a user.
func (a *Activities) ShipmentCreatedNotification(ctx context.Context, input ShipmentCreatedNotificationInput) error {
	err := a.sendMail(from, to,
		fmt.Sprintf("Shipment for order: %s", string(input.OrderID)),
		"Your order has been processed and shipping has been arranged with the courier. We'll be in touch once its dispatched.",
	)

	return err
}

// ShipmentDispatchedNotificationInput is the input for a ShipmentDispatched notification.
type ShipmentDispatchedNotificationInput struct {
	OrderID ordersapi.OrderID
}

// ShipmentDispatchedNotification sends a ShipmentDispatched notification to a user.
func (a *Activities) ShipmentDispatchedNotification(ctx context.Context, input ShipmentDispatchedNotificationInput) error {
	err := a.sendMail(from, to,
		fmt.Sprintf("Shipment dispatched for order: %s", string(input.OrderID)),
		"Your order has been dispatched.",
	)

	return err
}

// ShipmentDeliveredNotificationInput is the input for a ShipmentDelivered notification.
type ShipmentDeliveredNotificationInput struct {
	OrderID ordersapi.OrderID
}

// ShipmentDeliveredNotification sends a ShipmentDelivered notification to a user.
func (a *Activities) ShipmentDeliveredNotification(ctx context.Context, input ShipmentDeliveredNotificationInput) error {
	err := a.sendMail(from, to,
		fmt.Sprintf("Shipment delivered for order: %s", string(input.OrderID)),
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
