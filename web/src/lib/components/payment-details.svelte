<script lang="ts">
	import type { Payment } from '$lib/stores/order';

	export let payment: Payment;
</script>

<div class="container">
	<div class="row">
		<p class="label">Payment</p>
		<p><strong>{payment.status.toUpperCase()}</strong></p>
	</div>
	<div class="row">
		<p class="label">Subtotal</p>
		<p class="value">
			{payment.subTotal.toLocaleString('en-US', { style: 'currency', currency: 'USD' })}
		</p>
	</div>
	<div class="row">
		<p class="label">Tax</p>
		<p class="value">
			{payment.tax.toLocaleString('en-US', { style: 'currency', currency: 'USD' })}
		</p>
	</div>
	<div class="row">
		<p class="label">Shipping</p>
		<p class="value">
			{payment.shipping.toLocaleString('en-US', { style: 'currency', currency: 'USD' })}
		</p>
	</div>
	<div class="row">
		<p class="label">Total</p>
		<p class="value">
			{payment.total.toLocaleString('en-US', { style: 'currency', currency: 'USD' })}
		</p>
	</div>
</div>

<style>
	.container {
		margin-top: 1rem;
		font-size: 0.9rem;
		padding: 0.55rem;
		background-color: var(--color-bg-2);
		border-radius: 0.15rem;
	}
	.row {
		display: flex;
	}

	.label {
		width: 200px;
		font-weight: 700;
	}

	.value {
		width: 100%;
		text-align: right;
	}
</style>
