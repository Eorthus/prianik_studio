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

          <div ref="recaptchaContainer" class="tw-mb-4 recaptcha-container"></div>
          <!-- Сообщение об ошибке -->
          <p v-if="error" class="tw-text-red-500 tw-mt-2">{{ error }}</p>
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

// Функция для инициализации reCAPTCHA
const initRecaptcha = () => {
  if (!window.grecaptcha || !recaptchaContainer.value) {
    console.error("reCAPTCHA не загружена или контейнер не найден");
    return;
  }

  // Проверяем, что ключ доступен
  if (!siteKey) {
    console.error("reCAPTCHA site key не найден");
    return;
  }

  try {
    recaptchaId.value = window.grecaptcha.render(recaptchaContainer.value, {
      sitekey: siteKey,
      callback: (response: string) => {
        emit("verify", response);
        emit("close");
      },
      "expired-callback": () => {
        emit("expire");
      },
      "error-callback": () => {
        emit("error", "Ошибка reCAPTCHA");
      },
      size: "normal",
      theme: "light"
    });
  } catch (e) {
    console.error("Ошибка инициализации reCAPTCHA:", e);
  }
};

// Функция сброса капчи
const reset = () => {
  if (window.grecaptcha && recaptchaId.value !== null) {
    window.grecaptcha.reset(recaptchaId.value);
  }
};

// Функция для загрузки скрипта reCAPTCHA
const loadRecaptchaScript = () => {
  // Проверяем, был ли скрипт уже загружен
  if (window.grecaptcha) {
    initRecaptcha();
    return;
  }

  // Проверяем, что ключ загружен
  if (!siteKey) {
    console.error("reCAPTCHA site key не найден в конфигурации");
    return;
  }

  // Создаем скрипт reCAPTCHA с ключом в URL
  const script = document.createElement("script");
  script.src = `https://www.google.com/recaptcha/api.js?render=${siteKey}`;
  script.async = true;
  script.defer = true;

  // Обработчик загрузки скрипта
  script.onload = () => {
    // Ждем готовности grecaptcha
    if (window.grecaptcha && window.grecaptcha.ready) {
      window.grecaptcha.ready(() => {
        initRecaptcha();
      });
    } else {
      // Fallback если ready не доступен
      setTimeout(initRecaptcha, 100);
    }
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
  loadRecaptchaScript();
});

onUnmounted(() => {
  // Очистка - больше не нужна, так как не используем глобальный callback
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
    onRecaptchaLoaded?: () => void;
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