import { json, type RequestHandler } from '@sveltejs/kit';

export const POST: RequestHandler = async ({ request }) => {
	const { order } = await request.json();

	try {
		const response = await fetch(`http://127.0.0.1:8083/orders`, {
			method: 'POST',
			body: JSON.stringify(order)
		});
		return json({ status: 'ok', body: response });
	} catch (e) {
		return json({ status: 'error' });
	}
};
