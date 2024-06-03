<script lang="ts">
	import Logo from '$lib/components/logo.svelte';
	import './app.css';
</script>

<svelte:head>
	<title>OMS</title>
	<meta name="description" content="OMS App" />
</svelte:head>

<div class="app">
	<header>
		<nav>
			<Logo />
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
		position: relative;
	}

	header {
		padding: 1rem 2rem;
		display: flex;
		justify-content: end;
	}

	nav {
		margin-bottom: 0;
		display: flex;
		justify-content: space-between;
		width: 100%;
		height: 100px;
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
