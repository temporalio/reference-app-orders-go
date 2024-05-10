export const load = async (request) => {
	const { id } = request.params;

	const response = await fetch(`http://127.0.0.1:8081/shipments/${id}`);
	const shipment = await response.json();

	return { shipment };
};
