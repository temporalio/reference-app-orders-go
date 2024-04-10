<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import OrderDetails from '$lib/components/order-details.svelte';
	import OrderStatus from '$lib/components/order-status.svelte';
	import ShipmentDetails from '$lib/components/shipment-details.svelte';

	$: ({ order, shipments } = $page.data);
</script>

<svelte:head>
	<title>Tora</title>
	<meta name="description" content="Tora App" />
</svelte:head>

<section>
	<div class="container">
		<OrderDetails {order} />
		<ShipmentDetails {shipments} />
	</div>
	<OrderStatus {order} />
	<div class="container submit">
		<button on:click={() => goto('/order/select')}>Back</button>
		<button disabled class="disabled">Approve</button>
	</div>
</section>

<style>
	section {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}
	.container {
		display: flex;
		flex-direction: row;
		gap: 2rem;
		width: 100%;
	}

	.submit {
		justify-content: end;
	}
</style>
