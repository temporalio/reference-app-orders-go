<script lang="ts">
	export let statusColor = '#366ee9';
</script>

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
		<polygon points="75,60 90,50 75,40" fill={statusColor} stroke={statusColor} stroke-width="1" />
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
		<polygon points="60,75 50,90 40,75" fill={statusColor} stroke={statusColor} stroke-width="1" />
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
		<polygon points="25,60 10,50 25,40" fill={statusColor} stroke={statusColor} stroke-width="1" />
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
		<polygon points="60,25 50,10 40,25" fill={statusColor} stroke={statusColor} stroke-width="1" />
	</g>
</svg>
