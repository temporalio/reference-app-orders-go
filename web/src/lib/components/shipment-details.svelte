<script lang="ts">
	import { type Shipment } from '$lib/types/order';
	import ItemDetails from './item-details.svelte';
	import ShipmentStatus from './shipment-status.svelte';

	export let shipment: Shipment;
	export let status = '';
</script>

<div class="details">
	<ShipmentStatus {status} />
	<ItemDetails items={shipment.items} />
</div>

<style>
	.details {
		display: flex;
		flex-direction: column;
		gap: 2rem;
		background-color: white;
		align-items: center;
		width: 100%;
	}
</style>
