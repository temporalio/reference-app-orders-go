<script lang="ts">
	import { page } from '$app/stores';
	import { type Order } from '$lib/types/order';

	$: actionRequired = $page.data?.orders?.some(
		(o: Order) => o?.status === 'customerActionRequired'
	);
	$: statusColor = actionRequired ? '#EF8080' : '#366ee9';
</script>

<svg width="50" height="50" xmlns="http://www.w3.org/2000/svg">
	<path d="M 15,15 L 35,15 L 35,35 L 15,35 Z" fill={statusColor} stroke={statusColor} />
	<!-- Right arrow -->
	<g>
		<rect
			x="35"
			y="22.5"
			width="2.5"
			height="5"
			fill={statusColor}
			stroke={statusColor}
			stroke-width="1"
		/>
		<polygon
			points="37.5,30 45,25 37.5,20"
			fill={statusColor}
			stroke={statusColor}
			stroke-width="1"
		/>
	</g>
	<!-- Bottom arrow -->
	<g>
		<rect
			x="22.5"
			y="35"
			width="5"
			height="2.5"
			fill={statusColor}
			stroke={statusColor}
			stroke-width="1"
		/>
		<polygon
			points="30,37.5 25,45 20,37.5"
			fill={statusColor}
			stroke={statusColor}
			stroke-width="1"
		/>
	</g>
	<!-- Left arrow -->
	<g>
		<rect
			x="12.5"
			y="22.5"
			width="2.5"
			height="5"
			fill={statusColor}
			stroke={statusColor}
			stroke-width="1"
		/>
		<polygon
			points="12.5,30 5,25 12.5,20"
			fill={statusColor}
			stroke={statusColor}
			stroke-width="1"
		/>
	</g>
	<!-- Top arrow -->
	<g>
		<rect
			x="22.5"
			y="12.5"
			width="5"
			height="2.5"
			fill={statusColor}
			stroke={statusColor}
			stroke-width="1"
		/>
		<polygon
			points="30,12.5 25,5 20,12.5"
			fill={statusColor}
			stroke={statusColor}
			stroke-width="1"
		/>
	</g>
</svg>
