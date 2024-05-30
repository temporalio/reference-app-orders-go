import { ORDER_API_URL } from '$env/static/private';

export const load = async ({ params }) => {
	const { id } = params;

	const orderResponse = await fetch(`${ORDER_API_URL}/orders/${id}`);
	const order = await orderResponse.json();

	return { order };
};
