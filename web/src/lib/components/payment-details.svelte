<script lang="ts">
	import type { Payment } from '$lib/stores/order';

	export let payment: Payment;
</script>

<div>
	<div class="row">
		<p class="label">Payment</p>
		<p><strong>{payment.status.toUpperCase()}</strong></p>
	</div>
	<div class="row">
		<p class="label">Subtotal</p>
		<p>{payment.subTotal.toLocaleString('en-US', { style: 'currency', currency: 'USD' })}</p>
	</div>
	<div class="row">
		<p class="label">Tax</p>
		<p>{payment.tax.toLocaleString('en-US', { style: 'currency', currency: 'USD' })}</p>
	</div>
	<div class="row">
		<p class="label">Shipping</p>
		<p>{payment.shipping.toLocaleString('en-US', { style: 'currency', currency: 'USD' })}</p>
	</div>
	<div class="row">
		<p class="label">Total</p>
		<p>{payment.total.toLocaleString('en-US', { style: 'currency', currency: 'USD' })}</p>
	</div>
</div>

<style>
	.row {
		display: flex;
	}

	.label {
		width: 200px;
		font-weight: 700;
	}
</style>
