import { json, type RequestHandler } from '@sveltejs/kit';
import { ORDER_API_URL } from '$env/static/private';

export const POST: RequestHandler = async ({ request }) => {
	const { order } = await request.json();

	try {
		const response = await fetch(`${ORDER_API_URL}/orders`, {
			method: 'POST',
			body: JSON.stringify(order)
		});
		return json({ status: 'ok', body: response });
	} catch (e) {
		return json({ status: 'error' });
	}
};
