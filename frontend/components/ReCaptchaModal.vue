<template>
  <transition name="modal-fade">
    <div
      v-show="props.isOpen"
      class="tw-fixed tw-inset-0 tw-bg-black/50 tw-flex tw-items-center tw-justify-center tw-z-50"
    >
      <div
        class="tw-bg-white tw-rounded-lg tw-shadow-xl tw-w-full tw-p-6 tw-max-w-[350px]"
      >
        <!-- Заголовок модального окна -->
        <div class="tw-flex tw-justify-between tw-items-center tw-mb-4">
          <h3 class="tw-text-xl tw-font-semibold tw-text-gray-800">
            {{ $t("form.verify_captcha") }}
          </h3>
          <button
            @click="closeModal"
            class="tw-text-gray-500 hover:tw-text-gray-700 focus:tw-outline-none"
          >
            <CrossIcon class="tw-w-6 tw-h-6" />
          </button>
        </div>

        <!-- Содержимое модального окна -->
        <div class="tw-mb-4">
          <p class="tw-text-gray-600 tw-mb-4">
            {{ $t("form.captcha_explanation") }}
          </p>

          <!-- Контейнер для reCAPTCHA v2 -->
          <div ref="recaptchaContainer" class="tw-mb-4 recaptcha-container tw-flex tw-justify-center"></div>
          
          <!-- Сообщение об ошибке -->
          <p v-if="error" class="tw-text-red-500 tw-mt-2 tw-text-center">{{ error }}</p>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import CrossIcon from "./icons/CrossIcon.vue";
import { useRuntimeConfig } from "nuxt/app";

// Получаем конфигурацию
const { public: { recaptchaSiteKey } } = useRuntimeConfig();

// Props
const props = defineProps<{
  isOpen: boolean;
}>();

// Уникальный ID для этого экземпляра reCAPTCHA
const recaptchaId = ref<string | null>(null);
const recaptchaContainer = ref<HTMLElement | null>(null);
const error = ref<string>("");

// ВАЖНО: Используем ключ из runtime config
const siteKey = recaptchaSiteKey;

// Эмиты
const emit = defineEmits(["verify", "expire", "error", "close"]);

// Функция для инициализации reCAPTCHA v2
const initRecaptcha = () => {
  if (!window.grecaptcha || !recaptchaContainer.value) {
    console.error("reCAPTCHA не загружена или контейнер не найден");
    return;
  }

  // Проверяем, что ключ доступен
  if (!siteKey) {
    console.error("reCAPTCHA site key не найден");
    error.value = "Ошибка конфигурации reCAPTCHA";
    return;
  }

  try {
    // Очищаем контейнер перед рендерингом
    recaptchaContainer.value.innerHTML = '';
    
    recaptchaId.value = window.grecaptcha.render(recaptchaContainer.value, {
      sitekey: siteKey,
      callback: (response: string) => {
        console.log("reCAPTCHA v2 успешно пройдена");
        emit("verify", response);
        emit("close");
      },
      "expired-callback": () => {
        console.log("reCAPTCHA истекла");
        emit("expire");
        error.value = "Время проверки истекло. Попробуйте снова.";
      },
      "error-callback": () => {
        console.log("Ошибка reCAPTCHA");
        emit("error", "Ошибка reCAPTCHA");
        error.value = "Ошибка проверки. Попробуйте обновить страницу.";
      },
      size: "normal",
      theme: "light"
    });
    
    console.log("reCAPTCHA v2 инициализирована с ID:", recaptchaId.value);
  } catch (e) {
    console.error("Ошибка инициализации reCAPTCHA:", e);
    error.value = "Ошибка загрузки проверки безопасности";
  }
};

// Функция сброса капчи
const reset = () => {
  if (window.grecaptcha && recaptchaId.value !== null) {
    try {
      window.grecaptcha.reset(recaptchaId.value);
      error.value = "";
    } catch (e) {
      console.error("Ошибка сброса reCAPTCHA:", e);
    }
  }
};

// Функция для загрузки скрипта reCAPTCHA v2
const loadRecaptchaScript = () => {
  // Проверяем, был ли скрипт уже загружен
  if (window.grecaptcha) {
    console.log("reCAPTCHA уже загружена");
    initRecaptcha();
    return;
  }

  // Проверяем, что ключ загружен
  if (!siteKey) {
    console.error("reCAPTCHA site key не найден в конфигурации");
    error.value = "Ошибка конфигурации reCAPTCHA";
    return;
  }

  console.log("Загрузка reCAPTCHA v2 с ключом:", siteKey.substring(0, 20) + "...");

  // Создаем скрипт reCAPTCHA v2 (БЕЗ параметра render с ключом!)
  const script = document.createElement("script");
  script.src = `https://www.google.com/recaptcha/api.js?onload=onRecaptchaLoadedV2&render=explicit`;
  script.async = true;
  script.defer = true;

  // Определяем глобальную функцию обратного вызова
  window.onRecaptchaLoadedV2 = () => {
    console.log("reCAPTCHA v2 загружена, инициализация...");
    setTimeout(() => {
      initRecaptcha();
    }, 100);
  };

  script.onerror = () => {
    error.value = "Ошибка загрузки reCAPTCHA";
    console.error("Ошибка загрузки скрипта reCAPTCHA");
  };

  // Добавляем скрипт на страницу
  document.head.appendChild(script);
};

// Функция закрытия модального окна
const closeModal = () => {
  emit("close");
};

// Хук жизненного цикла
onMounted(() => {
  console.log("ReCaptchaModal смонтирован");
  loadRecaptchaScript();
});

onUnmounted(() => {
  // Очистка глобальной функции
  if (window.onRecaptchaLoadedV2) {
    window.onRecaptchaLoadedV2 = undefined;
  }
});

// Экспортируем методы
defineExpose({
  reset,
});
</script>

<script lang="ts">
// Определяем типы для интерфейса window
declare global {
  interface Window {
    grecaptcha: any;
    onRecaptchaLoadedV2?: () => void;
  }
}
</script>

<style scoped>
/* Скрываем предупреждающий текст для тестового ключа в режиме разработки */
.recaptcha-container :deep(.rc-anchor-error-msg-container) {
  display: none !important;
}

.recaptcha-container :deep(.rc-anchor-alert) {
  display: none !important;
}

/* Скрываем любые красные предупреждения */
.recaptcha-container :deep([style*="color: rgb(255, 0, 0)"]) {
  display: none !important;
}

.recaptcha-container :deep([style*="color:#ff0000"]) {
  display: none !important;
}

/* Анимация для модального окна */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>