<script lang="ts">
	import type { OrderItem } from '$lib/stores/order';

	export let items: OrderItem[];
</script>

<div class="container">
	{#each items as item}
		<div class="item">
			<div class="item-header">
				<div class="quantity">{item.quantity}</div>
				<h3 class="name">{item.sku}</h3>
			</div>
		</div>
	{/each}
</div>

<style>
	.container {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		width: 100%;
		align-items: start;
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
	.name {
		font-weight: 700;
	}
</style>
