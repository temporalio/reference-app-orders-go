<script lang="ts">
	import { goto } from '$app/navigation';
	import OrderDetails from '$lib/components/order-details.svelte';
	import { generateOrders, order, type Order } from '$lib/stores/order';


	const onItemClick = async (order: Order) => {
		if (order.OrderId === $order?.OrderId) {
			$order = undefined;
		} else {
			$order = order;
		};
	}

	const onSubmit = async () => {
		if ($order) {
			await fetch('/api/order', { method: 'POST', body: JSON.stringify({ order: $order  }) })
			goto(`/order/${$order.OrderId}/status`);
		}
	}

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
					class:active={_order.OrderId === $order?.OrderId}
					on:click={() => onItemClick(_order)}
				>
					<div class="name">Package {index + 1}</div>
				</button>
			{/each}
		</div>
		<OrderDetails order={$order} />
	</div>
	<div class="container submit">
		<button
			class="submit-button"
			disabled={!$order}
			class:disabled={!$order}
			on:click={onSubmit}>Submit</button
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
		width: 33%;
		height: 30rem;
		overflow: auto;
	}

	.item {
		height: 2rem;
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
	}

	.active {
		background-color: black;
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

	.submit-button {
		padding: 1rem 2rem;
		background-color: black;
		color: white;
		border: none;
		cursor: pointer;
	}

	.disabled {
		background-color: #ccc;
		color: #666;
		cursor: not-allowed;
	}
</style>
