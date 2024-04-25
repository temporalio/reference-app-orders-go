<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import FulfillmentDetails from '$lib/components/fulfillment-details.svelte';

	$: ({ order } = $page.data);
</script>

<svelte:head>
	<title>Tora</title>
	<meta name="description" content="Tora App" />
</svelte:head>

<section>
	<div class="container">
		<h1>{order.id}</h1>
		<FulfillmentDetails {order} />
		<div class="submit">
			<button on:click={() => goto('/orders')}>Back to Orders</button>
		</div>
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
		flex-direction: column;
		align-items: start;
		gap: 2rem;
		width: 100%;
		background-color: white;
		padding: 2rem;
	}

	.submit {
		width: 100%;
		display: flex;
		justify-content: end;
	}
</style>
