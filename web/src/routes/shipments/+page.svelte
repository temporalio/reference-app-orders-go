<script lang="ts">
	export let data;

	$: ({ shipments } = data);
</script>

<svelte:head>
	<title>Tora</title>
	<meta name="description" content="Tora App" />
</svelte:head>

<section>
	<nav>
		<h1>Shipments</h1>
	</nav>
	<table>
		<thead>
			<tr>
				<th>ID</th>
				<th>Status</th>
			</tr>
		</thead>
		<tbody>
			{#each shipments as shipment}
				<tr>
					<td style="width: 100%;"><a href={`/shipments/${shipment.id}`}>{shipment.id}</a></td>

					<td><div class="status">{shipment.status.toUpperCase()}</div></td>
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
