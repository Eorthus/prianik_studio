<template>
  <div>
    <!-- Уведомление/нотификация о процессе заказа -->
    <div
      class="tw-bg-gray-100 tw-border-l-4 tw-border-gray-800 tw-p-3 tw-rounded-md tw-max-w-2xl tw-shadow-sm tw-mt-4"
    >
      <div class="tw-flex tw-items-start">
        <div class="tw-flex-shrink-0 tw-pt-0.5">
          <InfoIcon class="tw-h-5 tw-w-5 tw-text-gray-800" />
        </div>
        <div class="tw-ml-3 tw-flex-1">
          <h3 class="tw-text-lg tw-font-medium tw-text-gray-800">
            {{ $t("order_form.notification_title") }}
          </h3>
          <div class="tw-mt-2 tw-text-gray-600">
            <p>{{ $t("order_form.notification_text") }}</p>
          </div>
        </div>
      </div>
    </div>
    <SubmitForm
      ref="submitFormRef"
      formType="order"
      :initialData="initialFormData"
      @submit="handleSubmit"
      @success="handleSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useI18n } from "vue-i18n";
import { useApiService } from "~/services/api";
import { useCart } from "~/shared/useCart";
import SubmitForm from "../contacts/SubmitForm.vue";
import InfoIcon from "../icons/InfoIcon.vue";
import type { LangType } from "../types";

const { t, locale } = useI18n();
const emit = defineEmits<{
  success: [];
}>();

// Возможно получение начальных данных извне
const props = defineProps<{
  initialData?: object;
  // Если нужно заказать конкретный товар (для быстрого заказа с карточки товара)
  productToOrder?: {
    id: number;
    quantity: number;
  };
}>();

// API сервис для отправки заказа
const { createOrder } = useApiService();

// Корзина (для обычного заказа из корзины)
const { cart, clearCart } = useCart();

// Ссылка на форму
const submitFormRef = ref(null);

// Начальные данные для формы
const initialFormData = ref(props.initialData || {});

// Проверяем, какие товары заказываются: один товар или корзина
const orderItems = computed(() => {
  if (props.productToOrder) {
    // Быстрый заказ одного товара
    return [
      {
        product_id: props.productToOrder.id,
        quantity: props.productToOrder.quantity || 1,
      },
    ];
  } else {
    // Заказ всех товаров из корзины
    return cart.value.map((item) => ({
      product_id: item.id,
      quantity: item.quantity,
    }));
  }
});

// Обработчик отправки формы
const handleSubmit = async ({
  formData,
  formState,
  handleValidationErrors,
}) => {
  // Проверяем наличие товаров для заказа
  if (orderItems.value.length === 0) {
    formState.value.errors.general = t("cart.error.empty_cart");
    return;
  }

  formState.value.isSubmitting = true;

  try {
    // Формируем данные заказа, включая ответ reCAPTCHA
    const orderData = {
      name: formData.name,
      email: formData.email,
      phone: formData.phone,
      comment: formData.comment || "",
      items: orderItems.value,
      language: locale.value as LangType,
      recaptchaResponse: "", // Добавляем ответ капчи
    };

    // Отправляем заказ на сервер
    const response = await createOrder(orderData);

    if (response.success) {
      // Если заказ был из корзины, очищаем её
      if (!props.productToOrder) {
        clearCart();
      }

      // Показываем успешное сообщение и сбрасываем форму
      submitFormRef.value.showSuccessAndReset();
    } else {
      // Обрабатываем ошибки
      if (response.validation_errors) {
        handleValidationErrors(response.validation_errors);
      } else {
        formState.value.errorMessage =
          response.error || t("cart.error.general");
        formState.value.submitted = true; // Показываем модальное окно с ошибкой
      }
      submitFormRef.value.showRejectAndReset();
    }
  } catch (error) {
    console.error("Ошибка при оформлении заказа:", error);
    formState.value.errorMessage = t("cart.error.general");
    formState.value.submitted = true; // Показываем модальное окно с ошибкой
    submitFormRef.value.showRejectAndReset();
  } finally {
    formState.value.isSubmitting = false;
  }
};

// Обработчик успешной отправки
const handleSuccess = () => {
  emit("success");
};
</script>
