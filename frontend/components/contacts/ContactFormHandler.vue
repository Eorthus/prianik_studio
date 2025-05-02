<template>
  <div class="tw-container tw-mx-auto tw-px-4">
    <h2
      class="tw-text-3xl tw-font-bold tw-text-gray-800 tw-text-center"
    >
      {{ $t("home.contact_us") }}
    </h2>
    <SubmitForm
      ref="submitFormRef"
      formType="contact"
      :initialData="initialFormData"
      class="tw-p-8 tw-mx-auto"
      @submit="handleSubmit"
      @success="handleSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useApiService } from '~/services/api';
import SubmitForm from './SubmitForm.vue';
import type { LangType } from '../types';

const { t, locale } = useI18n();
const emit = defineEmits<{
    success: []
}>();

// Возможно получение начальных данных извне
const props = defineProps<{
    initialData?: object
}>();

// API сервис для отправки формы
const { submitContactForm } = useApiService();

// Ссылка на форму
const submitFormRef = ref(null);

// Начальные данные для формы
const initialFormData = ref(props.initialData || {});

// Обработчик отправки формы
const handleSubmit = async ({ formData, formState, handleValidationErrors }) => {
  formState.value.isSubmitting = true;
  
  try {
    // Отправляем данные на сервер вместе с ответом reCAPTCHA
    const response = await submitContactForm({
      name: formData.name,
      email: formData.email,
      phone: formData.phone,
      message: formData.message,
      language: locale.value as LangType,
      recaptchaResponse: formData.recaptchaResponse // Передаем ответ капчи
    });
    
    if (response.success) {
      // Показываем успешное сообщение и сбрасываем форму
      submitFormRef.value.showSuccessAndReset();
    } else {
      // Обрабатываем ошибки
      if (response.validation_errors) {
        handleValidationErrors(response.validation_errors);
      } else {
        formState.value.errorMessage = response.error || t('form.error.general');
      }
      submitFormRef.value.showRejectAndReset();
    }
  } catch (error) {
    console.error('Ошибка при отправке формы обратной связи:', error);
    formState.value.errorMessage = t('form.error.general');
    formState.value.submitted = true;
    submitFormRef.value.showRejectAndReset();
  } finally {
    formState.value.isSubmitting = false;
  }
};

// Обработчик успешной отправки
const handleSuccess = () => {
  emit('success');
};
</script>