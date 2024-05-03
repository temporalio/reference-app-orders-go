<script lang="ts">
	import { type Shipment } from '$lib/stores/order';

	export let shipment: Shipment | undefined;
	export let status = '';

	const activeStatuses = ['booked', 'dispatched', 'delivered'];
	const noStatus = !shipment || !activeStatuses.includes(status || shipment?.status);
	const statuses = noStatus ? ['pending'] : activeStatuses;
	$: currentIndex = noStatus ? 0 : statuses.indexOf(status || shipment.status);
</script>

<ul>
	{#each statuses as s, index}
		<li
			class:active={currentIndex === index}
			class:completed={currentIndex > index}
			class:incomplete={currentIndex < index}
		>
			{s.toUpperCase()}
		</li>
	{/each}
</ul>

<style>
	ul {
		position: relative;
		list-style: none;
		display: inline-flex;
		border: 3px solid black;
		border-radius: 9999px;
		overflow: hidden;
	}

	li {
		padding: 0.75em 1.5em;
		position: relative;
		background: transparent;
		z-index: 1;
		font-weight: 700;
	}

	li::before {
		content: '';
		position: absolute;
		inset: 0;
		border-left: 3px solid black;
		transform: skew(30deg);
		z-index: -1;
	}

	li.completed::before {
		background: forestgreen;
	}

	li.active::before {
		background: lightgreen;
		animation-name: color;
		animation-duration: 2s;
		animation-iteration-count: infinite;
		animation-direction: alternate-reverse;
		animation-timing-function: ease;
	}

	@keyframes color {
		from {
			background-color: forestgreen;
		}
		to {
			background-color: lightgreen;
		}
	}

	li.incomplete::before {
		background: lightcoral;
	}

	li:first-child {
		/* extend the first item leftward to fill the rest of the space */
		margin-left: -4rem;
		padding-left: 4rem;
	}

	li:last-child {
		/* extend the last item rightward to fill the rest of the space */
		margin-right: -2rem;
		padding-right: 4rem;
	}

	@media (max-width: 640px) {
		li {
			font-size: 0.75rem;
			padding: 0.5em 1em;
		}

		li:first-child {
			/* extend the first item leftward to fill the rest of the space */
			margin-left: -3rem;
			padding-left: 2rem;
		}

		li:last-child {
			/* extend the last item rightward to fill the rest of the space */
			margin-right: -2rem;
			padding-right: 3rem;
		}
	}
</style>
