import { json, type RequestHandler } from '@sveltejs/kit';
import { SHIPMENT_API_URL } from '$env/static/private';

export const POST: RequestHandler = async ({ request }) => {
	const { signal, shipment } = await request.json();

	try {
		const response = await fetch(`${SHIPMENT_API_URL}/shipments/${shipment.id}/status`, {
			method: 'POST',
			body: JSON.stringify({ status: signal.status })
		});
		return json({ status: 'ok', body: response });
	} catch (e) {
		return json({ status: 'error' });
	}
};
