<script lang="ts">
	import { type Shipment } from '$lib/stores/order';

	export let shipments: Shipment[] = [];

	const dispatchShipment = async (shipment: Shipment) => {
		const signal = { name: 'ShipmentUpdate', status: 1 };
		await fetch('/api/shipment', { method: 'POST', body: JSON.stringify({ shipment, signal }) });
	};

	const deliverShipment = async (shipment: Shipment) => {
		const signal = { name: 'ShipmentUpdate', status: 2 };
		await fetch('/api/shipment', { method: 'POST', body: JSON.stringify({ shipment, signal }) });
	};
</script>

<div class="details">
	<h1 class="title">Shipments</h1>
	{#each shipments as shipment}
		<div class="shipment">
			<h3 class="name">{shipment.id}</h3>
			<button on:click={() => dispatchShipment(shipment)}>Dispatch</button>
			<button on:click={() => deliverShipment(shipment)}>Deliver</button>
		</div>
	{:else}
		<div class="shipment">
			<h3 class="name">Completed</h3>
		</div>
	{/each}
</div>

<style>
	.title {
		font-size: 1.5rem;
		text-decoration: underline;
	}

	.shipment {
		display: flex;
		width: 100%;
		justify-content: space-between;
		align-items: center;
	}

	.name {
		font-weight: bold;
		font-size: 1rem;
	}

	.details {
		display: flex;
		flex-direction: column;
		gap: 2rem;
		width: 66%;
		background-color: white;
		padding: 2rem;
		align-items: start;
	}

	.details .name {
		font-size: 1.5rem;
	}
</style>
