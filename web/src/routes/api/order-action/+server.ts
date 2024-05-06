import { json, type RequestHandler } from '@sveltejs/kit';

export const POST: RequestHandler = async ({ request }) => {
	const { id, action } = await request.json();

	try {
		const response = await fetch(`http://localhost:8083/orders/${id}/${action}`, {
			method: 'POST',
			body: JSON.stringify(action)
		});
		return json({ status: 'ok', body: response });
	} catch (e) {
		return json({ status: 'error' });
	}
};
