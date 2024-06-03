import { json, type RequestHandler } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const POST: RequestHandler = async ({ request }) => {
	const { id, action } = await request.json();

	try {
		const response = await fetch(`${env.ORDER_API_URL}/orders/${id}/action`, {
			method: 'POST',
			body: JSON.stringify({ action })
		});
		return json({ status: 'ok', body: response });
	} catch (e) {
		return json({ status: 'error' });
	}
};
