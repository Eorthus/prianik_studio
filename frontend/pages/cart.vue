<script setup lang="ts">
import CartIcon from "~/components/icons/CartIcon.vue";
import { computed, onMounted } from "vue";
import MinusIcon from "~/components/icons/MinusIcon.vue";
import PlusIcon from "~/components/icons/PlusIcon.vue";
import DeleteIcon from "~/components/icons/DeleteIcon.vue";
import { useCart } from "~/shared/useCart";
import { currencyMap } from "~/components";
import OrderFormHandler from "~/components/card/OrderFormHandler.vue";

// Подключаем корзину
const {
  cart,
  updateItemQuantity,
  removeFromCart,
  calculateTotal,
  loadCart,
} = useCart();

// Изменение количества товара
const updateQuantity = (id, newQuantity) => {
  if (newQuantity > 0) {
    updateItemQuantity(id, newQuantity);
  }
};

// Удаление товара из корзины
const removeItem = (id) => {
  removeFromCart(id);
};

// Общая стоимость корзины
const totalPrice = computed(() => {
  return calculateTotal();
});

// Проверка, пуста ли корзина
const isCartEmpty = computed(() => !cart.value || cart.value.length === 0);

// Обработчик успешного оформления заказа
const handleOrderSuccess = () => {
  // Дополнительные действия после успешного оформления заказа, если необходимо
  console.log("Заказ успешно оформлен");
};

// Загружаем корзину при монтировании компонента
onMounted(() => {
  // В случае, если корзина не была загружена ранее
  if (isCartEmpty.value) {
    loadCart();
  }
});
</script>

<template>
  <div class="tw-py-12">
    <div class="tw-container tw-mx-auto tw-px-4">
      <h1
        class="tw-text-3xl tw-font-bold tw-text-gray-800 tw-mb-8 tw-text-center"
      >
        {{ $t("cart.title") }}
      </h1>

      <div
        v-if="!isCartEmpty"
        class="tw-flex tw-flex-col lg:tw-flex-row tw-gap-8"
      >
        <!-- Список товаров в корзине -->
        <div class="lg:tw-w-2/3">
          <div class="tw-bg-white tw-overflow-hidden">
            <div class="tw-p-6">
              <div class="tw-space-y-4">
                <div
                  v-for="item in cart"
                  :key="item.id"
                  class="tw-flex tw-flex-col sm:tw-flex-row tw-items-center tw-gap-4 tw-pb-4 tw-border-gray-200"
                  :class="{
                    'tw-border-b': item.id !== cart[cart.length - 1].id,
                  }"
                >
                  <!-- Изображение товара -->
                  <div class="sm:tw-w-20 sm:tw-h-20">
                    <img
                      :src="
                        item.image ||
                        'https://s3.stroi-news.ru/img/masterskaya-kartinki-1.jpg'
                      "
                      :alt="item.name"
                      class="tw-w-full tw-h-full tw-object-cover"
                    />
                  </div>

                  <!-- Информация о товаре -->
                  <div class="tw-flex-grow tw-text-center sm:tw-text-left">
                    <h3
                      class="tw-text-base tw-font-medium tw-text-gray-800 tw-mb-1"
                    >
                      {{ item.name }}
                    </h3>
                    <p class="tw-text-gray-600">
                      {{ item.price }} {{ currencyMap[item.currency || "RUB"] }}
                    </p>
                  </div>

                  <!-- Количество -->
                  <div class="tw-flex tw-items-center">
                    <button
                      @click="updateQuantity(item.id, item.quantity - 1)"
                      class="tw-bg-gray-200 tw-text-gray-800 tw-px-3 tw-py-2 tw-rounded-l-md tw-transition-colors hover:tw-bg-gray-300"
                      :disabled="item.quantity <= 1"
                    >
                      <MinusIcon class="tw-w-3 tw-h-3" />
                    </button>
                    <span class="tw-w-8 tw-text-center tw-py-1">{{
                      item.quantity
                    }}</span>
                    <button
                      @click="updateQuantity(item.id, item.quantity + 1)"
                      class="tw-bg-gray-200 tw-text-gray-800 tw-px-3 tw-py-2 tw-rounded-r-md tw-transition-colors hover:tw-bg-gray-300"
                    >
                      <PlusIcon class="tw-w-3 tw-h-3" />
                    </button>
                  </div>

                  <!-- Сумма -->
                  <div class="tw-w-24 tw-text-right">
                    <p class="tw-font-semibold tw-text-gray-800">
                      {{ item.price * item.quantity }}
                      {{ currencyMap[item.currency || "RUB"] }}
                    </p>
                  </div>

                  <!-- Кнопка удаления -->
                  <button
                    @click="removeItem(item.id)"
                    class="tw-text-gray-500 hover:tw-text-red-500 tw-transition-colors"
                    :aria-label="$t('cart.remove_item')"
                  >
                    <DeleteIcon class="tw-w-5 tw-h-5" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Форма оформления заказа -->
          <div class="tw-mt-8 tw-bg-white tw-overflow-hidden">
            <div class="tw-p-6">
              <h2 class="tw-text-xl tw-font-semibold tw-text-gray-800">
                {{ $t("cart.order_info") }}
              </h2>

              <!-- Используем новый компонент обработчика формы заказа -->
              <OrderFormHandler @success="handleOrderSuccess" />
            </div>
          </div>
        </div>

        <!-- Итоговая информация и оформление заказа -->
        <div class="lg:tw-w-1/3">
          <div
            class="tw-bg-white tw-shadow-md tw-overflow-hidden tw-sticky tw-top-24"
          >
            <div class="tw-p-6">
              <h2 class="tw-text-xl tw-font-semibold tw-text-gray-800 tw-mb-4">
                {{ $t("cart.your_order") }}
              </h2>

              <div class="tw-space-y-4">
                <div
                  class="tw-flex tw-justify-between tw-border-b tw-border-gray-200 tw-pb-4"
                >
                  <span class="tw-text-gray-600"
                    >{{ $t("cart.items") }} ({{ cart.length }})</span
                  >
                  <span class="tw-font-medium tw-text-gray-800"
                    >{{ totalPrice }} {{  currencyMap[cart?.[0]?.currency || "RUB"]  }}</span
                  >
                </div>

                <div
                  class="tw-flex tw-justify-between tw-border-b tw-border-gray-200 tw-pb-4"
                >
                  <span class="tw-text-gray-600">{{
                    $t("cart.delivery")
                  }}</span>
                  <span class="tw-font-medium tw-text-gray-800">{{
                    $t("cart.delivery_cost")
                  }}</span>
                </div>

                <div class="tw-flex tw-justify-between tw-pt-2">
                  <span class="tw-text-lg tw-font-semibold tw-text-gray-800">{{
                    $t("cart.total")
                  }}</span>
                  <span class="tw-text-lg tw-font-semibold tw-text-gray-800"
                    >{{ totalPrice }} {{  currencyMap[cart?.[0]?.currency || "RUB"]  }}</span
                  >
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Пустая корзина -->
      <div v-else class="tw-bg-white tw-p-8 tw-text-center">
        <div class="tw-flex tw-justify-center tw-mb-4">
          <CartIcon class="tw-w-16 tw-h-16 tw-text-gray-400" />
        </div>

        <h2 class="tw-text-2xl tw-font-semibold tw-text-gray-800 tw-mb-2">
          {{ $t("cart.empty_cart") }}
        </h2>
        <p class="tw-text-gray-600 tw-mb-6">
          {{ $t("cart.add_items") }}
        </p>

        <NuxtLink
          to="/catalog"
          class="tw-bg-gray-800 tw-text-white tw-py-3 tw-px-6 tw-rounded-md tw-shadow-sm tw-font-medium tw-transition-colors tw-duration-300 hover:tw-bg-gray-700 focus:tw-outline-none focus:tw-ring-2 focus:tw-ring-offset-2 focus:tw-ring-gray-500 tw-inline-block"
        >
          {{ $t("cart.go_to_catalog") }}
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Стили для прилипающей боковой панели */
.tw-sticky {
  position: sticky;
}
</style>
