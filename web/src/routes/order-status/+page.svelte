<script lang="ts">
	import { goto } from "$app/navigation";
	import { order, type Order } from "$lib/stores/order";
</script>

<svelte:head>
	<title>Tora</title>
	<meta name="description" content="Tora App" />
</svelte:head>

<section>
	<div class="container">
		<div class="details">
			{#if $order}
				<h3 class="name">{$order.name}</h3>
				<div class="description">{$order.description}</div>
			{:else}
			<h3>Pick order to view details</h3>
			{/if}
		</div>
	</div>
	<div class="container submit">
    <button on:click={() => goto('/select-order')}>Back</button>
		<button disabled={!$order} class:disabled={!$order}>Submit</button>
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
		border-bottom: 1px solid #ccc;
		text-align: center;
		height: 2rem;
		display: flex;
		flex-direction: column;
		justify-content: center;
		cursor: pointer;
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
