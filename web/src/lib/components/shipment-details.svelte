<script lang="ts">
	import { type Order, type Shipment } from '$lib/stores/order';

	export let order: Order;

	$: shipments = order?.fulfillments?.map((fulfillment) => fulfillment.shipment) || [];

	const dispatchShipment = async (shipment: Shipment) => {
		const signal = { name: 'ShipmentUpdate', status: 'dispatched' };
		await fetch('/api/shipment', { method: 'POST', body: JSON.stringify({ shipment, signal }) });
	};

	const deliverShipment = async (shipment: Shipment) => {
		const signal = { name: 'ShipmentUpdate', status: 'delivered' };
		await fetch('/api/shipment', { method: 'POST', body: JSON.stringify({ shipment, signal }) });
	};
</script>

<div class="details">
	<h1 class="title">Shipments</h1>
	{#each shipments as shipment}
		<div class="shipment">
			{#if shipment}
				<h3>{shipment.id}</h3>
				<p>Status: <span class="status">{shipment.status}</span></p>
				<button on:click={() => dispatchShipment(shipment)}>Dispatch</button>
				<button on:click={() => deliverShipment(shipment)}>Deliver</button>
			{:else}
				<h3>Shipment not created</h3>
			{/if}
		</div>
	{:else}
		<div class="shipment">
			<h3>Waiting for shipments</h3>
		</div>
	{/each}
</div>

<style>
	.title {
		font-size: 1.5rem;
	}

	.details {
		display: flex;
		width: 50%;
		flex-direction: column;
		gap: 2rem;
		background-color: white;
		padding: 2rem;
		align-items: start;
	}

	.details .name {
		font-size: 1.5rem;
	}
</style>
