import { SHIPMENT_API_URL } from '$env/static/private';

export const load = async (request) => {
	const { id } = request.params;

	const response = await fetch(`${SHIPMENT_API_URL}/shipments/${id}`);
	const shipment = await response.json();

	return { shipment };
};
