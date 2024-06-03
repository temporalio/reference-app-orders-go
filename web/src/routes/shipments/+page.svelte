<script lang="ts">
	import StatusBadge from '$lib/components/status-badge.svelte';

	export let data;

	$: ({ shipments } = data);
</script>

<svelte:head>
	<title>OMS</title>
	<meta name="description" content="OMS App" />
</svelte:head>

<section>
	<nav>
		<h1>Shipments</h1>
	</nav>
	<table>
		<thead>
			<tr>
				<th>ID</th>
				<th style="text-align: center;">Status</th>
			</tr>
		</thead>
		<tbody>
			{#each shipments as shipment}
				<tr>
					<td style="width: 100%;"><a href={`/shipments/${shipment.id}`}>{shipment.id}</a></td>

					<td style="text-align: center;"><StatusBadge status={shipment.status} /></td>
				</tr>
			{:else}
				<tr>
					<td>No Active Shipments</td>
					<td />
				</tr>
			{/each}
		</tbody>
	</table>
</section>
