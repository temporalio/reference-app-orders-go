<script>
	import './app.css';

	let statusColor = '#788F9C';
</script>

<div class="app">
	<header>
		<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg">
			<path d="M 30,30 L 70,30 L 70,70 L 30,70 Z" fill={statusColor} stroke={statusColor}/>
			<!-- Right arrow -->
			<g>
				<rect x="70" y="45" width="10" height="10" fill={statusColor} stroke={statusColor} stroke-width="1" />
				<polygon points="80,60 90,50 80,40" fill={statusColor} stroke={statusColor} stroke-width="1" />	
			</g>
			  <!-- Bottom arrow -->
			<g>
				<rect x="45" y="70" width="10" height="10" fill={statusColor} stroke={statusColor} stroke-width="1" />
				<polygon points="60,80 50,90 40,80" fill={statusColor} stroke={statusColor} stroke-width="1" />	
			</g>
			  <!-- Left arrow -->
			<g>
				<rect x="20" y="45" width="10" height="10" fill={statusColor} stroke={statusColor} stroke-width="1" />
				<polygon points="20,60 10,50 20,40" fill={statusColor} stroke={statusColor} stroke-width="1" />	
			</g>
			  <!-- Top arrow -->
			<g>
				<rect x="45" y="20" width="10" height="10" fill={statusColor} stroke={statusColor} stroke-width="1" />
				<polygon points="60,20 50,10 40,20" fill={statusColor} stroke={statusColor} stroke-width="1" />	
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
