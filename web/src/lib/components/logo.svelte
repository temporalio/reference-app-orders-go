<script lang="ts">
	export let loading = false;
</script>

{#if !loading}
	<svg width="400" height="150" viewBox="0 0 400 150">
		<a href="/"><text x="120" y="100">Tora</text></a>
	</svg>
{:else}
	<svg width="400" height="400" viewBox="0 0 400 400">
		<radialGradient id="gradient" cx="50%" cy="50%" r="70%">
			<animate attributeName="r" values="0%;150%;100%;0%" dur="1s" repeatCount="indefinite" />
			<stop stop-color="#366ee9" offset="0">
				<animate
					attributeName="stop-color"
					values="#333;#366ee9;#366ee9;#333"
					dur="1s"
					repeatCount="indefinite"
				/>
			</stop>
			<stop stop-color="rgba(55,55,55,0)" offset="100%" />
		</radialGradient>
		<text text-anchor="middle" x="50%" y="50%" class="loading">Tora </text>
	</svg>
{/if}

<style>
	a:hover {
		text-decoration: none;
	}

	text {
		font-family: 'Kanit', sans-serif;
		font-size: 4.75rem;
		stroke-linejoin: round;
		text-anchor: middle;
		fill: var(--color-theme-2);
		paint-order: stroke fill;
		stroke: #000;
		stroke-width: 16px;
		letter-spacing: -14px;
		font-weight: 600;
		text-transform: uppercase;
	}

	.loading {
		fill: url(#gradient);
	}
</style>
