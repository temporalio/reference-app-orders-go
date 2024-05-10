<script lang="ts">
	export let id: string | undefined = undefined;
	export let status = '';

	if (id) {
		const broadcaster = new BroadcastChannel(`shipment-${id}`);
		broadcaster?.addEventListener('message', (event) => {
			status = event.data;
		});
	}

	const inactiveStatuses = ['pending', 'unavailable'];
	const activeStatuses = ['booked', 'dispatched', 'delivered'];

	$: finalStatus = status || 'pending';
	$: inactive = inactiveStatuses.includes(finalStatus);
	$: statuses = inactive ? [finalStatus] : activeStatuses;
	$: currentIndex = inactive ? 0 : activeStatuses.indexOf(finalStatus);
</script>

<ul>
	{#each statuses as s, index}
		<li
			class:active={!inactive && currentIndex === index}
			class:completed={!inactive && currentIndex > index}
			class:incomplete={!inactive && currentIndex < index}
			class:pending={s === 'pending'}
			class:unavailable={s === 'unavailable'}
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
	li.pending::before {
		background: lightgoldenrodyellow;
	}

	li.incomplete::before,
	li.unavailable::before {
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
