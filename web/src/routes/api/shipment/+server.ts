import { json, type RequestHandler } from '@sveltejs/kit';

export const POST: RequestHandler = async ({ request }) => {
	const { signal, shipment } = await request.json();

	try {
		const response = await fetch(
			`http://localhost:8234/api/v1/namespaces/default/workflows/${shipment.name}/signal/${signal.name}`,
			{
				method: 'POST',
				body: JSON.stringify({
					input: [{ Status: signal.status }]
				})
			}
		);
		return json({ status: 'ok', body: response });
	} catch (e) {
		return json({ status: 'error' });
	}
};
