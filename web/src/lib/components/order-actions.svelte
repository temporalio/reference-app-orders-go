<script lang="ts">
	import type { Action } from '$lib/types/order';

	export let id: string;

	const sendAction = async (action: Action) => {
		await fetch('/api/order-action', { method: 'POST', body: JSON.stringify({ id, action }) });
	};
</script>

<div class="actions">
	<button on:click={() => sendAction('amend')}>Amend</button>
	<button on:click={() => sendAction('cancel')}>Cancel</button>
</div>

<style>
	.actions {
		display: flex;
		align-items: center;
		gap: 1rem;
		justify-content: end;
		width: 100%;
	}

	@media (max-width: 640px) {
	}
</style>
