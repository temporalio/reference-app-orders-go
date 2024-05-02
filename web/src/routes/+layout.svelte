<script>
	import { page } from '$app/stores';
	import './app.css';

	let statusColor = '#366ee9';
</script>

<div class="app">
	<header>
		<nav>
			<a href="/"><svg width="400" height="150" viewBox="0 0 400 150">
				<text x="120" y="100">Tora</text>
			</svg></a>
			<div class="links">
				<a href="/orders" class:active={$page.url.pathname.includes('orders')}>Orders</a>
				<a href="/shipments" class:active={$page.url.pathname.includes('shipments')}>Shipments</a>
			</div>
		</nav>
		<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg" style="transform: scale(.5);">
			<path d="M 30,30 L 70,30 L 70,70 L 30,70 Z" fill={statusColor} stroke={statusColor} />
			<!-- Right arrow -->
			<g>
				<rect
					x="70"
					y="45"
					width="5"
					height="10"
					fill={statusColor}
					stroke={statusColor}
					stroke-width="1"
				/>
				<polygon
					points="75,60 90,50 75,40"
					fill={statusColor}
					stroke={statusColor}
					stroke-width="1"
				/>
			</g>
			<!-- Bottom arrow -->
			<g>
				<rect
					x="45"
					y="70"
					width="10"
					height="5"
					fill={statusColor}
					stroke={statusColor}
					stroke-width="1"
				/>
				<polygon
					points="60,75 50,90 40,75"
					fill={statusColor}
					stroke={statusColor}
					stroke-width="1"
				/>
			</g>
			<!-- Left arrow -->
			<g>
				<rect
					x="25"
					y="45"
					width="5"
					height="10"
					fill={statusColor}
					stroke={statusColor}
					stroke-width="1"
				/>
				<polygon
					points="25,60 10,50 25,40"
					fill={statusColor}
					stroke={statusColor}
					stroke-width="1"
				/>
			</g>
			<!-- Top arrow -->
			<g>
				<rect
					x="45"
					y="25"
					width="10"
					height="5"
					fill={statusColor}
					stroke={statusColor}
					stroke-width="1"
				/>
				<polygon
					points="60,25 50,10 40,25"
					fill={statusColor}
					stroke={statusColor}
					stroke-width="1"
				/>
			</g>
		</svg>
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

	nav {
		margin-bottom: 0;
		display: flex;
		justify-content: space-between;
		width: 100%;
		height: 100px;
		text-transform: uppercase;
	}
	
	nav a {
		letter-spacing: -1px;
		font-weight: 600;
	}

	text {
		font-family: "Kanit", sans-serif;
		font-size: 4.75rem;
		stroke-linejoin: round;
		text-anchor: middle;
		fill: var(--color-theme-2);
		paint-order: stroke fill;
		stroke: #000;
		stroke-width: 16px;
		letter-spacing: -14px;
		font-weight: 600;
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
