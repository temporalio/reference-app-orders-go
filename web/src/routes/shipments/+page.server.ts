import type { Shipment } from '$lib/types/order';
import { SHIPMENT_API_URL } from '$env/static/private';

export const load = async () => {
	const response = await fetch(`${SHIPMENT_API_URL}/shipments`);
	const shipments: Shipment[] = await response.json();

	return { shipments };
};
