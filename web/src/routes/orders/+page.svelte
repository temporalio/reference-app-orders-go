<script lang="ts">
	import { goto } from '$app/navigation';
	import StatusBadge from '$lib/components/status-badge.svelte';

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
				<th>ID</th>
				<th style="text-align: center;">Date & Time</th>
				<th style="text-align: center;">Status</th>
			</tr>
		</thead>
		<tbody>
			{#each orders as order}
				<tr>
					<td style="width: 100%;"><a href={`/orders/${order.id}`}>{order.id}</a></td>
					<td style="text-align: center;"
						><div style="width: 210px;">
							{new Date(order.startedAt).toLocaleDateString()}
							{new Date(order.startedAt).toLocaleTimeString()}
						</div></td
					>
					<td style="text-align: center;"><StatusBadge status={order.status} /></td>
				</tr>
			{:else}
				<tr>
					<td>No Active Orders</td>
					<td />
					<td />
				</tr>
			{/each}
		</tbody>
	</table>
</section>
