import { writable } from "svelte/store";

export type Order = { id: string, name: string, description: string}
export const order = writable<Order | undefined>();

export const items: Order[] = [
  {
    "id": "1",
    "name": "Nike Air Force Ones",
    "description": "The Nike Air Force Ones combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "id": "2",
    "name": "Adidas UltraBoost",
    "description": "The Adidas UltraBoost combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "id": "3",
    "name": "Reebok Classic Leather",
    "description": "The Reebok Classic Leather combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "id": "4",
    "name": "Puma Suede Classic",
    "description": "The Puma Suede Classic combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "id": "5",
    "name": "New Balance 574",
    "description": "The New Balance 574 combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "id": "6",
    "name": "Vans Old Skool",
    "description": "The Vans Old Skool combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "id": "7",
    "name": "Converse Chuck Taylor All Star",
    "description": "The Converse Chuck Taylor All Star combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "id": "8",
    "name": "Under Armour HOVR Sonic",
    "description": "The Under Armour HOVR Sonic combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "id": "9",
    "name": "Jordan Air Jordan 1",
    "description": "The Jordan Air Jordan 1 combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "id": "10",
    "name": "Asics GEL-Kayano",
    "description": "The Asics GEL-Kayano combines timeless style with modern comfort, featuring premium materials and cutting-edge technology for unmatched performance."
  },
  {
    "id": "11",
    "name": "Nike Air Force Ones",
    "description": "A second iteration of the classic, the Nike Air Force Ones Model 11 is redesigned for the modern athlete, offering enhanced cushioning and durability."
  },
  {
    "id": "12",
    "name": "Adidas UltraBoost",
    "description": "Adidas UltraBoost Model 12 brings you closer to the ground for a more responsive feel, ideal for runners seeking a blend of support and speed."
  },
  {
    "id": "13",
    "name": "Reebok Classic Leather",
    "description": "Reebok's Classic Leather Model 13 is the epitome of retro chic, offering unparalleled comfort and a sleek design for everyday wear."
  },
  {
    "id": "14",
    "name": "Puma Suede Classic",
    "description": "The Puma Suede Classic Model 14 updates the iconic design with advanced materials for better wearability and style."
  },
  {
    "id": "15",
    "name": "New Balance 574",
    "description": "This latest version of the New Balance 574, Model 15, combines heritage styling with modern technology for an improved fit and function."
  },
  {
    "id": "16",
    "name": "Vans Old Skool",
    "description": "Vans Old Skool Model 16 reintroduces the classic skate shoe with updated features for enhanced performance and comfort."
  },
  {
    "id": "17",
    "name": "Converse Chuck Taylor All Star",
    "description": "Model 17 of the Converse Chuck Taylor All Star elevates the iconic silhouette with premium materials and an improved insole for all-day comfort."
  },
  {
    "id": "18",
    "name": "Under Armour HOVR Sonic",
    "description": "The latest Under Armour HOVR Sonic, Model 18, offers an unparalleled ride, blending the perfect balance of cushioning and energy return."
  },
  {
    "id": "19",
    "name": "Jordan Air Jordan 1",
    "description": "Jordan Air Jordan 1 Model 19 continues the legacy, integrating classic design lines with modern technology for a timeless look and feel."
  },
  {
    "id": "20",
    "name": "Asics GEL-Kayano",
    "description": "Asics GEL-Kayano Model 20 is the latest in the series, offering improved stability and support for overpronators."
  }
]
