import { json, type RequestHandler } from '@sveltejs/kit';

export const POST: RequestHandler = async ({ request }) => {
	const { order } = await request.json();

	try {
		const response = await fetch(
			`http://localhost:8234/api/v1/namespaces/default/workflows/${order.OrderId}`,
			{
				method: 'POST',
				body: JSON.stringify({
					taskQueue: { name: 'order' },
					workflowType: { name: 'Order' },
					input: [order]
				})
			}
		);
		return json({ status: 'ok', body: response });
	} catch (e) {
		return json({ status: 'error' });
	}
};
