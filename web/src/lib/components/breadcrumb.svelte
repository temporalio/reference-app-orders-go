<script lang="ts">
	export let label: string;
	export let href: string;
</script>

<a {href}>{label}</a>
<slot />

<style>
	a {
		position: absolute;
		top: 4rem;
		right: 2rem;
		letter-spacing: -1px;
		font-weight: 600;
		text-transform: uppercase;
	}
</style>
