import { json, type RequestHandler } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const POST: RequestHandler = async ({ request }) => {
	const { order } = await request.json();

	try {
		const response = await fetch(`${env.ORDER_API_URL}/orders`, {
			method: 'POST',
			body: JSON.stringify(order)
		});
		return json({ status: 'ok', body: response });
	} catch (e) {
		return json({ status: 'error' });
	}
};
