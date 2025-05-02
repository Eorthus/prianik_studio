export type LangType = "en" | "es" | "ru";

export interface Product {
  id: number;
  name: string;
  price: number;
  description: string;
  currency: "USD" | "RUB" | "CRC";
  category_id: number;
  subcategory_id?: number;
  characteristics?: Record<string, string>;
  images: string[];
  related_products: Product[];
  translations?: Record<string, { name: string; description: string }>;
}

export interface Category {
  id: number;
  name: string;
  subcategories: Subcategory[];
}

export interface Subcategory {
  id: number;
  name: string;
}

export interface GalleryItem {
  id: number;
  title: string;
  description: string;
  full: string;
  thumbnail: string;
  category: string;
  category_id: number;
}

export interface ContactFormData {
  name: string;
  email: string;
  phone: string;
  language: LangType
  message: string;
  recaptchaResponse: string;
}

export interface OrderItem {
  product_id: number;
  quantity: number;
}

export interface OrderData {
  name: string;
  email: string;
  phone: string;
  comment?: string;
  items: OrderItem[];
  language: LangType
  recaptchaResponse: string;
}

export interface APIResponse<T> {
  success: boolean;
  data: T;
  error?: string;
  validation_errors?: { field: string; message: string }[];
}
