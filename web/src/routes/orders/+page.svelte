<script lang="ts">
	import { goto } from '$app/navigation';

	export let data;

	$: ({ orders } = data);
</script>

<svelte:head>
	<title>Tora</title>
	<meta name="description" content="Tora App" />
</svelte:head>

<section>
	<nav>
		<h1>Order History</h1>
		<button on:click={() => goto('/orders/new')}>New Order</button>
	</nav>
	<table>
		<thead>
			<tr>
				<th>Order ID</th>
				<th>Date</th>
			</tr>
		</thead>
		<tbody>
			{#each orders as order}
				<tr on:click={() => goto(`/orders/${order.id}/status`)}>
					<td>{order.id}</td>
					<td>{order.startedAt}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</section>

<style>
	nav {
		display: flex;
		justify-content: space-between;
		margin-bottom: 2rem;
	}

	table {
		width: 100%;
		border-collapse: collapse;
	}

	th,
	td {
		border: 1px solid black;
		padding: 8px;
		text-align: left;
		background-color: black;
		color: white;
	}

	tr {
		cursor: pointer;
	}

	td {
		background-color: #f2f2f2;
		color: black;
	}
</style>
