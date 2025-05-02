<template>
  <div class="tw-container tw-mx-auto tw-py-16 tw-px-4">
    <div class="tw-bg-white tw-shadow-md tw-rounded-lg tw-p-8 tw-text-center">
      <h1 class="tw-text-2xl tw-font-bold tw-text-gray-800 tw-mb-4">
        {{
          props.error.statusCode === 404 ? "Страница не найдена" : "Произошла ошибка"
        }}
      </h1>
      <p class="tw-text-gray-600 tw-mb-6">
        {{
          props.error.statusCode === 404
            ? "Запрашиваемая страница не существует"
            : "Произошла ошибка при обработке запроса"
        }}
      </p>
      <button
        @click="handleError"
        class="tw-bg-gray-800 tw-text-white tw-py-3 tw-px-6 tw-rounded-md tw-shadow-sm tw-font-medium tw-transition-colors tw-duration-300 hover:tw-bg-gray-700"
      >
        Вернуться на главную
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { clearError, useHead } from 'nuxt/app';

// Получаем данные об ошибке
const props = defineProps<{
  error: {
    statusCode:number
  };
}>();

// Обработка ошибки
const handleError = () => {
  clearError({ redirect: "/" });
};

// Устанавливаем HTTP статус-код
if (props.error?.statusCode) {
  useHead({
    title: `Ошибка ${props.error.statusCode}`,
    meta: [
      {
        name: "robots",
        content: "noindex, nofollow",
      },
    ],
  });
}
</script>
