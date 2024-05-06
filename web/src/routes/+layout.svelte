<script lang="ts">
	import { page } from '$app/stores';
	import type { LayoutData } from './$types';

	import Logo from '$lib/components/logo.svelte';
	import StatusIcon from '$lib/components/status-icon.svelte';
	import type { Order } from '$lib/types/order';
	import './app.css';

	export let data: LayoutData;

	$: actionRequired = data?.orders.some((o: Order) => o?.status === 'customerActionRequired');
</script>

<svelte:head>
	<title>Tora</title>
	<meta name="description" content="Tora App" />
</svelte:head>

<div class="app">
	<header>
		<nav>
			<Logo />
			<div class="action">
				<div class="links">
					<a href="/orders" class:active={$page.url.pathname.includes('orders')}>Orders</a>
					<a href="/shipments" class:active={$page.url.pathname.includes('shipments')}>Shipments</a>
				</div>
				<StatusIcon {actionRequired} />
			</div>
		</nav>
	</header>
	<main>
		<slot />
	</main>
</div>

<style>
	.app {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	header {
		padding: 1rem 2rem;
		display: flex;
		justify-content: end;
	}

	.action {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	nav {
		margin-bottom: 0;
		display: flex;
		justify-content: space-between;
		width: 100%;
		height: 100px;
	}

	nav a {
		letter-spacing: -1px;
		font-weight: 600;
		text-transform: uppercase;
	}

	.links {
		display: flex;
		gap: 1rem;
	}

	.active {
		text-decoration: underline;
	}

	main {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 1rem;
		width: 100%;
		max-width: 64rem;
		margin: 0 auto;
		box-sizing: border-box;
	}
</style>
