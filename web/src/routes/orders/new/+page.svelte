<script lang="ts">
	import { goto } from '$app/navigation';
	import OrderDetails from '$lib/components/order-details.svelte';
	import { generateOrders, order, type Order } from '$lib/stores/order';

	const onItemClick = async (order: Order) => {
		if (order.id === $order?.id) {
			$order = undefined;
		} else {
			$order = order;
		}
	};

	const onSubmit = async () => {
		if ($order) {
			await fetch('/api/order', { method: 'POST', body: JSON.stringify({ order: $order }) });
			// TODO(Rob): Add a status field to say an order is pending, or similar, so we know to refresh
			goto(`/orders/Order:${$order.id}`);
		}
	};

	const orders = generateOrders(20);
</script>

<svelte:head>
	<title>Tora</title>
	<meta name="description" content="Tora App" />
</svelte:head>

<section>
	<div class="container">
		<div class="list">
			{#each orders as _order, index}
				<button
					class="item"
					class:active={_order.id === $order?.id}
					on:click={() => onItemClick(_order)}
				>
					<div class="name">Order {index + 1}</div>
				</button>
			{/each}
		</div>
		<OrderDetails order={$order} />
	</div>
	<div class="container submit">
		<button class="submit-button" disabled={!$order} on:click={onSubmit}
			>Submit</button
		>
	</div>
</section>

<style>
	.container {
		display: flex;
		flex-direction: row;
		gap: 2rem;
		width: 100%;
	}	

	.submit {
		margin: 1rem 0;
		justify-content: end;
	}

	.list {
		background-color: white;
		width: 20vw;
		height: 55vh;
		overflow: auto;
		border: 3px solid black;
		border-radius: .25rem;
	}

	@media (max-width: 640px) {
		.container {
			flex-direction: column;
		}

		.list {
			width: 100%;
			height: 15rem;
		}
	}
	

	.item {
		height: 2.25rem;
		width: 100%;
		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;
		cursor: pointer;
		border-radius: 0;
		border: none;
		border-bottom: 2px solid #ccc;
		background-color: white;
		color: black;
		text-transform: uppercase;
	}

	.active {
		background-color: var(--color-theme-2);
		color: white;
	}

	.name {
		font-weight: bold;
	}
</style>
