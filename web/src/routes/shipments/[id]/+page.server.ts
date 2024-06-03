import { env } from '$env/dynamic/private';

export const load = async (request) => {
	const { id } = request.params;

	const response = await fetch(`${env.SHIPMENT_API_URL}/shipments/${id}`);
	const shipment = await response.json();

	return { shipment };
};
