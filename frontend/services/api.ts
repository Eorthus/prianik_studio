// services/api.ts
import { ref } from "vue";
import { useRuntimeConfig } from "nuxt/app";

import type {
  APIResponse,
  Category,
  Product,
  GalleryItem,
  ContactFormData,
  OrderData,
} from "../components";

// API сервис
export const useApiService = () => {
  // Базовый URL API из переменных окружения
  const API_BASE_URL = useRuntimeConfig().public.apiBaseUrl;
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // Функция для выполнения HTTP запросов
  const fetchApi = async <T>(
    endpoint: string,
    options: RequestInit = {},
    language: string = "ru"
  ): Promise<APIResponse<T>> => {
    isLoading.value = true;
    error.value = null;

    try {
      // Добавляем параметр языка к эндпоинтам, которые его поддерживают
      const needsLanguage =
        !endpoint.includes("orders") && !endpoint.includes("contact");
      const separator = endpoint.includes("?") ? "&" : "?";
      const finalEndpoint = needsLanguage
        ? `${endpoint}${separator}language=${language}`
        : endpoint;

      const response = await fetch(`${API_BASE_URL}${finalEndpoint}`, {
        ...options,
        headers: {
          "Content-Type": "application/json",
          ...options.headers,
        },
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(
          data.error || "Произошла ошибка при выполнении запроса"
        );
      }

      return data;
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "Неизвестная ошибка";
      error.value = errorMessage;
      console.error("API error:", errorMessage);

      return {
        success: false,
        data: null as any,
        error: errorMessage,
      };
    } finally {
      isLoading.value = false;
    }
  };

  // Получение списка категорий
  const getCategories = (language: string = "ru") => {
    return fetchApi<Category[]>("/categories", {}, language);
  };

  // Получение списка товаров с пагинацией и фильтрацией
  const getProducts = (
    page = 1,
    pageSize = 10,
    categoryId?: number,
    subcategoryId?: number,
    search?: string,
    sortPrice?: "asc" | "desc",
    language: string = "ru"
  ) => {
    let endpoint = `/products?page=${page}&page_size=${pageSize}`;

    if (categoryId) endpoint += `&category=${categoryId}`;
    if (subcategoryId) endpoint += `&subcategory=${subcategoryId}`;
    if (search) endpoint += `&search=${encodeURIComponent(search)}`;
    if (sortPrice) endpoint += `&sort_price=${sortPrice}`;

    return fetchApi<{ items: Product[]; total_items: number }>(
      endpoint,
      {},
      language
    );
  };

  // Получение товара по ID
  const getProductById = (id: number, language: string = "ru") => {
    return fetchApi<Product>(`/products/${id}`, {}, language);
  };

  // Получение связанных товаров
  const getRelatedProducts = (
    id: number,
    limit = 5,
    language: string = "ru"
  ) => {
    return fetchApi<Product[]>(
      `/products/${id}/related?limit=${limit}`,
      {},
      language
    );
  };

  // Получение элементов галереи
  const getGalleryItems = (categoryId?: number, language: string = "ru") => {
    let endpoint = "/gallery";
    if (categoryId) endpoint += `?category=${categoryId}`;

    return fetchApi<GalleryItem[]>(endpoint, {}, language);
  };

  // Обработчик возможных ошибок API
  const handleApiError = (error: any): string => {
    if (typeof error === "string") return error;
    if (error?.message) return error.message;
    return "Произошла неизвестная ошибка";
  };

  // Отправка формы обратной связи
  const submitContactForm = (formData: ContactFormData) => {
    return fetchApi<{ message: string }>("/contact", {
      method: "POST",
      body: JSON.stringify(formData),
    });
  };

  // Создание заказа
  const createOrder = (orderData: OrderData) => {
    return fetchApi<{ order_id: number; message: string }>("/orders", {
      method: "POST",
      body: JSON.stringify(orderData),
    });
  };

  return {
    isLoading,
    error,
    getCategories,
    getProducts,
    getProductById,
    getRelatedProducts,
    getGalleryItems,
    submitContactForm,
    createOrder,
    handleApiError,
  };
};
