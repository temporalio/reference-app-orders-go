<script lang="ts">
	import { type Order } from '$lib/stores/order';
	import ShipmentStatus from './shipment-status.svelte';

	export let order: Order;

	$: fulfillments = order?.fulfillments || [];
</script>

<div class="details">
	{#each fulfillments as fulfillment}
		<div class="fulfillment">
			<p class="location">{fulfillment.Location}</p>
			<ShipmentStatus shipment={fulfillment.shipment} />
		</div>
		{#each fulfillment.items as item}
			<div class="item">
				<div class="item-header">
					<div class="quantity">{item.quantity}</div>
					<h3 class="name">{item.sku}</h3>
				</div>
			</div>
		{/each}
	{/each}
</div>

<style>
	.fulfillment {
		display: flex;
		align-items: center;
		justify-content: space-between;
		width: 100%;
	}
	.item {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.item-header {
		display: flex;
		gap: 1rem;
		align-items: center;
	}

	.location {
		font-size: 24px;
		font-weight: 700;
	}
	.name {
		font-weight: 500;
	}

	.details {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		width: 100%;
		align-items: start;
	}
</style>
