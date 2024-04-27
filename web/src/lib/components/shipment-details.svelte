<script lang="ts">
	import { type Shipment } from '$lib/stores/order';
	import ItemDetails from './item-details.svelte';
	import ShipmentStatus from './shipment-status.svelte';

	export let shipment: Shipment;
	export let status = '';
</script>

<div class="details">
	<ShipmentStatus {shipment} {status} />
	<ItemDetails items={shipment.items} />
</div>

<style>
	.details {
		display: flex;
		width: 50%;
		flex-direction: column;
		gap: 2rem;
		background-color: white;
		padding: 2rem;
		align-items: start;
	}
</style>
