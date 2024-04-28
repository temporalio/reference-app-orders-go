<script lang="ts">
	import { items, type Order } from '$lib/stores/order';
	import ItemDetails from './item-details.svelte';
	import PaymentDetails from './payment-details.svelte';
	import ShipmentStatus from './shipment-status.svelte';

	export let order: Order;

	$: fulfillments = order?.fulfillments || [];
</script>

<div class="details">
	{#each fulfillments as fulfillment}
	<div class="container">
		<div class="fulfillment">
			<p class="location">{fulfillment.location}</p>
			<ShipmentStatus shipment={fulfillment.shipment} />
		</div>
		<ItemDetails items={fulfillment.items} />
		{#if fulfillment.payment}
		<PaymentDetails payment={fulfillment.payment} />
		{/if}
		</div>
	{/each}
</div>

<style>


.fulfillment {
	display: flex;
	align-items: center;
	justify-content: space-between;
	width: 100%;
}

@media (max-width: 640px) {
	.fulfillment {
		flex-direction: column;
	}
}

	.container {
		display: flex;
		flex-direction: column;
		align-items: start;
		justify-content: space-between;
		width: 100%;
		border: 3px solid black;
		padding: .5rem;
		border-radius: .15rem;
	}
	.location {
		font-size: 1.5rem;
		font-weight: 700;
	}

	.details {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		width: 100%;
		align-items: start;
	}
</style>
