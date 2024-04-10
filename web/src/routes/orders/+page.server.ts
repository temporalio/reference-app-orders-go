export const load = async () => {
	const response = await fetch(
		`http://localhost:8234/api/v1/namespaces/default/workflows?query=WorkflowType="Order"`
	);
	const workflows = await response.json();

	return { orders: workflows.executions };
};
