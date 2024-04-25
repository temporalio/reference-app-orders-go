import type { Order, Fulfillment } from '$lib/stores/order';

export const load = async () => {
	const response = await fetch(`http://localhost:8083/orders`);
	const orders = await response.json();
	const orderIds = orders.map((order: Order) => order.id);

	const shipments = await Promise.all(
		orderIds.map(async (orderId: string) => {
			const response = await fetch(`http://localhost:8083/order/${orderId}`);
			const order = await response.json();
			const orderShipments = order.fulfillments?.map((fullment: Fulfillment) => {
				return fullment.shipment;
			});
			return orderShipments;
		})
	);
	return { shipments: shipments.flat() };
};
