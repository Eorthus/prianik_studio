import { ref, watch } from "vue";

// Типы данных для корзины
export interface CartItem {
  id: number;
  name: string;
  price: number;
  image?: string;
  quantity: number;
  currency?: string;
}

// Используем локальное хранилище для сохранения корзины между сеансами
const CART_STORAGE_KEY = 'prianik-cart';

// Начальное состояние корзины
const cart = ref<CartItem[]>([]);
const cartCounter = ref(0);

// Загрузка корзины из localStorage при инициализации
const loadCart = () => {
  if (process.client) {
    try {
      const savedCart = localStorage.getItem(CART_STORAGE_KEY);
      if (savedCart) {
        cart.value = JSON.parse(savedCart);
        updateCartCounter();
      }
    } catch (error) {
      console.error('Ошибка при загрузке корзины:', error);
    }
  }
};

// Сохранение корзины в localStorage
const saveCart = () => {
  if (process.client) {
    try {
      localStorage.setItem(CART_STORAGE_KEY, JSON.stringify(cart.value));
    } catch (error) {
      console.error('Ошибка при сохранении корзины:', error);
    }
  }
};

// Обновление счетчика корзины
const updateCartCounter = () => {
  cartCounter.value = cart.value.reduce((total, item) => total + item.quantity, 0);
};

// Добавление товара в корзину
const addProductToCart = (product: Partial<CartItem> & { id: number, name: string, price: number }, quantity: number = 1) => {
  const existingItemIndex = cart.value.findIndex(item => item.id === product.id);
  
  if (existingItemIndex >= 0) {
    // Если товар уже есть в корзине, увеличиваем количество
    cart.value[existingItemIndex].quantity += quantity;
  } else {
    // Иначе добавляем новый товар
    cart.value.push({
      id: product.id,
      name: product.name,
      price: product.price,
      image: product.image || '',
      quantity: quantity,
      currency: product.currency || 'RUB'
    });
  }
  
  // Обновляем счетчик и сохраняем корзину
  updateCartCounter();
  saveCart();
  
  return true;
};

// Обновление количества товара в корзине
const updateItemQuantity = (productId: number, quantity: number) => {
  const itemIndex = cart.value.findIndex(item => item.id === productId);
  
  if (itemIndex >= 0) {
    if (quantity > 0) {
      cart.value[itemIndex].quantity = quantity;
    } else {
      // Если количество 0 или отрицательное, удаляем товар
      removeFromCart(productId);
      return;
    }
    
    updateCartCounter();
    saveCart();
  }
};

// Удаление товара из корзины
const removeFromCart = (productId: number) => {
  cart.value = cart.value.filter(item => item.id !== productId);
  updateCartCounter();
  saveCart();
};

// Очистка корзины
const clearCart = () => {
  cart.value = [];
  updateCartCounter();
  saveCart();
};

// Расчет общей стоимости корзины
const calculateTotal = () => {
  return cart.value.reduce((total, item) => total + (item.price * item.quantity), 0);
};

// Автоматически загружаем корзину при создании модуля
if (process.client) {
  loadCart();
  
  // Следим за изменениями в корзине для обновления счетчика и сохранения
  watch(cart.value, () => {
    updateCartCounter();
    saveCart();
  }, { deep: true });
}

export const useCart = () => {
  return {
    cart,
    cartCounter,
    addProductToCart,
    updateItemQuantity,
    removeFromCart,
    clearCart,
    calculateTotal,
    loadCart
  };
};