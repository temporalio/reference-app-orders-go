export const load = async (request) => {
	const { id } = request.params;
	const workflowResponse = await fetch(
		`http://localhost:8234/api/v1/namespaces/default/workflows/${id}`
	);
	const workflow = await workflowResponse.json();

	const historyResponse = await await fetch(
		`http://localhost:8234/api/v1/namespaces/default/workflows/${id}/history`
	);
	const { history } = await historyResponse.json();
	const order = history?.events[0]?.workflowExecutionStartedEventAttributes?.input[0];

	const childWorkflowIds = workflow?.pendingChildren?.map((child: any) => child.workflowId) ?? [];
	const childWorkflows = await Promise.all(
		childWorkflowIds.map(async (id: string) => {
			const response = await fetch(
				`http://localhost:8234/api/v1/namespaces/default/workflows/${id}`
			);
			const childWorkflow = await response.json();
			return childWorkflow;
		})
	);
	const shipments = childWorkflows.map((child: any) => {
		return { name: child.workflowExecutionInfo.execution.workflowId };
	});

	return { workflow, order, shipments };
};
