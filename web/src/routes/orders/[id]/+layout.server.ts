export const load = async ({ params }) => {
	const { id } = params;

	const orderResponse = await fetch(`http://localhost:8083/orders/${id}`);
	const order = await orderResponse.json();

	return { order };
};
