import type { Shipment } from '$lib/types/order';
import { env } from '$env/dynamic/private';

export const load = async () => {
	const response = await fetch(`${env.SHIPMENT_API_URL}/shipments`);
	const shipments: Shipment[] = await response.json();

	return { shipments };
};
