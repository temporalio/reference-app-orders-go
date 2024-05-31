<script lang="ts">
	import { type Fulfillment, type Order } from '$lib/types/order';
	import ItemDetails from './item-details.svelte';
	import PaymentDetails from './payment-details.svelte';
	import ShipmentStatus from './shipment-status.svelte';

	export let order: Order;

	$: fulfillments = order?.fulfillments?.filter((f) => f.status != 'cancelled') || [];

	$: getStatus = (fulfillment: Fulfillment) => {
		if (!fulfillment.shipment) return fulfillment.status;
		return fulfillment.shipment.status;
	};
</script>

<div class="details">
	{#each fulfillments as fulfillment}
		<div class="container">
			<div class="fulfillment">
				<p class="location">
					{#if fulfillment?.location}
						<i>{fulfillment.location}</i>
					{:else}
						<strong>Action Required</strong>
					{/if}
				</p>
				<ShipmentStatus id={fulfillment?.shipment?.id} status={getStatus(fulfillment)} />
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
		align-items: end;
		justify-content: space-between;
		width: 100%;
		border-bottom: 2px solid black;
		margin-bottom: 1rem;
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
		padding: 0.5rem;
		border-radius: 0.15rem;
	}
	.location {
		font-size: 1.5rem;
	}

	.details {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		width: 100%;
		align-items: start;
	}
</style>
