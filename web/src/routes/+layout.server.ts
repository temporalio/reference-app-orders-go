export const load = async ({ depends }) => {
	const response = await fetch(`http://localhost:8083/orders`);
	const orders = await response.json();

	return { orders };
};
