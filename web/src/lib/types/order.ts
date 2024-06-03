import { writable } from 'svelte/store';

export interface Item {
	sku: string;
	description: string;
}
export interface OrderItem extends Item {
	sku: string;
	quantity: number;
}

export type OrderStatus =
	| 'pending'
	| 'processing'
	| 'customerActionRequired'
	| 'completed'
	| 'cancelled';

export interface Order {
	id: string;
	customerId: string;
	items: OrderItem[];
	fulfillments?: Fulfillment[];
	status?: OrderStatus;
}

export type FulfillmentStatus =
	| 'unavailable'
	| 'pending'
	| 'processing'
	| 'completed'
	| 'cancelled'
	| 'failed';

export interface Fulfillment {
	id: string;
	shipment?: Shipment;
	items: OrderItem[];
	payment?: Payment;
	location: string;
	status?: FulfillmentStatus;
}

export interface Shipment {
	id: string;
	status: string;
	items: OrderItem[];
	updatedAt: string;
}

export type PaymentStatus = 'pending' | 'success' | 'failed';

export interface Payment {
	shipping: number;
	tax: number;
	subTotal: number;
	total: number;
	status: PaymentStatus;
}

export type Action = 'amend' | 'cancel';

export const order = writable<Order | undefined>();

export const generateOrders = (quantity: number): Order[] => {
	const orders = [];
	for (let i = 0; i < quantity; i++) {
		const shuffledItems = items.sort(() => 0.5 - Math.random());
		const n = Math.floor(Math.random() * 5) + 2;
		const selected = shuffledItems.slice(0, n);
		orders.push({
			id: `A${i + 1}-${Date.now()}`,
			customerId: '1234',
			items: selected.map((item) => ({ ...item, quantity: Math.floor(Math.random() * 3) + 1 }))
		});
	}
	// For demo purposes, we'll ensure that the first order always
	// contains a single specific item (whose SKU does not trigger
	// item unavailability). This will make it easy to demonstrate
	// the simplest possible case.
	const sortedItems = items.sort(function (a, b) {
		return a.sku < b.sku ? -1 : a.sku > b.sku ? 1 : 0;
	});
	orders[0].items = [
		{
			sku: sortedItems[sortedItems.length - 1].sku,
			description: items[sortedItems.length - 1].description,
			quantity: 1
		}
	];
	return orders;
};

export const items: Item[] = [
	{
		sku: 'Nike Air Force Ones',
		description:
			'The Nike Air Force Ones combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance.'
	},
	{
		sku: 'Adidas UltraBoost',
		description:
			'The Adidas UltraBoost combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance.'
	},
	{
		sku: 'Reebok Classic Leather White',
		description:
			'The Reebok Classic Leather combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance.'
	},
	{
		sku: 'Puma Suede Classic',
		description:
			'The Puma Suede Classic combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance.'
	},
	{
		sku: 'New Balance 574',
		description:
			'The New Balance 574 combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance.'
	},
	{
		sku: 'Vans Old Skool',
		description:
			'The Vans Old Skool combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance.'
	},
	{
		sku: 'Converse Chuck Taylor All Star',
		description:
			'The Converse Chuck Taylor All Star combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance.'
	},
	{
		sku: 'Under Armour HOVR Sonic',
		description:
			'The Under Armour HOVR Sonic combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance.'
	},
	{
		sku: 'Jordan Air Jordan 1',
		description:
			'The Jordan Air Jordan 1 combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance.'
	},
	{
		sku: 'Asics GEL-Kayano',
		description:
			'The Asics GEL-Kayano combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance.'
	},
	{
		sku: 'Nike Air Force Ones',
		description:
			'A second iteration of the classic, the Nike Air Force Ones Model 11 is redesigned for the modern athlete, offering enhanced cushioning and durability.'
	},
	{
		sku: 'Adidas UltraBoost',
		description:
			'Adidas UltraBoost Model 12 brings you closer to the ground for a more responsive feel, ideal for runners seeking a blend of support and speed.'
	},
	{
		sku: 'Reebok Classic Leather Black',
		description:
			"Reebok's Classic Leather Model 13 is the epitome of retro chic, offering unparalleled comfort and a sleek design for everyday wear."
	},
	{
		sku: 'Puma Suede Classic',
		description:
			'The Puma Suede Classic Model 14 updates the iconic design with advanced materials for better wearability and style.'
	},
	{
		sku: 'New Balance 574',
		description:
			'This latest version of the New Balance 574, Model 15, combines heritage styling with modern technology for an improved fit and function.'
	},
	{
		sku: 'Vans New Skool',
		description:
			'Vans New Skool Model 16 reintroduces the classic skate shoe with updated features for enhanced performance and comfort.'
	},
	{
		sku: 'Converse Chuck Taylor All Star',
		description:
			'Model 17 of the Converse Chuck Taylor All Star elevates the iconic silhouette with premium materials and an improved insole for all-day comfort.'
	},
	{
		sku: 'Under Armour HOVR Sonic',
		description:
			'The latest Under Armour HOVR Sonic, Model 18, offers an unparalleled ride, blending the perfect balance of cushioning and energy return.'
	},
	{
		sku: 'Jordan Air Jordan 2',
		description:
			'Jordan Air Jordan 2 Model 19 continues the legacy, integrating classic design lines with modern technology for a timeless look and feel.'
	},
	{
		sku: 'Asics GEL-Tres',
		description:
			'Asics GEL-Tres Model 20 is the latest in the series, offering improved stability and support for overpronators.'
	}
];
