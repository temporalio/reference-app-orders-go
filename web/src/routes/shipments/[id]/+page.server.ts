export const load = async (request) => {
	const { id } = request.params;

	const response = await fetch(`http://localhost:8081/shipments/${id}`);
	const shipment = await response.json();

	return { shipment };
};
