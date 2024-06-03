<script lang="ts">
	import { goto } from '$app/navigation';
	import OrderDetails from '$lib/components/order-details.svelte';
	import { generateOrders, order, type Order } from '$lib/types/order';
	import Logo from '$lib/components/logo.svelte';

	const onItemClick = async (order: Order) => {
		if (order.id === $order?.id) {
			$order = undefined;
		} else {
			$order = order;
		}
	};

	const onSubmit = async () => {
		if ($order) {
			loading = true;
			await fetch('/api/order', { method: 'POST', body: JSON.stringify({ order: $order }) });
			setTimeout(() => {
				goto(`/orders/${$order.id}`, { invalidateAll: true });
				$order = undefined;
			}, 1000);
		}
	};

	const orders = generateOrders(20);

	let loading = false;
</script>

<svelte:head>
	<title>OMS</title>
	<meta name="description" content="OMS App" />
</svelte:head>

<section>
	{#if loading}
		<div class="container loading"><Logo loading /></div>
	{:else}
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
			<button class="submit-button" disabled={!$order} on:click={onSubmit}>Submit</button>
		</div>
	{/if}
</section>

<style>
	.container {
		display: flex;
		flex-direction: row;
		justify-content: center;
		gap: 2rem;
		width: 100%;
	}

	.loading {
		align-items: center;
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
		border-radius: 0.25rem;
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
		height: 2.5rem;
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
	}

	.active {
		background-color: var(--color-theme-2);
		color: white;
	}

	.name {
		font-weight: bold;
	}
</style>
