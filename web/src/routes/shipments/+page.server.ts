import type { Shipment } from '$lib/types/order';

export const load = async () => {
	const response = await fetch(`http://localhost:8081/shipments`);
	const shipments: Shipment[] = await response.json();

	return { shipments };
};
