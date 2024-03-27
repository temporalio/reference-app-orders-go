<script lang="ts">
	import { goto } from "$app/navigation";
	import { generateOrders, items, order, type Order } from "$lib/stores/order";

	const onItemClick = (order: Order) => {
		if (order.id === $order?.id) {
			$order = undefined;
		} else {
			$order = order;
		}
	}

  const orders = generateOrders(20)
</script>

<svelte:head>
	<title>Tora</title>
	<meta name="description" content="Tora App" />
</svelte:head>

<section>
	<div class="container">
		<div class="list">
			{#each orders as order}
				<button class="item" class:active={false} on:click={() => onItemClick(order)}>
					<div class="name">{order.id}</div>
        </button>
			{/each}
		</div>
		<div class="details">
			{#if $order}
				<h3 class="name">{$order.id}</h3>
        {#each $order.items as item}
          <div class="description">{item.name}</div>
          <div class="description">{item.description}</div>
          <div class="description">{item.quantity}</div>
        {/each}
			{:else}
			<h3>Pick order to view details</h3>
			{/if}
		</div>
	</div>
	<div class="container submit">
		<button class="submit-button" disabled={!$order} class:disabled={!$order} on:click={() => goto('/order-status')}>Submit</button>
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
