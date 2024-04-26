import type { Order, Fulfillment } from '$lib/stores/order';

export const load = async () => {
	const response = await fetch(`http://localhost:8081/shipments`);
	const shipments = await response.json();

	return { shipments };
};
