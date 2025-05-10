<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import LeftArrowIcon from "~/components/icons/LeftArrowIcon.vue";
import RightArrowIcon from "~/components/icons/RightArrowIcon.vue";
import { useNuxtApp } from "nuxt/app";
import { useApiService } from "~/services/api";
import type { GalleryItem } from "~/components";
import { useWindowSize } from "@vueuse/core";
import ContactFormHandler from "~/components/contacts/ContactFormHandler.vue";
import LoaderView from "~/components/LoaderView.vue";
import type { GsapMethods } from "~/plugins/gsap";

declare module "nuxt/app" {
  interface NuxtApp {
    $gsap: GsapMethods;
  }
}

const { $gsap } = useNuxtApp();
// Используем встроенный композабл для i18n
//@ts-expect-error no type for usei18n
const { locale } = useI18n();
const { width: windowWidth } = useWindowSize();
// API сервис для загрузки данных галереи
const { getGalleryItems, isLoading } = useApiService();

// Ссылки на DOM-элементы для анимаций
const bannerSection = ref(null);
const bannerImage = ref(null);
const aboutSection = ref(null);
const aboutTitle = ref(null);
const aboutParagraphs = ref(null);
const aboutVideo = ref(null);
const sliderSection = ref(null);
const contactSection = ref(null);

// Данные для слайдера работ из API
const galleryItems = ref<GalleryItem[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);

// Загрузка данных галереи
const loadGalleryItems = async () => {
  loading.value = true;
  error.value = null;

  try {
    const response = await getGalleryItems(null, locale.value);

    if (response.success && response.data) {
      galleryItems.value = response.data;
      console.log("Загружено элементов галереи:", galleryItems.value.length);
    } else {
      error.value = response.error || "Не удалось загрузить элементы галереи";
      console.error("Ошибка загрузки галереи:", response.error);
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Неизвестная ошибка";
    console.error("Ошибка при загрузке галереи:", err);
  } finally {
    loading.value = false;
  }
};

// Текущий слайд в основном слайдере
const currentSlide = ref(0);

// Расчет дублированного массива элементов для бесконечного слайдера
const duplicatedItems = computed(() => {
  if (galleryItems.value.length === 0) return [];
  // Дублируем элементы для создания "бесконечного" эффекта
  // Добавляем больше дубликатов для обоих слайдеров, чтобы избежать пустых мест
  return [
    ...galleryItems.value,
    ...galleryItems.value,
    ...galleryItems.value.slice(0, 3),
  ];
});

// Вычисляем смещение для анимации основного слайдера
const slideOffset = computed(() => {
  // Используем ширину слайда в зависимости от размера экрана
  // И учитываем отступы между слайдами
  const slideWidth = windowWidth.value >= 1024 ? 50 : 100;
  const slideGap = 6; // соответствует tw-gap-6 (в процентах от ширины)

  // Учитываем отступы в расчете общего смещения
  return currentSlide.value * (slideWidth + slideGap);
});

// Автоматическое переключение слайдов
let autoSlideInterval: ReturnType<typeof setInterval> | null = null;

const startAutoSlide = () => {
  if (autoSlideInterval) clearInterval(autoSlideInterval);

  autoSlideInterval = setInterval(() => {
    nextSlide();
  }, 3000); // Переключение каждые 3 секунд
};

const stopAutoSlide = () => {
  if (autoSlideInterval) {
    clearInterval(autoSlideInterval);
    autoSlideInterval = null;
  }
};

// Функция для перехода к предыдущему слайду
const prevSlide = () => {
  // Останавливаем автоматическое переключение при ручном управлении
  stopAutoSlide();

  // Если мы на первом слайде, переходим к последнему
  if (currentSlide.value === 0) {
    // Сначала устанавливаем текущий слайд на "клон" последнего элемента
    currentSlide.value = galleryItems.value.length;
    // Затем после небольшой задержки (для анимации) переходим к реальному последнему элементу
    setTimeout(() => {
      // Отключаем анимацию для мгновенного перехода
      document.querySelector(".main-slider")?.classList.add("no-transition");
      currentSlide.value = galleryItems.value.length - 1;
      // Включаем анимацию обратно
      setTimeout(() => {
        document
          .querySelector(".main-slider")
          ?.classList.remove("no-transition");
      }, 50);
    }, 500);
  } else {
    currentSlide.value--;
  }

  // Перезапускаем автоматическое переключение после ручного взаимодействия
  startAutoSlide();
};

// Функция для перехода к следующему слайду
const nextSlide = () => {
  // Если мы достигли конца массива с дубликатами
  if (currentSlide.value >= galleryItems.value.length) {
    // Сначала показываем первый дубликат для плавного перехода
    currentSlide.value++;

    // Затем после завершения анимации перехода сбрасываем на начало без анимации
    setTimeout(() => {
      // Отключаем анимацию для мгновенного перехода
      document.querySelector(".main-slider")?.classList.add("no-transition");
      currentSlide.value = 0;
      // Включаем анимацию обратно
      setTimeout(() => {
        document
          .querySelector(".main-slider")
          ?.classList.remove("no-transition");
      }, 50);
    }, 500);
  } else {
    // Обычный переход к следующему слайду
    currentSlide.value++;
  }
};

// Текущий слайд во втором слайдере (движется в противоположном направлении)
const currentReverseSlide = ref(0);

// Вычисляем смещение для анимации второго слайдера
const reverseSlideOffset = computed(() => {
  // Для обратного слайдера используем меньшие слайды
  const slideWidth = windowWidth.value >= 1024 ? 25 : 50;
  const slideGap = 2; // соответствует tw-gap-2 (в процентах от ширины)

  // Отрицательное значение смещения с учетом отступов
  return -currentReverseSlide.value * (slideWidth + slideGap);
});

// Автоматическое переключение слайдов для второго слайдера
let autoReverseSlideInterval: ReturnType<typeof setInterval> | null = null;

const startAutoReverseSlide = () => {
  if (autoReverseSlideInterval) clearInterval(autoReverseSlideInterval);

  autoReverseSlideInterval = setInterval(() => {
    // Обратное направление
    nextReverseSlide();
  }, 2000); // Немного быстрее основного слайдера
};

const stopAutoReverseSlide = () => {
  if (autoReverseSlideInterval) {
    clearInterval(autoReverseSlideInterval);
    autoReverseSlideInterval = null;
  }
};

// Функция для перехода к предыдущему слайду (для второго слайдера)
const prevReverseSlide = () => {
  // Останавливаем автоматическое переключение при ручном управлении
  stopAutoReverseSlide();

  // Если мы на первом слайде, переходим к последнему
  if (currentReverseSlide.value === 0) {
    // Сначала устанавливаем текущий слайд на "клон" последнего элемента
    currentReverseSlide.value = galleryItems.value.length;
    // Затем после небольшой задержки (для анимации) переходим к реальному последнему элементу
    setTimeout(() => {
      // Отключаем анимацию для мгновенного перехода
      document.querySelector(".reverse-slider")?.classList.add("no-transition");
      currentReverseSlide.value = galleryItems.value.length - 1;
      // Включаем анимацию обратно
      setTimeout(() => {
        document
          .querySelector(".reverse-slider")
          ?.classList.remove("no-transition");
      }, 50);
    }, 500);
  } else {
    currentReverseSlide.value--;
  }

  // Перезапускаем автоматическое переключение после ручного взаимодействия
  startAutoReverseSlide();
};

// Функция для перехода к следующему слайду (для второго слайдера)
const nextReverseSlide = () => {
  // Останавливаем автоматическое переключение при ручном управлении
  stopAutoReverseSlide();

  // Если мы достигли конца массива с дубликатами
  if (currentReverseSlide.value >= galleryItems.value.length) {
    // Сначала показываем первый дубликат для плавного перехода
    currentReverseSlide.value++;

    // Затем после завершения анимации перехода сбрасываем на начало без анимации
    setTimeout(() => {
      // Отключаем анимацию для мгновенного перехода
      document.querySelector(".reverse-slider")?.classList.add("no-transition");
      currentReverseSlide.value = 0;
      // Включаем анимацию обратно
      setTimeout(() => {
        document
          .querySelector(".reverse-slider")
          ?.classList.remove("no-transition");
      }, 50);
    }, 500);
  } else {
    // Обычный переход к следующему слайду
    currentReverseSlide.value++;
  }

  // Перезапускаем автоматическое переключение после ручного взаимодействия
  startAutoReverseSlide();
};

// Добавим переменную для хранения инстанса анимации
let videoAnimation = null;

// Функция для создания или обновления анимации видео
const updateVideoAnimation = async () => {
  if (!aboutSection.value || !aboutVideo.value) return;

  try {
    const { gsap } = await $gsap.getGSAP();

    // Очищаем предыдущую анимацию, если она существует
    if (videoAnimation) {
      videoAnimation.kill();
      videoAnimation = null;
    }

    // Создаем анимацию только если ширина экрана >= 1024px
    if (windowWidth.value >= 1024) {
      // Анимация смещения видео при скролле (эффект параллакса)
      videoAnimation = gsap.to(aboutVideo.value, {
        xPercent: -15,
        ease: "none",
        scrollTrigger: {
          trigger: aboutSection.value,
          start: "top bottom",
          end: "bottom top",
          scrub: true,
        },
      });

      console.log("Анимация видео активирована для Desktop");
    } else {
      console.log("Анимация видео отключена для мобильных устройств");
    }
  } catch (error) {
    console.error("Ошибка инициализации анимации видео:", error);
  }
};

// Функция для проверки загрузки видео и его воспроизведения
const handleVideoLoad = () => {
  if (aboutVideo.value) {
    // Правильно обращаемся к DOM-элементу video, находящемуся внутри aboutVideo.value
    const videoElement = aboutVideo.value.querySelector("video");
    if (videoElement) {
      // Проверяем поддержку autoplay
      const promise = videoElement.play();

      if (promise !== undefined) {
        promise.catch((error) => {
          console.log("Автовоспроизведение видео не поддерживается:", error);
          // Если автовоспроизведение не поддерживается, показываем элементы управления
          videoElement.controls = true;
        });
      }
    }
  }
};

// Наблюдаем за изменением ширины окна
watch(
  windowWidth,
  () => {
    updateVideoAnimation();
  },
  { immediate: false }
);

// В функции onMounted заменяем вызов initVideoAnimation на updateVideoAnimation
onMounted(async () => {
  // Загружаем данные галереи
  await loadGalleryItems();

  if (bannerSection.value && bannerImage.value) {
    // Инициализируем параллакс для баннера
    await $gsap.initParallaxBanner(bannerSection.value, bannerImage.value);
  }

  if (aboutSection.value) {
    // Анимация заголовка "О нас"
    await $gsap.initScrollReveal(aboutTitle.value, {
      y: 30,
      duration: 0.7,
      scrollTrigger: {
        trigger: aboutSection.value,
        start: "top 80%",
      },
    });

    // Анимация параграфов в секции "О нас"
    if (aboutParagraphs.value) {
      await $gsap.initStaggerAnimation(
        aboutParagraphs.value,
        aboutSection.value,
        {
          y: 30,
          delay: 0.2,
          stagger: 0.2,
        }
      );
    }

    // Инициализируем анимацию для видео с учетом ширины экрана
    await updateVideoAnimation();

    // Обработка загрузки видео
    handleVideoLoad();
  }

  // Запускаем анимацию для блока категорий
  const categoryCards = document.querySelectorAll(".category-card");
  if (categoryCards.length > 0) {
    const { gsap } = await $gsap.getGSAP();

    // Анимация для основных категорий (появление сверху вниз с задержкой)
    gsap.from(".category-card", {
      y: 50,
      opacity: 0,
      duration: 0.8,
      stagger: 0.2, // Интервал между появлением каждой карточки
      ease: "power2.out",
      scrollTrigger: {
        trigger: ".category-card",
        start: "top 85%",
      },
    });
  }

  // Анимация секции слайдера
  if (sliderSection.value) {
    const { gsap } = await $gsap.getGSAP();

    // Анимация заголовка и подзаголовка слайдера
    gsap.from(sliderSection.value.querySelectorAll("h2, p"), {
      y: 30,
      opacity: 0,
      duration: 0.7,
      stagger: 0.2,
      scrollTrigger: {
        trigger: sliderSection.value,
        start: "top 80%",
      },
    });

    // Анимация основного слайдера
    gsap.from(".main-slider > div", {
      scale: 0.9,
      opacity: 0,
      duration: 1,
      stagger: 0.1,
      delay: 0.3,
      scrollTrigger: {
        trigger: ".main-slider",
        start: "top 80%",
      },
    });

    // Анимация слайдера миниатюр
    gsap.from(".reverse-slider > div", {
      scale: 0.9,
      opacity: 0,
      duration: 0.8,
      stagger: 0.05,
      delay: 0.5,
      scrollTrigger: {
        trigger: ".reverse-slider",
        start: "top 90%",
      },
    });
  }

  // Анимация секции контактов
  if (contactSection.value) {
    await $gsap.initScrollReveal(contactSection.value, {
      opacity: 0,
      y: 40,
      duration: 0.8,
    });
  }

  // Запускаем автоматическое переключение слайдов
  startAutoSlide();
  startAutoReverseSlide();
});

// В функции onUnmounted добавляем очистку анимации и интервалов
onUnmounted(() => {
  // Очищаем анимацию видео
  if (videoAnimation) {
    videoAnimation.kill();
    videoAnimation = null;
  }

  // Останавливаем автоматическое переключение слайдов
  stopAutoSlide();
  stopAutoReverseSlide();

  // Останавливаем видео
  if (aboutVideo.value) {
    const videoElement = aboutVideo.value.querySelector("video");
    if (videoElement) {
      videoElement.pause();
      videoElement.removeAttribute("src");
      videoElement.load();
    }
  }
});
</script>

<template>
  <LoaderView v-if="isLoading" />
  <div v-else>
    <!-- Главный постер с параллаксом -->
    <section
      ref="bannerSection"
      class="tw-relative tw-flex tw-items-center tw-justify-center tw-overflow-hidden"
      :style="{ height: 'calc(100vh - 128px)' }"
    >
      <div class="tw-absolute tw-inset-0 tw-z-0">
        <div class="tw-w-full tw-h-full tw-bg-gray-100">
          <img
            ref="bannerImage"
            src="https://s3.stroi-news.ru/img/masterskaya-kartinki-1.jpg"
            :alt="$t('home.title')"
            class="tw-w-full tw-h-full tw-object-cover"
          />
          <!-- Затемняющий градиент -->
          <div
            class="tw-absolute tw-inset-0 tw-bg-gradient-to-t tw-from-black/70 tw-via-black/40 tw-to-black/20"
          ></div>
        </div>
      </div>

      <div
        class="tw-container tw-mx-auto tw-px-4 tw-relative tw-z-10 tw-text-center"
      >
        <h1
          class="tw-text-4xl md:tw-text-6xl tw-font-bold tw-text-white tw-mb-4 tw-shadow-text tw-uppercase"
        >
          {{ $t("home.title") }}
        </h1>
        <p
          class="tw-text-xl md:tw-text-2xl tw-text-white tw-max-w-2xl tw-mx-auto tw-shadow-text"
        >
          {{ $t("home.subtitle") }}
        </p>
        <div class="tw-mt-8">
          <NuxtLink
            to="/catalog"
            class="tw-bg-white tw-uppercase tw-text-gray-800 tw-px-8 tw-py-3 tw-rounded-md tw-text-lg tw-font-medium tw-inline-block tw-shadow-md tw-transition-all tw-duration-300 hover:tw-shadow-lg hover:tw-opacity-90"
          >
            {{ $t("home.go_to_catalog") }}
          </NuxtLink>
        </div>
      </div>
    </section>

    <!-- Модернизированный блок "О нас" с видео справа -->
    <section
      ref="aboutSection"
      class="tw-relative tw-py-16 tw-overflow-hidden"
      id="about"
    >
      <!-- Градиентный фон слева -->
      <div class="tw-absolute tw-inset-0 tw-z-0">
        <div
          class="tw-absolute tw-inset-0 tw-bg-gradient-to-r tw-from-white tw-via-white/95 tw-to-transparent"
        ></div>
      </div>

      <!-- Видео справа (фоновое) -->
      <div
        class="tw-absolute tw-right-0 tw-bottom-0 md:tw-top-0 lg:tw-bottom-unset tw-w-full md:tw-w-[85%] md:tw-mr-[-10%] tw-z-0 tw-overflow-hidden tw-h-[70%] md:tw-h-full"
      >
        <!-- Градиентный оверлей -->
        <div
          class="gradient-horizontal tw-pointer-events-none tw-absolute tw-inset-0 tw-z-10 tw-hidden lg:tw-block"
        ></div>
        <div ref="aboutVideo" class="tw-relative tw-w-full tw-h-full">
          <!-- Само видео -->
          <video
            class="tw-object-cover tw-w-full tw-h-full"
            autoplay
            loop
            muted
            playsinline
          >
            <source src="../assets/img/banner_video.mp4" type="video/mp4" />
            <img
              src="../assets/img/engraving3.png"
              alt="Мастерская"
              class="tw-object-cover tw-w-full tw-h-full"
            />
          </video>

          <!-- Градиентный оверлей -->
          <div
            class="tw-pointer-events-none tw-absolute tw-inset-0 tw-z-10"
            :style="{
              background: `linear-gradient(
                to ${windowWidth > 768 ? 'right' : 'top'},
                rgba(255, 255, 255, 1) 0%,
                rgba(255, 255, 255, 0) 20%,
                rgba(255, 255, 255, 0) 80%,
                rgba(255, 255, 255, 1) 100%
              )`,
            }"
          ></div>
        </div>
      </div>

      <!-- Контент блока -->
      <div class="tw-container tw-mx-auto tw-px-4 tw-relative tw-z-10">
        <div class="tw-flex tw-flex-col md:tw-flex-row tw-items-center">
          <!-- Текстовая часть (слева) -->
          <div class="md:tw-w-1/2 tw-py-8 md:tw-py-16">
            <h2
              ref="aboutTitle"
              class="tw-text-3xl tw-font-bold tw-text-gray-800 tw-mb-8"
            >
              {{ $t("home.about_us") }}
            </h2>

            <div class="tw-max-w-md">
              <p
                ref="aboutParagraphs"
                class="tw-text-lg tw-text-gray-600 tw-mb-6"
              >
                {{ $t("home.about_text1") }}
              </p>
              <p
                ref="aboutParagraphs"
                class="tw-text-lg tw-text-gray-600 tw-mb-6"
              >
                {{ $t("home.about_text2") }}
              </p>
              <p ref="aboutParagraphs" class="tw-text-lg tw-text-gray-600">
                {{ $t("home.about_text3") }}
              </p>

              <!-- Кнопка перехода в каталог -->
              <div class="tw-mt-8">
                <NuxtLink
                  to="/catalog"
                  class="tw-bg-gray-800 tw-uppercase tw-text-white tw-px-6 tw-py-3 tw-rounded-md tw-text-base tw-font-medium tw-inline-block tw-shadow-md tw-transition-all tw-duration-300 hover:tw-shadow-lg hover:tw-bg-gray-700"
                >
                  {{ $t("home.discover_products") }}
                </NuxtLink>
              </div>
            </div>
          </div>

          <!-- Пустое пространство для видео (на мобильных) -->
          <div
            class="md:tw-w-1/2 tw-h-64 md:tw-h-auto tw-relative tw-z-0 md:tw-hidden"
          >
            <!-- Это пространство будет заполнено фоновым видео на мобильных устройствах -->
          </div>
        </div>
      </div>
    </section>

    <!-- Плитка категорий -->
    <section class="tw-bg-white">
      <div class="tw-mx-auto">
        <!-- Основные категории (2 в ряду) -->
        <div class="tw-grid tw-grid-cols-1 md:tw-grid-cols-2 tw-gap-4 tw-mb-4">
          <!-- 3D печать -->
          <NuxtLink
            to="/catalog?category=3d-printing"
            class="category-card tw-relative tw-overflow-hidden tw-shadow-md tw-aspect-w-16 tw-aspect-h-9 tw-block"
          >
            <img
              src="../assets/img/3d_printing.png"
              alt="3D Печать"
              class="tw-w-full tw-h-full tw-object-cover tw-transition-transform tw-duration-500"
            />
            <div
              class="tw-absolute tw-inset-0 tw-bg-gradient-to-t tw-from-black/70 tw-to-transparent"
            ></div>
            <div
              class="tw-absolute tw-bottom-0 tw-left-0 tw-w-full tw-p-6 tw-flex tw-flex-col tw-items-center tw-justify-center"
            >
              <h3 class="tw-text-6xl tw-font-bold tw-text-white">
                {{ $t("catalog.categories.3d-printing") }}
              </h3>
              <p class="tw-text-gray-200 tw-mt-2">
                {{ $t("home.categories.3d_description") }}
              </p>
            </div>
          </NuxtLink>

          <!-- Выжигание -->
          <NuxtLink
            to="/catalog?category=engraving"
            class="category-card tw-relative tw-overflow-hidden tw-shadow-md tw-aspect-w-16 tw-aspect-h-9 tw-block"
          >
            <img
              src="../assets/img/engraving.png"
              alt="Выжигание"
              class="tw-w-full tw-h-full tw-object-cover tw-transition-transform tw-duration-500"
            />
            <div
              class="tw-absolute tw-inset-0 tw-bg-gradient-to-t tw-from-black/70 tw-to-transparent"
            ></div>
            <div
              class="tw-absolute tw-bottom-0 tw-left-0 tw-w-full tw-p-6 tw-flex tw-flex-col tw-items-center tw-justify-center"
            >
              <h3 class="tw-text-6xl tw-font-bold tw-text-white">
                {{ $t("catalog.categories.engraving") }}
              </h3>
              <p class="tw-text-gray-200 tw-mt-2">
                {{ $t("home.categories.engraving_description") }}
              </p>
            </div>
          </NuxtLink>
        </div>

        <!-- Подкатегории (3 в ряду) -->
        <div
          class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 md:tw-grid-cols-3 tw-gap-4"
        >
          <!-- Дерево -->
          <NuxtLink
            to="/catalog?category=engraving&subcategory=wood"
            class="category-card tw-relative tw-overflow-hidden tw-shadow-md tw-aspect-w-4 tw-aspect-h-3 tw-block"
          >
            <img
              src="../assets/img/wood.png"
              alt="Дерево"
              class="tw-w-full tw-h-full tw-object-cover tw-transition-transform tw-duration-500"
            />
            <div
              class="tw-absolute tw-inset-0 tw-bg-gradient-to-t tw-from-black/70 tw-to-transparent"
            ></div>
            <div
              class="tw-absolute tw-bottom-0 tw-left-0 tw-w-full tw-p-6 tw-flex tw-flex-col tw-items-center tw-justify-center"
            >
              <h3 class="tw-text-4xl tw-font-bold tw-text-white">
                {{ $t("catalog.subcategories.wood") }}
              </h3>
              <p class="tw-text-gray-200 tw-mt-2">
                {{ $t("home.categories.wood_description") }}
              </p>
            </div>
          </NuxtLink>

          <!-- Металл -->
          <NuxtLink
            to="/catalog?category=engraving&subcategory=metal"
            class="category-card tw-relative tw-overflow-hidden tw-shadow-md tw-aspect-w-4 tw-aspect-h-3 tw-block"
          >
            <img
              src="../assets/img/metal.png"
              alt="Металл"
              class="tw-w-full tw-h-full tw-object-cover tw-transition-transform tw-duration-500"
            />
            <div
              class="tw-absolute tw-inset-0 tw-bg-gradient-to-t tw-from-black/70 tw-to-transparent"
            ></div>
            <div
              class="tw-absolute tw-bottom-0 tw-left-0 tw-w-full tw-flex tw-p-6 tw-flex-col tw-items-center tw-justify-center"
            >
              <h3 class="tw-text-4xl tw-font-bold tw-text-white">
                {{ $t("catalog.subcategories.metal") }}
              </h3>
              <p class="tw-text-gray-200 tw-mt-2">
                {{ $t("home.categories.metal_description") }}
              </p>
            </div>
          </NuxtLink>

          <!-- Другое -->
          <NuxtLink
            to="/catalog?category=engraving&subcategory=other"
            class="category-card tw-relative tw-overflow-hidden tw-shadow-md tw-aspect-w-4 tw-aspect-h-3 tw-block"
          >
            <img
              src="../assets/img/glass.png"
              alt="Другое"
              class="tw-w-full tw-h-full tw-object-cover tw-transition-transform tw-duration-500"
            />
            <div
              class="tw-absolute tw-inset-0 tw-bg-gradient-to-t tw-from-black/70 tw-to-transparent"
            ></div>
            <div
              class="tw-absolute tw-bottom-0 tw-left-0 tw-w-full tw-flex tw-p-6 tw-flex-col tw-items-center tw-justify-center"
            >
              <h3 class="tw-text-4xl tw-font-bold tw-text-white">
                {{ $t("catalog.subcategories.other") }}
              </h3>
              <p class="tw-text-gray-200 tw-mt-2">
                {{ $t("home.categories.other_description") }}
              </p>
            </div>
          </NuxtLink>
        </div>
      </div>
    </section>

    <!-- Улучшенный слайдер с примерами работ - основная часть -->
    <section
      ref="sliderSection"
      class="tw-relative tw-mx-auto tw-overflow-hidden tw-pt-6 tw-pb-8 tw-bg-gray-100"
    >
      <div v-if="loading" class="tw-flex tw-justify-center tw-py-16">
        <div
          class="tw-animate-spin tw-rounded-full tw-h-12 tw-w-12 tw-border-t-2 tw-border-b-2 tw-border-gray-900"
        ></div>
      </div>

      <template v-else>
        <!-- Основной слайдер -->
        <div class="tw-overflow-hidden tw-mb-2 tw-relative">
          <div
            class="main-slider tw-transition-all tw-duration-500 tw-ease-in-out tw-flex tw-gap-6"
            :style="`transform: translateX(-${slideOffset}%);`"
            @mouseenter="stopAutoSlide"
            @mouseleave="startAutoSlide"
          >
            <div
              v-for="(item, index) in duplicatedItems"
              :key="`main-${item.id}-${index}`"
              class="tw-w-full lg:tw-w-1/2 tw-flex-shrink-0"
            >
              <div
                class="tw-relative tw-overflow-hidden tw-h-[700px] tw-shadow-lg"
              >
                <img
                  :src="item.full"
                  :alt="item.title"
                  class="tw-w-full tw-h-full tw-object-cover"
                />
                <div
                  class="tw-absolute tw-inset-0 tw-bg-gradient-to-t tw-from-black/80 tw-via-black/40 tw-to-transparent"
                ></div>
                <div class="tw-absolute tw-bottom-0 tw-left-0 tw-w-full tw-p-6">
                  <span
                    class="tw-text-sm tw-font-medium tw-text-white tw-bg-black/40 tw-py-1 tw-px-3 tw-rounded-full"
                  >
                    {{ item.category }}
                  </span>
                  <h3 class="tw-text-2xl tw-font-bold tw-text-white tw-mt-2">
                    {{ item.title }}
                  </h3>
                  <p class="tw-text-gray-200 tw-line-clamp-2 tw-mt-2">
                    {{ item.description }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Кнопки навигации - возвращены в исходный стиль и положение -->
          <button
            @click="prevSlide"
            class="tw-absolute tw-left-4 tw-top-1/2 tw-transform -tw-translate-y-1/2 tw-p-2 tw-rounded-full tw-text-white tw-transition-all tw-duration-300 hover:tw-text-gray-200 focus:tw-outline-none"
            :aria-label="$t('slider.prev_slide')"
          >
            <LeftArrowIcon class="tw-w-10 tw-h-10" />
          </button>
          <button
            @click="nextSlide"
            class="tw-absolute tw-right-4 tw-top-1/2 tw-transform -tw-translate-y-1/2 tw-p-2 tw-rounded-full tw-text-white tw-transition-all tw-duration-300 hover:tw-text-gray-200 focus:tw-outline-none"
            :aria-label="$t('slider.next_slide')"
          >
            <RightArrowIcon class="tw-w-10 tw-h-10" />
          </button>

          <!-- Индикаторы позиции слайдера -->
          <div class="tw-flex tw-justify-center tw-gap-2 tw-mt-2">
            <div class="tw-flex tw-items-center tw-gap-2">
              <span
                v-for="(_, i) in galleryItems"
                :key="`dot-${i}`"
                class="tw-w-2 tw-h-2 tw-rounded-full tw-transition-all tw-duration-300"
                :class="
                  i === currentSlide % galleryItems.length
                    ? 'tw-bg-gray-800 tw-w-4'
                    : 'tw-bg-gray-400'
                "
                @click="currentSlide = i"
              ></span>
            </div>
          </div>
        </div>

        <!-- Второй слайдер (в обратном направлении) - миниатюры -->
        <div class="tw-overflow-hidden tw-relative">
          <div
            class="reverse-slider tw-transition-all tw-duration-500 tw-ease-in-out tw-flex tw-gap-2"
            :style="`transform: translateX(${reverseSlideOffset}%);`"
            @mouseenter="stopAutoReverseSlide"
            @mouseleave="startAutoReverseSlide"
          >
            <div
              v-for="(item, index) in duplicatedItems"
              :key="`reverse-${item.id}-${index}`"
              class="tw-w-1/2 md:tw-w-1/3 lg:tw-w-1/4 tw-flex-shrink-0"
            >
              <div
                class="tw-relative tw-overflow-hidden tw-h-[300px] tw-shadow-md tw-cursor-pointer tw-transition-all tw-duration-300 hover:tw-shadow-lg"
                @click="currentSlide = index % galleryItems.length"
              >
                <img
                  :src="item.thumbnail || item.full"
                  :alt="item.title"
                  class="tw-w-full tw-h-full tw-object-cover"
                />
                <div
                  class="tw-absolute tw-inset-0 tw-bg-gradient-to-t tw-from-black/60 tw-to-transparent"
                ></div>
                <div class="tw-absolute tw-bottom-0 tw-left-0 tw-w-full tw-p-3">
                  <h4
                    class="tw-text-sm tw-font-medium tw-text-white tw-line-clamp-1"
                  >
                    {{ item.title }}
                  </h4>
                </div>
              </div>
            </div>
          </div>

          <!-- Кнопки навигации для второго слайдера -->
          <button
            @click="nextReverseSlide"
            class="tw-absolute tw-left-4 tw-top-1/2 tw-transform -tw-translate-y-1/2 tw-p-2 tw-rounded-full tw-text-white tw-transition-all tw-duration-300 hover:tw-text-gray-200 focus:tw-outline-none"
            :aria-label="$t('slider.prev_slide')"
          >
            <LeftArrowIcon class="tw-w-8 tw-h-8" />
          </button>
          <button
            @click="prevReverseSlide"
            class="tw-absolute tw-right-4 tw-top-1/2 tw-transform -tw-translate-y-1/2 tw-p-2 tw-rounded-full tw-text-white tw-transition-all tw-duration-300 hover:tw-text-gray-200 focus:tw-outline-none"
            :aria-label="$t('slider.next_slide')"
          >
            <RightArrowIcon class="tw-w-8 tw-h-8" />
          </button>

          <!-- Индикаторы позиции для второго слайдера -->
          <div class="tw-flex tw-justify-center tw-gap-2 tw-mt-2">
            <div class="tw-flex tw-items-center tw-gap-2">
              <span
                v-for="(_, i) in galleryItems"
                :key="`reverse-dot-${i}`"
                class="tw-w-2 tw-h-2 tw-rounded-full tw-transition-all tw-duration-300"
                :class="
                  i === currentReverseSlide % galleryItems.length
                    ? 'tw-bg-gray-800 tw-w-4'
                    : 'tw-bg-gray-400'
                "
                @click="currentReverseSlide = i"
              ></span>
            </div>
          </div>
        </div>

        <!-- Ссылка на галерею -->
        <div class="tw-flex tw-justify-center tw-mt-10">
          <NuxtLink
            to="/gallery"
            class="tw-bg-gray-800 tw-uppercase tw-text-white tw-px-6 tw-py-3 tw-rounded-md tw-text-base tw-font-medium tw-inline-block tw-shadow-md tw-transition-all tw-duration-300 hover:tw-shadow-lg hover:tw-bg-gray-700"
          >
            {{ $t("home.view_all_gallery") }}
          </NuxtLink>
        </div>
      </template>
    </section>

    <!-- Форма обратной связи -->
    <section ref="contactSection" class="tw-py-16 tw-bg-white">
      <ContactFormHandler />
    </section>
  </div>
</template>

<style scoped>
/* Эффект увеличения при наведении */
.category-card:hover img {
  transform: scale(1.05);
}

/* Стили для соотношения сторон */
.tw-aspect-w-16 {
  position: relative;
  padding-bottom: 70%; /* 16:9 */
}

.tw-aspect-w-4 {
  position: relative;
  padding-bottom: 50%; /* 4:3 */
}

.tw-aspect-w-16 > *,
.tw-aspect-w-4 > * {
  position: absolute;
  height: 100%;
  width: 100%;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
}

/* Стили для бесконечных слайдеров */
.main-slider,
.reverse-slider {
  transition: transform 0.5s ease-in-out;
}

/* Класс для отключения анимации при сбросе */
.no-transition {
  transition: none !important;
}

/* Подсветка активного слайда */
.active-slide {
  opacity: 1;
  transform: scale(1.05);
  z-index: 10;
}

/* Ограничение текста до 2 строк */
.tw-line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.tw-line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Анимация для точек навигации */
.dot-indicator {
  transition: all 0.3s ease;
}

/* Увеличение при наведении на миниатюры */
.reverse-slider > div:hover {
  transform: translateY(-5px);
  transition: transform 0.3s ease;
}

/* Тень для текста на постере */
.tw-shadow-text {
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
}

.gradient-horizontal {
  background: linear-gradient(
    to right,
    rgba(255, 255, 255, 1) 0%,
    rgba(255, 255, 255, 0) 20%,
    rgba(255, 255, 255, 0) 80%,
    rgba(255, 255, 255, 1) 100%
  );
}
</style>
