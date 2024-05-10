<script lang="ts">
	export let status: string;

	const capitalize = (word?: string): string => {
		if (!word) return '';
		return word[0].toUpperCase() + word.slice(1);
	};

	const spaceBetweenCapitalLetters = (word: string): string => {
		return capitalize(word)
			.replace(/([A-Z])/g, ' $1')
			.trim()
			.toUpperCase();
	};

	$: formattedStatus =
		status === 'customerActionRequired'
			? spaceBetweenCapitalLetters('actionRequired')
			: spaceBetweenCapitalLetters(status);
</script>

<div class="status {status}">
	{formattedStatus}
</div>

<style>
	.status {
		border: 2px solid black;
		border-radius: 9999px;
		padding: 0.25rem 1rem;
		margin: 0 auto;
		background-color: lightgreen;
		font-size: 14px;
		white-space: nowrap;
		font-weight: 700;
	}

	.customerActionRequired {
		background-color: lightcoral;
	}
</style>
