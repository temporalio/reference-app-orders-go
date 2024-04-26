export const load = async (request) => {
	const { id } = request.params;

	const orderResponse = await fetch(`http://localhost:8083/orders/${id}`);
	const order = await orderResponse.json();
	const shipments =
		order?.fulfillments?.map((f: any) => {
			return f.shipment;
		}) ?? [];

	return { order, shipments };
};
