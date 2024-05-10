<script lang="ts">
	import type { OrderItem } from '$lib/types/order';

	export let items: OrderItem[];
</script>

<div class="container">
	{#each items as item}
		<div class="item">
			<div class="item-header">
				<div class="quantity">{item.quantity}</div>
				<p>{item.sku}</p>
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
</style>
