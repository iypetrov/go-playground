function paymentForm() {
	return {
		errorMessage: '',
		async submitPayment() {
			const form = document.getElementById('payment-form');
			const elements = {};  // Replace this with the actual Stripe elements initialization if needed

			try {
				const { error } = await stripe.confirmPayment({
					elements,
					confirmParams: {
						return_url: window.location.href.split('?')[0] + 'complete.html',
					},
				});

				if (error) {
					this.errorMessage = error.message;
				}
			} catch (err) {
				this.errorMessage = 'Payment failed. Please try again.';
			}
		},
	};
}

document.addEventListener('DOMContentLoaded', async () => {
	const {publishableKey} = await fetch("/payments/config").then(r => r.json());
	const stripe = Stripe (publishableKey);

	const product = {
		name: "stripe t-shirt",
		description: "nice t-shirt",
		price: 19.99
	};
	const {clientSecret} = await fetch("/payments/intent", { 
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify(product)
	}).then(r => r.json());
	
	const elements = stripe.elements({ clientSecret });
	const paymentElement = elements.create('payment');
	paymentElement.mount('#payment-element');
})
