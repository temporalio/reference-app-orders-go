import { env } from '$env/dynamic/private';

export const load = async ({ params }) => {
	const { id } = params;

	const orderResponse = await fetch(`${env.ORDER_API_URL}/orders/${id}`);
	const order = await orderResponse.json();

	return { order };
};
