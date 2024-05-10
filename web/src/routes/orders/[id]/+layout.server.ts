export const load = async (request) => {
	const { id } = request.params;

	const orderResponse = await fetch(`http://127.0.0.1:8083/orders/${id}`);
	const order = await orderResponse.json();

	return { order };
};
