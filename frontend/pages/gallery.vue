<template>
  <div>
    <LoaderView v-if="isLoading" />
    <div v-else class="tw-container tw-mx-auto tw-px-4">
      <div class="tw-grid tw-grid-cols-1 lg:tw-grid-cols-2 tw-gap-2">
        <div
          v-for="image in galleryItems"
          :key="image.id"
          class="fade-in-element tw-cursor-pointer tw-overflow-hidden tw-transition-all tw-duration-300 hover:tw-shadow-lg"
          @click="openLightbox(image.id)"
        >
          <div class="tw-relative tw-pb-[75%]">
            <img
              :src="image.full"
              :alt="image.title"
              class="tw-absolute tw-inset-0 tw-w-full tw-h-full tw-object-cover tw-transition-transform tw-duration-500 hover:tw-scale-105"
            />
          </div>
        </div>
      </div>
    </div>

    <ImageViewer
      :images="galleryItems.map((image) => image.full)"
      :is-open="lightboxVisible"
      :initial-index="currentImageIndex"
      @close="closeLightbox"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, nextTick, watch } from "vue";
import ImageViewer from "~/components/card/ImageViewer.vue";
import { useApiService } from "~/services/api";

const { getGalleryItems, isLoading } = useApiService();

// Используем встроенный композабл для i18n
//@ts-expect-error no type for usei18n
const { locale } = useI18n();

// Состояние галереи
const galleryItems = ref([]);
const lightboxVisible = ref(false);
const currentLightboxId = ref(null);

// Переменная для хранения observer
let observer = null;

// Индекс текущего изображения для лайтбокса
const currentImageIndex = computed(() => {
  if (!currentLightboxId.value) return 0;
  const index = galleryItems.value.findIndex(
    (item) => item.id === currentLightboxId.value
  );
  return index > -1 ? index : 0;
});

// Загрузка данных галереи
const loadGalleryItems = async (categoryId = null) => {
  try {
    const response = await getGalleryItems(categoryId, locale.value);
    if (response.success && response.data) {
      galleryItems.value = response.data;
      
      // Инициализируем анимации после загрузки данных
      await nextTick();
      initScrollAnimations();
    } else {
      console.error("Ошибка загрузки галереи:", response.error);
    }
  } catch (err) {
    console.error("Ошибка при загрузке галереи:", err);
  }
};

// Функция инициализации анимаций скролла
const initScrollAnimations = () => {
  // Если observer уже существует, отключаем его
  if (observer) {
    observer.disconnect();
  }

  // Создаем новый Intersection Observer
  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          // Добавляем класс для анимации
          entry.target.classList.add("visible");
          // Прекращаем наблюдение за этим элементом
          observer.unobserve(entry.target);
        }
      });
    },
    {
      root: null, // используем viewport в качестве контейнера
      threshold: 0.1, // когда 10% элемента видно
      rootMargin: "0px 0px -50px 0px", // смещение
    }
  );

  // Находим все элементы с классом fade-in-element и начинаем наблюдение
  const elements = document.querySelectorAll(".fade-in-element");
  elements.forEach((el) => {
    observer.observe(el);
  });
};

// Открытие лайтбокса
const openLightbox = (imageId) => {
  currentLightboxId.value = imageId;
  lightboxVisible.value = true;
  document.body.style.overflow = "hidden";
};

// Закрытие лайтбокса
const closeLightbox = () => {
  lightboxVisible.value = false;
  document.body.style.overflow = "";
};

// Отслеживание изменений в galleryItems для переинициализации анимаций
watch(galleryItems, async () => {
  if (galleryItems.value.length > 0) {
    await nextTick();
    initScrollAnimations();
  }
}, { deep: true });

onMounted(async () => {
  // Загружаем данные галереи
  await loadGalleryItems();
});

onUnmounted(() => {
  // Очищаем observer при размонтировании
  if (observer) {
    observer.disconnect();
  }
  
  // Восстанавливаем скролл
  document.body.style.overflow = "";
});
</script>

<style scoped>
/* Анимация для лайтбокса */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Стили для анимации появления при скролле */
.fade-in-element {
  opacity: 0;
  transform: translateY(20px);
  transition: opacity 0.8s ease, transform 0.8s ease;
}

.fade-in-element.visible {
  opacity: 1;
  transform: translateY(0);
}
</style>