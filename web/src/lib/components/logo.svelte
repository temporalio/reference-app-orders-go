<script lang="ts">
	import { page } from '$app/stores';
	import { type Order } from '$lib/types/order';

	export let loading = false;
	export let loadingText = 'Loading';

	$: actionRequired = $page.data?.orders?.some(
		(o: Order) => o?.status === 'customerActionRequired'
	);
	$: statusColor = actionRequired ? '#EF8080' : '#366ee9';
</script>

{#if !loading}
	<svg width="400" height="150" viewBox="0 0 400 150">
		<a href="/"><text fill={statusColor} x="120" y="100">OMS</text></a>
	</svg>
{:else}
	<svg width="400" height="150" viewBox="0 0 400 150">
		<radialGradient id="gradient" cx="50%" cy="50%" r="70%">
			<animate attributeName="r" values="0%;150%;100%;0%" dur="1s" repeatCount="indefinite" />
			<stop stop-color="#366ee9" offset="0">
				<animate
					attributeName="stop-color"
					values="#eee;#366ee9;#366ee9;#eee"
					dur="1s"
					repeatCount="indefinite"
				/>
			</stop>
			<stop stop-color="#366ee9" offset="100%" />
		</radialGradient>
		<text text-anchor="middle" x="50%" y="50%" class="loading">{loadingText}</text>
	</svg>
{/if}

<style>
	a:hover {
		text-decoration: none;
	}

	text {
		font-family: 'Kanit', sans-serif;
		font-size: 3.5rem;
		stroke-linejoin: round;
		text-anchor: middle;
		paint-order: stroke fill;
		stroke: #000;
		stroke-width: 16px;
		letter-spacing: -10px;
		font-weight: 600;
		text-transform: uppercase;
	}

	.loading {
		fill: url(#gradient);
	}
</style>
