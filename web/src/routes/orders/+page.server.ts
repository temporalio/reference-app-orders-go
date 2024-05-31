export const load = async () => {
	const response = await fetch(`http://127.0.0.1:8083/orders`);
	const orders = await response.json();

	return { orders };
};
