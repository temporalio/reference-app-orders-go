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
		<h1>Orders</h1>
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
				<tr>
					<td><a href={`/orders/${order.id}`}>{order.id}</a></td>
					<td
						>{new Date(order.startedAt).toLocaleDateString('en-US', {
							weekday: 'short',
							year: 'numeric',
							month: 'long',
							day: 'numeric',
							hour: 'numeric',
							minute: 'numeric',
							second: 'numeric'
						})}</td
					>
				</tr>
			{/each}
		</tbody>
	</table>
</section>
