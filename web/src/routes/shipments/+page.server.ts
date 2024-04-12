export const load = async () => {
	const response = await fetch(
		`http://localhost:8234/api/v1/namespaces/default/workflows?query=WorkflowType="Shipment"`
	);
	const workflows = await response.json();

	return { shipments: workflows.executions };
};
