<template>
  <!-- Форма обратной связи -->
  <div class="tw-max-w-2xl tw-rounded-lg">
    <!-- Модальное окно результата отправки формы -->
    <SubmitModal
      :is-open="formState.submitted"
      :error-message="formState.errorMessage"
    />

    <!-- Модальное окно reCAPTCHA -->
    <ReCaptchaModal
      :is-open="showCaptchaModal"
      @close="closeCaptchaModal"
      @verify="onCaptchaVerify"
      @error="onCaptchaError"
    />

    <!-- Форма -->
    <form @submit.prevent="startFormSubmit" class="tw-space-y-6">
      <input type="hidden" name="_csrf" :value="csrfToken" />

      <!-- Имя -->
      <LInput
        label="form.name"
        placeholder="form.name_placeholder"
        v-model="formData.name"
        :errors="formState.errors"
        type="text"
        name="name"
      />

      <!-- Email -->
      <LInput
        label="form.email"
        :placeholder="emailComputed"
        v-model="formData.email"
        :errors="formState.errors.email"
        type="email"
        name="email"
      />

      <!-- Телефон -->
      <LInput
        label="form.phone"
        placeholder="form.phone_placeholder"
        v-model="formData.phone"
        :errors="formState.errors.phone"
        type="tel"
        name="phone"
      />

      <!-- Сообщение (текстовое поле) или Комментарий (textarea) в зависимости от типа формы -->
      <LTextarea
        v-if="props.formType === 'contact'"
        name="message"
        v-model="formData.message"
        label="form.message"
        placeholder="form.message_placeholder"
        :error="formState.errors.message"
        required
      />

      <!-- Комментарий к заказу (необязательное поле) -->
      <LTextarea
        v-if="props.formType === 'order'"
        name="comment"
        v-model="formData.comment"
        label="cart.comment"
        placeholder="cart.comment_placeholder"
        :error="formState.errors.comment"
      />

      <!-- Общая ошибка -->
      <div
        v-if="formState.errors.general"
        class="tw-bg-red-50 tw-p-4 tw-rounded-md"
      >
        <p class="tw-text-red-500">{{ formState.errors.general }}</p>
      </div>

      <!-- Кнопка отправки -->
      <div>
        <button
          type="submit"
          class="tw-w-full tw-uppercase tw-bg-gray-800 tw-text-white tw-py-3 tw-px-4 tw-rounded-md tw-shadow-sm tw-text-lg tw-font-medium tw-transition-colors tw-duration-300 hover:tw-bg-gray-700 focus:tw-outline-none focus:tw-ring-2 focus:tw-ring-offset-2 focus:tw-ring-gray-500"
          :disabled="formState.isSubmitting"
        >
          <span
            v-if="formState.isSubmitting"
            class="tw-flex tw-items-center tw-justify-center"
          >
            <SpinnerIcon class="tw-animate-spin tw-h-5 tw-w-5 tw-mr-3" />
            {{ $t("form.sending") }}
          </span>
          <span v-else>{{ $t("form.submit") }}</span>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import SubmitModal from "./SubmitModal.vue";
import ReCaptchaModal from "../ReCaptchaModal.vue";
import { useCsrf } from "~/shared/useCsrf";
import { sanitizeInput } from "~/utils/security";
import SpinnerIcon from "../icons/SpinnerIcon.vue";
import LInput from "~/lib/LInput.vue";
import LTextarea from "~/lib/LTextarea.vue";

//@ts-expect-error 
const { t } = useI18n();
const emit = defineEmits<{
  submit: [formData: object];
  success: [];
}>();

// Создаем и получаем CSRF-токен
const { csrfToken, generateToken } = useCsrf();

// Состояние модального окна reCAPTCHA
const showCaptchaModal = ref(false);
const captchaResponse = ref("");
const captchaVerified = ref(false);

// Props компонента
const props = withDefaults(
  defineProps<{
    formType: "contact" | "order";
    // Начальные данные формы (можно передать извне)
    initialData: object;
  }>(),
  {
    formType: "contact",
  }
);

// Данные формы
const formData = ref({
  name: "",
  email: "",
  phone: "",
  message: "", // Для формы обратной связи
  comment: "", // Для формы заказа
});

// Состояние формы
const formState = ref({
  isSubmitting: false,
  submitted: false,
  errorMessage: "",
  errors: {
    name: "",
    email: "",
    phone: "",
    message: "",
    comment: "",
    recaptcha: "",
    general: "",
  },
});

const emailComputed = computed(()=>`${t('form.email_placeholder_prefix')}@${t('form.email_placeholder_main')}`)

// Открытие/закрытие модального окна с reCAPTCHA
const openCaptchaModal = () => {
  showCaptchaModal.value = true;
};

const closeCaptchaModal = () => {
  showCaptchaModal.value = false;

  // Если капча не была пройдена, сбрасываем состояние отправки
  if (!captchaVerified.value) {
    formState.value.isSubmitting = false;
  }
};

// Обработчики событий reCAPTCHA
const onCaptchaVerify = (response: string) => {
  captchaResponse.value = response;
  captchaVerified.value = true;
  formState.value.errors.recaptcha = "";

  // После успешной верификации продолжаем отправку формы
  onSubmit();
};

const onCaptchaError = (error: string) => {
  captchaVerified.value = false;
  captchaResponse.value = "";
  formState.value.errors.recaptcha = t("form.error.recaptcha_error");
  formState.value.isSubmitting = false;
  console.error("reCAPTCHA error:", error);
};

// Валидация email
const validateEmail = (email: string) => {
  const re = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return re.test(email);
};

// Валидация телефона
const validatePhone = (phone: string) => {
  // Упрощенная проверка - ожидаем наличие цифр и общую длину > 7
  return /^[+]?[\d\s()-]{7,}$/.test(phone);
};

// Очистка ошибок
const clearErrors = () => {
  formState.value.errors = {
    name: "",
    email: "",
    phone: "",
    message: "",
    comment: "",
    recaptcha: "",
    general: "",
  };
  formState.value.errorMessage = "";
};

// Сброс формы
const resetForm = () => {
  formData.value = {
    name: "",
    email: "",
    phone: "",
    message: "",
    comment: "",
  };
  clearErrors();
  formState.value.submitted = false;

  // Сбрасываем состояние капчи
  captchaVerified.value = false;
  captchaResponse.value = "";
};

// Валидация формы
const validateForm = () => {
  clearErrors();
  let isValid = true;

  // Проверяем имя
  if (!formData.value.name.trim()) {
    formState.value.errors.name = t("form.error.name_required");
    isValid = false;
  }

  // Проверяем email
  if (!formData.value.email.trim()) {
    formState.value.errors.email = t("form.error.email_required");
    isValid = false;
  } else if (!validateEmail(formData.value.email)) {
    formState.value.errors.email = t("form.error.email_invalid");
    isValid = false;
  }

  // Проверяем телефон
  if (!formData.value.phone.trim()) {
    formState.value.errors.phone = t("form.error.phone_required");
    isValid = false;
  } else if (!validatePhone(formData.value.phone)) {
    formState.value.errors.phone = t("form.error.phone_invalid");
    isValid = false;
  }

  // Проверяем сообщение для формы обратной связи
  if (props.formType === "contact" && !formData.value.message.trim()) {
    formState.value.errors.message = t("form.error.message_required");
    isValid = false;
  }

  return isValid;
};

// Начало процесса отправки формы
const startFormSubmit = () => {
  // Проверяем валидность формы
  if (!validateForm()) return;

  // Отмечаем начало отправки
  formState.value.isSubmitting = true;

  // Открываем модальное окно с капчей
  openCaptchaModal();
};

// Санитизация данных формы
const sanitizeFormData = () => {
  const sanitizedData: Record<string, string> = {};

  // Общие поля для всех типов форм
  sanitizedData.name = sanitizeInput(formData.value.name);
  sanitizedData.email = sanitizeInput(formData.value.email);
  sanitizedData.phone = sanitizeInput(formData.value.phone);
  sanitizedData.recaptchaResponse = captchaResponse.value;

  // Специфичные поля в зависимости от типа формы
  if (props.formType === "contact") {
    sanitizedData.message = sanitizeInput(formData.value.message);
  } else if (props.formType === "order") {
    sanitizedData.comment = formData.value.comment
      ? sanitizeInput(formData.value.comment)
      : "";
  }

  return sanitizedData;
};

// Обработка ошибок валидации от сервера
const handleValidationErrors = (
  validationErrors: Array<{ field: string; message: string }>
) => {
  if (Array.isArray(validationErrors)) {
    validationErrors.forEach((error) => {
      // Преобразуем названия полей к нижнему регистру, так как бэкенд может вернуть их с заглавной буквы
      const fieldName = error.field.toLowerCase();
      if (fieldName in formState.value.errors) {
        formState.value.errors[fieldName] = error.message;
      }
    });
  }
};

// Обработчик отправки формы (вызывается после успешной капчи)
const onSubmit = () => {
  // Санитизируем данные
  const sanitizedData = sanitizeFormData();

  // Эмитим событие submit с данными формы для обработки родительским компонентом
  emit("submit", {
    formData: sanitizedData,
    formState,
    handleValidationErrors,
  });
};

// При успешной отправке показываем модальное окно и сбрасываем форму
const showSuccessAndReset = () => {
  formState.value.submitted = true;
  formState.value.isSubmitting = false;

  // Через 3 секунды скрываем модальное окно и сбрасываем форму
  setTimeout(() => {
    resetForm();
    emit("success");
  }, 3000);
};

const showRejectAndReset = () => {
  formState.value.submitted = true;
  formState.value.isSubmitting = false;

  setTimeout(() => {
    formState.value.submitted = false;
    emit("success");
  }, 3000);
};

// Инициализация компонента
onMounted(() => {
  generateToken(); // Генерируем CSRF токен

  // Заполняем начальные данные, если они переданы
  if (props.initialData) {
    formData.value = { ...formData.value, ...props.initialData };
  }
});

// Экспортируем методы для использования в родительском компоненте
defineExpose({
  resetForm,
  showSuccessAndReset,
  showRejectAndReset,
  formState,
  formData,
});
</script>
