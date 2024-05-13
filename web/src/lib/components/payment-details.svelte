<script lang="ts">
	import type { Payment } from '$lib/types/order';
	import StatusBadge from './status-badge.svelte';

	export let payment: Payment;
</script>

<div class="container">
	<div class="row">
		<p class="label">Payment</p>
		<StatusBadge status={payment.status} />
	</div>
	<div class="row">
		<p class="label">Subtotal</p>
		<p class="value">
			{(payment.subTotal / 100).toLocaleString('en-US', { style: 'currency', currency: 'USD' })}
		</p>
	</div>
	<div class="row">
		<p class="label">Tax</p>
		<p class="value">
			{(payment.tax / 100).toLocaleString('en-US', { style: 'currency', currency: 'USD' })}
		</p>
	</div>
	<div class="row" style="border-bottom: 2px solid black;">
		<p class="label">Shipping</p>
		<p class="value">
			{(payment.shipping / 100).toLocaleString('en-US', { style: 'currency', currency: 'USD' })}
		</p>
	</div>
	<div class="row">
		<p class="label">Total</p>
		<p class="value">
			{(payment.total / 100).toLocaleString('en-US', { style: 'currency', currency: 'USD' })}
		</p>
	</div>
</div>

<style>
	.container {
		margin-top: 1rem;
		padding: 0.55rem;
		background-color: var(--color-bg-2);
		border-radius: 0.15rem;
	}
	.row {
		display: flex;
	}

	.label {
		width: 240px;
		font-weight: 700;
	}

	.value {
		width: 100%;
		text-align: right;
		padding-right: 0.15rem;
	}
</style>
