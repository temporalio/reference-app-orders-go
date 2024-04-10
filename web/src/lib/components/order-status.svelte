<script lang="ts">
	import { type Order } from '$lib/stores/order';

	export let order: Order | undefined;
</script>

<div class="status">
	<h3>Status</h3>
	<div class="status-bar">
		<div class="bar green">Order entered</div>
	</div>
</div>

<style>
	.status {
		display: flex;
		gap: 1rem;
		align-items: center;
		padding: 0 1rem;
	}

	.status-bar {
		display: flex;
		gap: 1rem;
		align-items: center;
		width: 100%;
	}

	.bar {
		width: 100%;
		height: 40px;
		display: flex;
		align-items: center;
		justify-content: center;
		border: 2px black solid;
		border-radius: 4px;
	}

	.green {
		background: green;
		color: white;
	}
</style>
