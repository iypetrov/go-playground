document.addEventListener('DOMContentLoaded', async () => {
	const {publishableKey} = await fetch("/payments/config").then(r => r.json()) 
	const stripe = Stripe(publishableKey)
	const product = {
		name: "Pizza Margherita",
		description: "Classic pizza with fresh tomatoes, mozzarella, and basil.",
		price: 9.99
	};
	const {clientSecret} = await fetch("/payments/intent", { 
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify(product),
	}).then(r => r.json())

	const elements = stripe.elements ({ clientSecret })
	const paymentElement = elements.create('payment')
	paymentElement.mount('#payment-element')

	const form = document.getElementById('payment-form') 
	form.addEventListener('submit', async (e) => {
		e.preventDefault()
		const {error} = await stripe.confirmPayment({
			elements,
			confirmParams: {
				return_url: window.location.href.split('?')[0] + '/result'
			}
		})

		if (error) {
			const messages = document.getElementById('error-message')
			messages.innerText = error.message;
		}
	})
})