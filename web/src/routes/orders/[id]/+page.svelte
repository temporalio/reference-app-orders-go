<script lang="ts">
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
		<h1>{$page.params.id}</h1>
		<FulfillmentDetails {order} />
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
		background-color: white;
		padding: 2rem;
		border-radius: 0.5rem;
	}
</style>
