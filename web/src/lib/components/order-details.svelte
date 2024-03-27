<script lang="ts">
	import { type Order } from '$lib/stores/order';

	export let order: Order | undefined;
</script>

<div class="details">
	{#if order}
		{#each order.items as item}
			<div class="title">
				<div class="quantity">{item.quantity}</div>
				<h3 class="name">{item.name}</h3>
			</div>
			<div class="description">{item.description}</div>
		{/each}
	{:else}
		<h3>Pick order to view details</h3>
	{/if}
</div>

<style>
	.title {
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
		gap: 1rem;
		width: 66%;
		background-color: white;
		padding: 2rem;
	}

	.description {
		font-size: 1rem;
	}
	.details .name {
		font-size: 1.5rem;
	}
</style>
