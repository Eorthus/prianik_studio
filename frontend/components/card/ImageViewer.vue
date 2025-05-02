<template>
  <transition name="modal">
    <div
      v-if="props.isOpen"
      class="tw-fixed tw-inset-0 tw-bg-black/90 tw-flex tw-items-center tw-justify-center tw-z-50 tw-p-4"
      @click.self="close"
    >
      <!-- Кнопка закрытия -->
      <button
        @click="close"
        class="tw-absolute tw-top-4 tw-right-4 tw-text-white tw-transition-colors tw-duration-300 hover:tw-text-gray-300 focus:tw-outline-none"
        aria-label="Закрыть"
      >
        <CrossIcon class="tw-w-8 tw-h-8" />
      </button>

      <div
        class="tw-relative tw-max-w-5xl tw-w-full tw-h-full tw-max-h-[80vh] tw-flex tw-items-center tw-justify-center"
      >
        <!-- Навигация влево -->
        <button
          v-if="props.images.length > 1"
          @click.stop="prevImage"
          class="tw-absolute tw-left-2 tw-z-20 tw-p-2 tw-rounded-full tw-text-white tw-transition-colors tw-duration-300 hover:tw-text-gray-300 focus:tw-outline-none"
          aria-label="Предыдущее изображение"
        >
          <LeftArrowIcon class="tw-w-10 tw-h-10" />
        </button>

        <!-- Контейнер изображения -->
        <div
          class="tw-overflow-hidden tw-relative tw-flex tw-items-center tw-justify-center tw-w-full tw-h-full"
          :class="{
            'tw-cursor-zoom-in': !isZoomed,
            'tw-cursor-zoom-out': isZoomed,
          }"
          @click.stop="toggleZoom"
        >
          <div
            :class="{ 'tw-transition-transform tw-duration-300': !isDragging }"
            :style="containerStyle"
            @mousedown="startDrag"
            @mousemove="drag"
            @mouseup="endDrag"
            @mouseleave="endDrag"
            @touchstart="startDrag"
            @touchmove="drag"
            @touchend="endDrag"
            @touchcancel="endDrag"
          >
            <img
              :src="currentImage"
              :alt="`Просмотр изображения ${currentIndex + 1} из ${
                images.length
              }`"
              class="tw-max-w-full tw-max-h-full tw-object-contain"
              @dragstart.prevent
            />
          </div>
        </div>

        <!-- Навигация вправо -->
        <button
          v-if="props.images.length > 1"
          @click.stop="nextImage"
          class="tw-absolute tw-right-2 tw-z-20 tw-p-2 tw-rounded-full tw-text-white tw-transition-colors tw-duration-300 hover:tw-text-gray-300 focus:tw-outline-none"
          aria-label="Следующее изображение"
        >
          <RightArrowIcon class="tw-w-10 tw-h-10" />
        </button>

        <!-- Индикаторы -->
        <div
          v-if="props.images.length > 1"
          class="tw-absolute tw-bottom-4 tw-left-0 tw-right-0 tw-flex tw-justify-center tw-space-x-2"
        >
          <button
            v-for="(_, index) in props.images"
            :key="index"
            @click.stop="currentIndex = index"
            :class="[
              'tw-w-2 tw-h-2 tw-rounded-full tw-transition-all',
              currentIndex === index ? 'tw-bg-white' : 'tw-bg-gray-400',
            ]"
            :aria-label="`Изображение ${index + 1}`"
          ></button>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import CrossIcon from "../icons/CrossIcon.vue";
import LeftArrowIcon from "../icons/LeftArrowIcon.vue";
import RightArrowIcon from "../icons/RightArrowIcon.vue";

const props = defineProps<{
  images: string[];
  isOpen: boolean;
  initialIndex: number;
}>();

const emit = defineEmits<{
  close: [];
}>();

// Текущий индекс изображения
const currentIndex = ref(props.initialIndex);

// Текущее состояние увеличения
const isZoomed = ref(false);

// Положение изображения при увеличении
const position = ref({ x: 0, y: 0 });

// Состояние перетаскивания
const isDragging = ref(false);
const dragStart = ref({ x: 0, y: 0 });

// Текущее изображение
const currentImage = computed(() => {
  return props.images[currentIndex.value];
});

// Стиль контейнера в зависимости от увеличения и позиции
const containerStyle = computed(() => {
  if (!isZoomed.value) {
    return {
      transform: "scale(1) translate(0px, 0px)",
    };
  }

  return {
    transform: `scale(2) translate(${position.value.x}px, ${position.value.y}px)`,
  };
});

// Закрытие просмотрщика
const close = () => {
  emit("close");
  // Сбросить увеличение при закрытии
  isZoomed.value = false;
  position.value = { x: 0, y: 0 };
};

// Переключение между обычным и увеличенным режимом
const toggleZoom = () => {
  isZoomed.value = !isZoomed.value;

  // Сбросить позицию при уменьшении
  if (!isZoomed.value) {
    position.value = { x: 0, y: 0 };
  }
};

// Перейти к предыдущему изображению
const prevImage = () => {
  currentIndex.value =
    (currentIndex.value - 1 + props.images.length) % props.images.length;
  // Сбросить увеличение при переключении
  isZoomed.value = false;
  position.value = { x: 0, y: 0 };
};

// Перейти к следующему изображению
const nextImage = () => {
  currentIndex.value = (currentIndex.value + 1) % props.images.length;
  // Сбросить увеличение при переключении
  isZoomed.value = false;
  position.value = { x: 0, y: 0 };
};

// Функции для перетаскивания увеличенного изображения
const startDrag = (event) => {
  if (!isZoomed.value) return;

  isDragging.value = true;

  // Получаем начальные координаты в зависимости от типа события (мышь или тач)
  if (event.type.startsWith("touch")) {
    dragStart.value = {
      x: event.touches[0].clientX,
      y: event.touches[0].clientY,
    };
  } else {
    dragStart.value = {
      x: event.clientX,
      y: event.clientY,
    };
  }
};

const drag = (event) => {
  if (!isDragging.value) return;

  // Предотвращаем скролл страницы при перетаскивании
  event.preventDefault();

  let currentX, currentY;

  // Получаем текущие координаты в зависимости от типа события
  if (event.type.startsWith("touch")) {
    currentX = event.touches[0].clientX;
    currentY = event.touches[0].clientY;
  } else {
    currentX = event.clientX;
    currentY = event.clientY;
  }

  // Вычисляем разницу с начальными координатами
  const deltaX = (currentX - dragStart.value.x) / 2; // Делим на 2 для более плавного движения
  const deltaY = (currentY - dragStart.value.y) / 2;

  // Обновляем позицию и начальные координаты
  position.value = {
    x: position.value.x + deltaX / 20, // Ограничиваем скорость перемещения
    y: position.value.y + deltaY / 20,
  };

  // Ограничиваем перемещение, чтобы не выйти далеко за пределы изображения
  position.value.x = Math.max(Math.min(position.value.x, 2), -2);
  position.value.y = Math.max(Math.min(position.value.y, 2), -2);

  // Обновляем начальную позицию для следующего движения
  dragStart.value = { x: currentX, y: currentY };
};

const endDrag = () => {
  isDragging.value = false;
};

// Сброс индекса при изменении начального индекса
watch(
  () => props.initialIndex,
  (newVal) => {
    currentIndex.value = newVal;
    isZoomed.value = false;
    position.value = { x: 0, y: 0 };
  }
);

// Сброс увеличения при закрытии модального окна
watch(
  () => props.isOpen,
  (newVal) => {
    if (!newVal) {
      isZoomed.value = false;
      position.value = { x: 0, y: 0 };
    }
  }
);
</script>

<style scoped>
/* Анимация для модального окна */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
