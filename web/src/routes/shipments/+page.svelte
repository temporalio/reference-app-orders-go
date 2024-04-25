<script lang="ts">
	import type { Shipment } from '$lib/stores/order';

	export let data;

	$: ({ shipments } = data);

	let shipmentStatuses: Record<string, string> = {};
	$: {
		shipmentStatuses = shipments.reduce((acc, shipment) => {
			acc[shipment.id] = shipment.status;
			return acc;
		}, {});
	}

	const dispatchShipment = async (shipment: Shipment) => {
		const signal = { name: 'ShipmentUpdate', status: 'dispatched' };
		await fetch('/api/shipment', { method: 'POST', body: JSON.stringify({ shipment, signal }) });
		shipmentStatuses[shipment.id] = 'dispatched';
	};

	const deliverShipment = async (shipment: Shipment) => {
		const signal = { name: 'ShipmentUpdate', status: 'delivered' };
		await fetch('/api/shipment', { method: 'POST', body: JSON.stringify({ shipment, signal }) });
		shipmentStatuses[shipment.id] = 'delivered';
	};
</script>

<svelte:head>
	<title>Tora</title>
	<meta name="description" content="Tora App" />
</svelte:head>

<section>
	<nav>
		<h1>Shipments</h1>
	</nav>
	<table>
		<thead>
			<tr>
				<th>Order ID</th>
				<th>Status</th>
				<th></th>
			</tr>
		</thead>
		<tbody>
			{#each shipments as shipment}
				<tr>
					<td style="width: 100%;">{shipment.id}</td>
					<td><div class="status">{shipmentStatuses[shipment.id].toUpperCase()}</div></td>
					<td
						><div class="action-btns">
							<button on:click={() => dispatchShipment(shipment)}>Dispatch</button><button
								on:click={() => deliverShipment(shipment)}>Deliver</button
							>
						</div></td
					>
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
		padding: 0.25rem 1rem;
		text-align: left;
		background-color: black;
		color: white;
		white-space: nowrap;
	}

	tr {
		cursor: pointer;
	}

	td {
		background-color: #f2f2f2;
		color: black;
	}

	.status {
		position: relative;
		display: inline-flex;
		border: 3px solid black;
		border-radius: 9999px;
		padding: 0.25rem 1rem;
		margin: 0 auto;
		background-color: lightgreen;
	}

	.action-btns {
		display: flex;
		justify-content: end;
		gap: 0.5rem;
	}
</style>
