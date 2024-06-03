<script lang="ts">
	import { page } from '$app/stores';
	import ShipmentDetails from '$lib/components/shipment-details.svelte';
	import type { Shipment } from '$lib/types/order';

	$: ({ shipment } = $page.data);

	$: status = shipment?.status;

	let broadcaster: BroadcastChannel;
	$: {
		if (shipment?.id) {
			broadcaster = new BroadcastChannel(`shipment-${shipment.id}`);
		}
	}

	const dispatchShipment = async (shipment: Shipment) => {
		const signal = { name: 'ShipmentUpdate', status: 'dispatched' };
		await fetch('/api/shipment', { method: 'POST', body: JSON.stringify({ shipment, signal }) });
		status = 'dispatched';
		broadcaster?.postMessage(status);
	};

	const deliverShipment = async (shipment: Shipment) => {
		const signal = { name: 'ShipmentUpdate', status: 'delivered' };
		await fetch('/api/shipment', { method: 'POST', body: JSON.stringify({ shipment, signal }) });
		status = 'delivered';
		broadcaster?.postMessage(status);
	};
</script>

<svelte:head>
	<title>OMS</title>
	<meta name="description" content="OMS App" />
</svelte:head>

<div class="container">
	<h1>{shipment.id}</h1>
	<ShipmentDetails {shipment} {status} />
	<div class="action-btns">
		<button
			class="submit"
			disabled={status !== 'booked'}
			on:click={() => dispatchShipment(shipment)}>Dispatch</button
		><button
			class="submit"
			disabled={status !== 'dispatched'}
			on:click={() => deliverShipment(shipment)}>Deliver</button
		>
	</div>
</div>

<style>
	.container {
		display: flex;
		flex-direction: column;
		align-items: start;
		gap: 2rem;
		background-color: white;
		padding: 2rem;
		border-radius: 0.5rem;
	}
	.action-btns {
		display: flex;
		justify-content: end;
		gap: 0.5rem;
		width: 100%;
	}
</style>
