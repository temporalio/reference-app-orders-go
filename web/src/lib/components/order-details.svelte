<script lang="ts">
	import { type Order } from '$lib/stores/order';

	export let order: Order | undefined;
</script>

<div class="details">
	<h1 class="title">Order</h1>
	{#if order}
		{#each order.Items as item}
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
	.title {
		font-size: 1.5rem;
		text-decoration: underline;
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

	.quantity {
		font-size: 0.85rem;
		border-radius: 9999px;
		background-color: #788f9c;
		width: 2rem;
		height: 2rem;
		display: flex;
		justify-content: center;
		align-items: center;
		color: white;
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

	.description {
		font-size: 1rem;
	}
	.details .name {
		font-size: 1.5rem;
	}
</style>
