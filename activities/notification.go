package activities

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

const from = "orders@reference-app.example"
const to = "customer@reference-app.example"

type ShipmentCreatedNotificationInput struct {
	OrderID ordersapi.OrderID
}

func (a *Activities) ShipmentCreatedNotification(ctx context.Context, input ShipmentCreatedNotificationInput) error {
	err := a.sendMail(from, to,
		fmt.Sprintf("Shipment for order: %s", string(input.OrderID)),
		"Your order has been processed and shipping has been arranged with the courier. We'll be in touch once its dispatched.",
	)

	return err
}

type ShipmentDispatchedNotificationInput struct {
	OrderID ordersapi.OrderID
}

func (a *Activities) ShipmentDispatchedNotification(ctx context.Context, input ShipmentDispatchedNotificationInput) error {
	err := a.sendMail(from, to,
		fmt.Sprintf("Shipment dispatched for order: %s", string(input.OrderID)),
		"Your order has been dispatched.",
	)

	return err
}

type ShipmentDeliveredNotificationInput struct {
	OrderID ordersapi.OrderID
}

func (a *Activities) ShipmentDeliveredNotification(ctx context.Context, input ShipmentDeliveredNotificationInput) error {
	err := a.sendMail(from, to,
		fmt.Sprintf("Shipment delivered for order: %s", string(input.OrderID)),
		"Your order has been delivered.",
	)

	return err
}

func (a *Activities) sendMail(from string, to string, subject string, body string) error {
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", a.SMTPHost, a.SMTPPort),
		nil,
		from,
		[]string{to},
		[]byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)),
	)
}
