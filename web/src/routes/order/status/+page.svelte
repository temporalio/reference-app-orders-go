<script lang="ts">
	import { goto } from "$app/navigation";
	import OrderDetails from "$lib/components/order-details.svelte";
	import OrderStatus from "$lib/components/order-status.svelte";
	import { order, type Order } from "$lib/stores/order";
</script>

<svelte:head>
	<title>Tora</title>
	<meta name="description" content="Tora App" />
</svelte:head>

<section>
	<div class="container">
		<OrderDetails order={$order} />
	</div>
	<OrderStatus order={$order} />
	<div class="container submit">
    <button on:click={() => goto('/order/select')}>Back</button>
		<button disabled={!$order} class:disabled={!$order}>Submit</button>
	</div>
</section>

<style>
	section {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}
	.container {
		display: flex;
		flex-direction: row;
		gap: 2rem;
		width: 100%;
	}

	.submit {
		justify-content: end;
	}

	button {
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
