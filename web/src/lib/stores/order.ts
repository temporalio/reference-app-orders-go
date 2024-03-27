import { writable } from "svelte/store";

export interface Item { sku: string, name: string, description: string };
export interface OrderItem extends Item { quantity: number };
export interface Order { id: string, items: OrderItem[] };

export const order = writable<Order | undefined>();

export const generateOrders = (quantity: number): Order[] => {
  const orders = [];
  for (let i = 0; i < quantity; i++) {
    const shuffledItems = items.sort(() => 0.5 - Math.random());
    const n = Math.floor(Math.random() * 10) + 1
    const selected = shuffledItems.slice(0, n);
    orders.push({
      id: `order-${i + 1}`,
      items: selected.map((item) => ({ ...item, quantity: Math.floor(Math.random() * 5) + 1 }))
    });
  }
  return orders;
}

export const items: Item[] = [
  {
    "sku": "nike-air-force-ones",
    "name": "Nike Air Force Ones",
    "description": "The Nike Air Force Ones combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "sku": "adidas-ultraboost",
    "name": "Adidas UltraBoost",
    "description": "The Adidas UltraBoost combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "sku": "reebok-classic-leather-white",
    "name": "Reebok Classic Leather White",
    "description": "The Reebok Classic Leather combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "sku": "puma-suede-classic",
    "name": "Puma Suede Classic",
    "description": "The Puma Suede Classic combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "sku": "new-balance-574",
    "name": "New Balance 574",
    "description": "The New Balance 574 combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "sku": "vans-old-skool",
    "name": "Vans Old Skool",
    "description": "The Vans Old Skool combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "sku": "converse-chuck-taylor-all-star",
    "name": "Converse Chuck Taylor All Star",
    "description": "The Converse Chuck Taylor All Star combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "sku": "under-armour-hovr-sonic",
    "name": "Under Armour HOVR Sonic",
    "description": "The Under Armour HOVR Sonic combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "sku": "jordan-air-jordan-1",
    "name": "Jordan Air Jordan 1",
    "description": "The Jordan Air Jordan 1 combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "sku": "asics-gel-kayano",
    "name": "Asics GEL-Kayano",
    "description": "The Asics GEL-Kayano combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "sku": "nike-air-force-ones",
    "name": "Nike Air Force Ones",
    "description": "A second iteration of the classic, the Nike Air Force Ones Model 11 is redesigned for the modern athlete, offering enhanced cushioning and durability."
  },
  {
    "sku": "adidas-ultraboost",
    "name": "Adidas UltraBoost",
    "description": "Adidas UltraBoost Model 12 brings you closer to the ground for a more responsive feel, ideal for runners seeking a blend of support and speed."
  },
  {
    "sku": "reebok-classic-leather-black",
    "name": "Reebok Classic Leather Black",
    "description": "Reebok's Classic Leather Model 13 is the epitome of retro chic, offering unparalleled comfort and a sleek design for everyday wear."
  },
  {
    "sku": "puma-suede-classic",
    "name": "Puma Suede Classic",
    "description": "The Puma Suede Classic Model 14 updates the iconic design with advanced materials for better wearability and style."
  },
  {
    "sku": "new-balance-574",
    "name": "New Balance 574",
    "description": "This latest version of the New Balance 574, Model 15, combines heritage styling with modern technology for an improved fit and function."
  },
  {
    "sku": "vans-new-skool",
    "name": "Vans New Skool",
    "description": "Vans New Skool Model 16 reintroduces the classic skate shoe with updated features for enhanced performance and comfort."
  },
  {
    "sku": "converse-chuck-taylor-all-star",
    "name": "Converse Chuck Taylor All Star",
    "description": "Model 17 of the Converse Chuck Taylor All Star elevates the iconic silhouette with premium materials and an improved insole for all-day comfort."
  },
  {
    "sku": "under-armour-hovr-sonic",
    "name": "Under Armour HOVR Sonic",
    "description": "The latest Under Armour HOVR Sonic, Model 18, offers an unparalleled ride, blending the perfect balance of cushioning and energy return."
  },
  {
    "sku": "jordan-air-jordan-2",
    "name": "Jordan Air Jordan 2",
    "description": "Jordan Air Jordan 2 Model 19 continues the legacy, integrating classic design lines with modern technology for a timeless look and feel."
  },
  {
    "sku": "asics-gel-tres",
    "name": "Asics GEL-Tres",
    "description": "Asics GEL-Tres Model 20 is the latest in the series, offering improved stability and support for overpronators."
  }
]
