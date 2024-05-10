<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Action } from '$lib/types/order';
	import Logo from './logo.svelte';

	export let id: string;

	let loading = false;

	const sendAction = async (action: Action) => {
		loading = true;
		await fetch('/api/order-action', { method: 'POST', body: JSON.stringify({ id, action }) });
		setTimeout(() => {
			loading = false;
			if (action === 'cancel') {
				goto(`/orders`, { invalidateAll: true });
			} else {
				goto(`/orders/${id}`, { invalidateAll: true });
			}
		}, 1000);
	};
</script>

<div class="actions">
	{#if loading}
		<Logo loading loadingText="Processing" />
	{:else}
		<button on:click={() => sendAction('amend')}>Amend</button>
		<button on:click={() => sendAction('cancel')}>Cancel</button>
	{/if}
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
