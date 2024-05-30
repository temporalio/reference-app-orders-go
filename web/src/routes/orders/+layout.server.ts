import { ORDER_API_URL } from '$env/static/private';

export const load = async () => {
	const response = await fetch(`${ORDER_API_URL}/orders`);
	const orders = await response.json();

	return { orders };
};
