<template>
  <div class="tw-mb-16">
    <h2
      class="tw-text-2xl tw-font-bold tw-text-gray-800 tw-mb-8 tw-text-center"
    >
      {{ props.title }}
    </h2>

    <div class="tw-relative tw-mx-auto">
      <!-- Слайдер -->
      <div class="tw-overflow-hidden tw-pb-3" ref="sliderContainer">
        <div
          class="tw-flex tw-transition-all tw-duration-500 tw-ease-out"
          :style="`transform: translateX(-${slideOffset}%);`"
        >
          <div
            v-for="product in props.products"
            :key="product.id"
            class="tw-prod-slide"
          >
            <div
              class="tw-bg-white tw-shadow-md tw-overflow-hidden tw-transition-all tw-duration-300 hover:tw-shadow-lg tw-mx-2 tw-h-full tw-flex tw-flex-col tw-w-full"
            >
              <NuxtLink
                :to="`/catalog/${product.id}`"
                class="tw-h-full tw-flex tw-flex-col"
              >
                <div class="tw-aspect-w-1 tw-aspect-h-1 tw-overflow-hidden">
                  <img
                    :src="product.images?.[0]"
                    :alt="product.name"
                    class="tw-w-full tw-h-full tw-object-cover tw-transition-transform tw-duration-500 hover:tw-scale-110"
                  />
                </div>
                <div class="tw-p-4 tw-flex tw-flex-col tw-flex-grow">
                  <h3
                    class="tw-text-lg tw-font-medium tw-text-gray-800 tw-mb-2 tw-line-clamp-2"
                  >
                    {{ product.name }}
                  </h3>
                  <p
                    class="tw-text-lg tw-font-semibold tw-text-gray-800 tw-mt-auto"
                  >
                    {{ product.price }} {{ currencyMap[product.currency || "RUB"] }}
                  </p>
                </div>
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>

      <!-- Кнопки навигации (показываем только если товаров больше чем может быть отображено) -->
      <template v-if="props.products.length > itemsPerView">
        <button
          @click="prevSlide"
          class="tw-absolute tw-left-4 tw-top-1/2 tw-transform -tw-translate-y-1/2 tw-p-2 tw-rounded-full tw-text-white tw-transition-all tw-duration-300 hover:tw-text-gray-200 focus:tw-outline-none"
          aria-label="Предыдущий слайд"
          :disabled="currentSlide === 0"
          :class="{ 'tw-opacity-50': currentSlide === 0 }"
        >
          <LeftArrowIcon class="tw-w-10 tw-h-10" />
        </button>
        <button
          @click="nextSlide"
          class="tw-absolute tw-right-4 tw-top-1/2 tw-transform -tw-translate-y-1/2 tw-p-2 tw-rounded-full tw-text-white tw-transition-all tw-duration-300 hover:tw-text-gray-200 focus:tw-outline-none"
          aria-label="Следующий слайд"
          :disabled="currentSlide >= maxSlide"
          :class="{ 'tw-opacity-50': currentSlide >= maxSlide }"
        >
          <RightArrowIcon class="tw-w-10 tw-h-10" />
        </button>
      </template>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import LeftArrowIcon from "../icons/LeftArrowIcon.vue";
import RightArrowIcon from "../icons/RightArrowIcon.vue";
import { currencyMap, type Product } from "~/components";

const props = withDefaults(
  defineProps<{
    products: Product[];
    title: string;
  }>(),
  {
    title: "Похожие товары",
  }
);

const sliderContainer = ref(null);

// Текущий слайд
const currentSlide = ref(0);

// Количество отображаемых элементов в зависимости от размера экрана
const itemsPerView = ref(4);

// Максимальное количество слайдов
const maxSlide = computed(() => {
  return Math.max(0, props.products.length - itemsPerView.value);
});

// Смещение слайдера
const slideOffset = computed(() => {
  // Вычисляем ширину каждого слайда (в процентах)
  const slideWidth = 100 / itemsPerView.value;
  return currentSlide.value * slideWidth;
});

// Функция для перехода к предыдущему слайду
const prevSlide = () => {
  if (currentSlide.value > 0) {
    currentSlide.value = Math.max(currentSlide.value - 1, 0);
  } else {
    // Циклический переход к последнему слайду
    currentSlide.value = maxSlide.value;
  }
};

// Функция для перехода к следующему слайду
const nextSlide = () => {
  if (currentSlide.value < maxSlide.value) {
    currentSlide.value = Math.min(currentSlide.value + 1, maxSlide.value);
  } else {
    // Циклический переход к первому слайду
    currentSlide.value = 0;
  }
};

// Определение количества видимых элементов в зависимости от ширины экрана
const updateItemsPerView = () => {
  const width = window?.innerWidth;
  if (width >= 1280) {
    // xl и больше
    itemsPerView.value = 4;
  } else if (width >= 1024) {
    // lg
    itemsPerView.value = 3;
  } else if (width >= 768) {
    // md
    itemsPerView.value = 2;
  } else if (width >= 640) {
    // sm
    itemsPerView.value = 2;
  } else {
    // xs
    itemsPerView.value = 1.2; // Показываем один товар и часть следующего
  }

  // Корректируем текущий слайд, если после изменения размера экрана он вышел за допустимые пределы
  if (currentSlide.value > maxSlide.value) {
    currentSlide.value = maxSlide.value;
  }
};

// Автоматическое переключение слайдов
let autoSlideInterval = null;
const startAutoSlide = () => {
  // Автоматически переключать слайды каждые 5 секунд
  stopAutoSlide(); // Сначала останавливаем, если уже запущено

  autoSlideInterval = setInterval(() => {
    nextSlide();
  }, 5000);
};

const stopAutoSlide = () => {
  if (autoSlideInterval) {
    clearInterval(autoSlideInterval);
    autoSlideInterval = null;
  }
};

const pauseAutoSlide = () => {
  stopAutoSlide();
};

const resumeAutoSlide = () => {
  startAutoSlide();
};

onMounted(() => {
  // Инициализация количества отображаемых элементов
  updateItemsPerView();

  // Обновление при изменении размера окна
  window.addEventListener("resize", updateItemsPerView);

  // Запуск автоматического переключения слайдов
  startAutoSlide();

  // Приостанавливать автопрокрутку при наведении мыши
  if (sliderContainer.value) {
    sliderContainer.value.addEventListener("mouseenter", pauseAutoSlide);
    sliderContainer.value.addEventListener("mouseleave", resumeAutoSlide);

    // Также обрабатываем касания для мобильных устройств
    sliderContainer.value.addEventListener("touchstart", pauseAutoSlide);
    sliderContainer.value.addEventListener("touchend", resumeAutoSlide);
  }
});

onUnmounted(() => {
  // Удаление обработчиков событий и таймеров при уничтожении компонента
  window.removeEventListener("resize", updateItemsPerView);
  stopAutoSlide();

  if (sliderContainer.value) {
    sliderContainer.value.removeEventListener("mouseenter", pauseAutoSlide);
    sliderContainer.value.removeEventListener("mouseleave", resumeAutoSlide);
    sliderContainer.value.removeEventListener("touchstart", pauseAutoSlide);
    sliderContainer.value.removeEventListener("touchend", resumeAutoSlide);
  }
});
</script>

<style scoped>
/* Стиль для карточек товаров с фиксированным соотношением сторон */
.tw-aspect-w-1 {
  position: relative;
  padding-bottom: 100%; /* 1:1 */
}

.tw-aspect-w-1 > * {
  position: absolute;
  height: 100%;
  width: 100%;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  object-fit: cover;
}

/* Ограничение высоты текста (2 строки) */
.tw-line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Стиль для слайдов товаров */
.tw-prod-slide {
  flex: 0 0 auto;
  width: calc(100% / v-bind(itemsPerView));
  height: auto;
  display: flex;
}

/* Улучшение стиля для кнопок навигации на маленьких экранах */
@media (max-width: 640px) {
  button[aria-label="Предыдущий слайд"],
  button[aria-label="Следующий слайд"] {
    padding: 0.5rem;
    background-color: rgba(31, 41, 55, 0.8);
  }
}
</style>
