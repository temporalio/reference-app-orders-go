import { env } from '$env/dynamic/private';

export const load = async () => {
	const response = await fetch(`${env.ORDER_API_URL}/orders`);
	const orders = await response.json();

	return { orders };
};
