<script lang="ts">
	export let data;

	$: ({ shipments } = data);
</script>

<svelte:head>
	<title>Tora</title>
	<meta name="description" content="Tora App" />
</svelte:head>

<section>
	<nav>
		<h1>Shipment History</h1>
	</nav>
	<table>
		<thead>
			<tr>
				<th>Order ID</th>
				<th>Date</th>
			</tr>
		</thead>
		<tbody>
			{#each shipments as shipment}
				<tr>
					<td>{shipment.execution.workflowId}</td>
					<td>{shipment.startTime}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</section>

<style>
	nav {
		display: flex;
		justify-content: space-between;
		margin-bottom: 2rem;
	}

	table {
		width: 100%;
		border-collapse: collapse;
	}

	th,
	td {
		border: 1px solid black;
		padding: 8px;
		text-align: left;
		background-color: black;
		color: white;
	}

	tr {
		cursor: pointer;
	}

	td {
		background-color: #f2f2f2;
		color: black;
	}
</style>
