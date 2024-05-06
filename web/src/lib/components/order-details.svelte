<script lang="ts">
	import { type Order } from '$lib/types/order';

	export let order: Order | undefined;
</script>

<div class="details">
	{#if order}
		{#each order.items as item}
			<div class="item">
				<div class="item-header">
					<div class="quantity">{item.quantity}</div>
					<h3 class="name">{item.sku}</h3>
				</div>
				<div class="description">{item.description}</div>
			</div>
		{/each}
	{:else}
		<h3>Pick order to view details</h3>
	{/if}
</div>

<style>
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

	.details {
		display: flex;
		flex-direction: column;
		gap: 2rem;
		background-color: white;
		padding: 2rem;
		align-items: start;
		width: 66vw;
	}
	@media (max-width: 640px) {
		.details {
			width: 100%;
			max-width: fit-content;
		}
	}
</style>
