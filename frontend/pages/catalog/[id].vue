<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { ref, computed, onMounted, watch } from "vue";
import { useRoute } from "vue-router";
import { useCart } from "~/shared/useCart";
import ImageViewer from "~/components/card/ImageViewer.vue";
import RelatedProductsSlider from "~/components/card/RelatedProductsSlider.vue";
import CrossIcon from "~/components/icons/CrossIcon.vue";
import MinusIcon from "~/components/icons/MinusIcon.vue";
import PlusIcon from "~/components/icons/PlusIcon.vue";
import MarkIcon from "~/components/icons/MarkIcon.vue";
import { getSeoMetadata } from "~/shared/useSeo";
import { useSeoMeta } from "nuxt/app";
import { useApiService } from "~/services/api";
import { useCategories } from "~/shared/useCategories";
import { currencyMap } from "~/components";
import OrderFormHandler from "~/components/card/OrderFormHandler.vue";
import LoaderView from "~/components/LoaderView.vue";
import type { Product } from "~/components";

const route = useRoute();
const { t, locale } = useI18n();
const productId = computed(() => Number(route.params.id));

// API сервис
const { getProductById, isLoading } = useApiService();
const { addProductToCart } = useCart();
const { getCategory, getSubcategory } = useCategories();

// Данные продукта
const product = ref<Product | null>(null);
const productError = ref<string | null>(null);

// Связанные товары
const relatedProducts = ref([]);
const relatedError = ref<string | null>(null);

// SEO метаданные
const setupSeo = () => {
  if (!product.value) return;

  // Формируем специальное описание для SEO
  const seoDescription = computed(() => {
    if (!product.value) return "";
    const description =
      product.value.translations &&
      product.value.translations[locale.value]?.description
        ? product.value.translations[locale.value].description
        : product.value.description;
    const material = product.value.characteristics?.material || "";
    return `${product.value.name} - ${description.slice(
      0,
      150
    )}... Материал: ${material}. Купить с доставкой по России.`;
  });

  // Получаем метаданные для SEO с переопределенными значениями
  const seoData = getSeoMetadata(route, t, locale.value, {
    // Переопределяем заголовок
    title: `${product.value.name} | Prianik Studio`,
    // Используем специальное описание
    description: seoDescription.value,
    // Используем первое изображение товара
    image:
      product.value.images && product.value.images.length > 0
        ? product.value.images[0]
        : "",
    // Указываем тип контента
    type: "product",
  });

  // Применяем SEO-метатеги
  useSeoMeta({
    title: seoData.title,
    description: seoData.description,

    // Open Graph метатеги
    ogTitle: seoData.title,
    ogDescription: seoData.description,
    ogImage: seoData.image,
    ogUrl: seoData.pageUrl,
    //@ts-expect-error
    ogType: seoData.type,
    ogLocale: seoData.locale,

    // Twitter метатеги
    twitterTitle: seoData.title,
    twitterDescription: seoData.description,
    twitterImage: seoData.image,
    //@ts-expect-error
    twitterCard: seoData.twitterCard,

    // Дополнительные метатеги для товара
    // Для улучшения отображения в поисковых системах
    itemProp: [
      { name: "name", content: product.value.name },
      { name: "description", content: seoDescription.value },
      {
        name: "image",
        content:
          product.value.images && product.value.images.length > 0
            ? product.value.images[0]
            : "",
      },
    ],

    // Микроразметка для товара в формате JSON-LD
    script: [
      {
        type: "application/ld+json",
        innerHTML: JSON.stringify({
          "@context": "https://schema.org",
          "@type": "Product",
          name: product.value.name,
          description: product.value.description,
          image: product.value.images,
          offers: {
            "@type": "Offer",
            price: product.value.price,
            priceCurrency: "RUB",
            availability: "https://schema.org/InStock",
          },
        }),
      },
    ],
  });
};

// Форма заказа
const isOrderFormVisible = ref(false);
const currentProductToOrder = ref(null);

// Прямая покупка (открывает форму заказа)
const buyNow = () => {
  if (product.value) {
    currentProductToOrder.value = {
      id: product.value.id,
      quantity: quantity.value,
    };
    isOrderFormVisible.value = true;
  }
};

// Обработчик успешного оформления заказа
const handleOrderSuccess = () => {
  isOrderFormVisible.value = false;
  // Можно добавить дополнительные действия после успешного заказа
};

// Количество товара для заказа
const quantity = ref(1);
const increaseQuantity = () => {
  quantity.value += 1;
};
const decreaseQuantity = () => {
  if (quantity.value > 1) quantity.value -= 1;
};

// Активное изображение товара
const activeImageIndex = ref(0);
const setActiveImage = (index) => {
  activeImageIndex.value = index;
};

// Состояние модального окна для просмотра изображений
const isImageViewerOpen = ref(false);
const openImageViewer = (index) => {
  activeImageIndex.value = index;
  isImageViewerOpen.value = true;
};
const closeImageViewer = () => {
  isImageViewerOpen.value = false;
};

const isAddedToCart = ref(false);

// Добавление товара в корзину
const addToCart = () => {
  if (!product.value) return;

  // Добавляем товар в корзину через composable
  const success = addProductToCart(
    {
      id: product.value.id,
      name: product.value.name,
      price: product.value.price,
      image:
        product.value.images && product.value.images.length > 0
          ? product.value.images[0]
          : undefined,
      currency: product.value.currency,
    },
    quantity.value
  );

  // Показываем индикатор успешного добавления
  if (success) {
    isAddedToCart.value = true;

    // Возвращаем кнопку в исходное состояние через некоторое время
    setTimeout(() => {
      isAddedToCart.value = false;
    }, 3000);
  }
};

// Загрузка данных продукта
const loadProduct = async () => {
  productError.value = null;

  try {
    const response = await getProductById(productId.value, locale.value);

    if (response.success && response.data) {
      product.value = response.data;
      relatedProducts.value = response.data.related_products;
      // Сброс индекса активного изображения при загрузке нового продукта
      activeImageIndex.value = 0;

      // Настройка SEO-метаданных после загрузки данных
      setupSeo();
    } else {
      productError.value = response.error || "Не удалось загрузить товар";
      console.error("Ошибка загрузки товара:", response.error);
    }
  } catch (err) {
    productError.value =
      err instanceof Error ? err.message : "Неизвестная ошибка";
    console.error("Ошибка при загрузке товара:", err);
  }
};

// Загрузка данных при монтировании компонента
onMounted(loadProduct);

// Отслеживание изменения ID продукта для перезагрузки данных
watch(
  () => [productId.value, locale.value],

  loadProduct
);
</script>

<template>
  <div class="tw-py-6">
    <div class="tw-container tw-mx-auto tw-px-4">
      <!-- Состояние загрузки -->
      <LoaderView v-if="isLoading" />

      <!-- Ошибка загрузки -->
      <div
        v-else-if="productError"
        class="tw-bg-red-50 tw-p-8 tw-rounded-md tw-text-red-700 tw-text-center"
      >
        <p class="tw-text-xl tw-font-semibold tw-mb-2">
          {{ $t("product.load_error") }}
        </p>
        <p>{{ productError }}</p>
        <router-link
          to="/catalog"
          class="tw-inline-block tw-mt-4 tw-bg-gray-800 tw-text-white tw-py-2 tw-px-4 tw-rounded-md"
        >
          {{ $t("product.back_to_catalog") }}
        </router-link>
      </div>

      <!-- Содержимое страницы товара -->
      <div v-else-if="product">
        <!-- Хлебные крошки -->
        <div class="tw-mb-8">
          <div class="tw-flex tw-items-center tw-text-sm tw-text-gray-600">
            <NuxtLink
              to="/"
              class="hover:tw-text-gray-800 tw-transition-colors"
              >{{ $t("product.main") }}</NuxtLink
            >
            <span class="tw-mx-2">/</span>
            <NuxtLink
              to="/catalog"
              class="hover:tw-text-gray-800 tw-transition-colors"
              >{{ $t("product.catalog") }}</NuxtLink
            >
            <span class="tw-mx-2">/</span>
            <NuxtLink
              :to="`/catalog?category=${getCategory(
                product.category_id
              )?.toLowerCase()}`"
              class="hover:tw-text-gray-800 tw-transition-colors"
            >
              {{ getCategory(product.category_id) }}
            </NuxtLink>
            <span v-if="product.subcategory_id" class="tw-mx-2">/</span>
            <NuxtLink
              v-if="product.subcategory_id"
              :to="`/catalog?category=${getCategory(
                product.category_id
              )?.toLowerCase()}&subcategory=${getSubcategory(
                product.subcategory_id,
                product.category_id
              )?.toLowerCase()}`"
              class="hover:tw-text-gray-800 tw-transition-colors"
            >
              {{ getSubcategory(product.subcategory_id, product.category_id) }}
            </NuxtLink>
          </div>
        </div>

        <!-- Основная информация о товаре -->
        <div class="tw-flex tw-flex-col md:tw-flex-row tw-gap-8 tw-mb-16">
          <!-- Галерея изображений -->
          <div class="md:tw-w-1/2">
            <div
              v-if="product.images && product.images.length > 0"
              class="tw-bg-gray-50 tw-shadow-md tw-overflow-hidden tw-mb-4 tw-cursor-pointer"
              @click="openImageViewer(activeImageIndex)"
            >
              <img
                :src="product.images[activeImageIndex]"
                :alt="product.name"
                class="tw-w-full tw-h-auto tw-object-contain"
              />
            </div>

            <!-- Миниатюры для галереи -->
            <div
              v-if="product.images && product.images.length > 1"
              class="tw-flex tw-gap-4"
            >
              <button
                v-for="(image, index) in product.images"
                :key="index"
                @click="setActiveImage(index)"
                class="tw-w-20 tw-h-20 tw-bg-gray-50 tw-overflow-hidden tw-transition-all tw-duration-300"
                :class="
                  activeImageIndex === index
                    ? 'tw-ring-2 tw-ring-gray-800'
                    : 'tw-opacity-70 hover:tw-opacity-100'
                "
              >
                <img
                  :src="image"
                  :alt="`${product.name} - ${$t('image_viewer.image_of', {
                    current: index + 1,
                    total: product.images.length,
                  })}`"
                  class="tw-w-full tw-h-full tw-object-cover"
                />
              </button>
            </div>
          </div>

          <!-- Информация о товаре -->
          <div class="md:tw-w-1/2">
            <h1 class="tw-text-3xl tw-font-bold tw-text-gray-800 tw-mb-4">
              {{ product.name }}
            </h1>

            <div class="tw-mb-6">
              <p class="tw-text-2xl tw-font-bold tw-text-gray-800">
                {{ product.price }} {{ currencyMap[product.currency] }}
              </p>
            </div>

            <div class="tw-mb-6">
              <p class="tw-text-gray-600">{{ product.description }}</p>
            </div>

            <div v-if="product.characteristics" class="tw-mb-8">
              <h2 class="tw-text-lg tw-font-semibold tw-text-gray-800 tw-mb-2">
                {{ $t("product.characteristics") }}
              </h2>
              <div class="tw-space-y-2">
                <div
                  v-for="(value, key) in product.characteristics"
                  :key="key"
                  class="tw-flex"
                >
                  <div class="tw-w-1/3 tw-text-gray-600">
                    {{ $t(`product.${key}`) }}
                  </div>
                  <div class="tw-w-2/3 tw-text-gray-800">{{ value }}</div>
                </div>
              </div>
            </div>

            <!-- Выбор количества и добавление в корзину -->
            <div class="tw-mb-8">
              <label
                for="quantity"
                class="tw-block tw-text-sm tw-font-medium tw-text-gray-700 tw-mb-2"
                >{{ $t("product.quantity") }}</label
              >
              <div class="tw-flex tw-items-center tw-w-full tw-mb-4">
                <button
                  @click="decreaseQuantity"
                  class="tw-bg-gray-200 tw-text-gray-800 tw-p-3 tw-rounded-l-md tw-transition-colors hover:tw-bg-gray-300"
                  :disabled="quantity <= 1"
                >
                  <MinusIcon class="tw-w-4 tw-h-4" />
                </button>
                <input
                  type="number"
                  id="quantity"
                  v-model="quantity"
                  min="1"
                  class="tw-w-16 tw-text-center tw-py-2 tw-border-none focus:tw-outline-none"
                />
                <button
                  @click="increaseQuantity"
                  class="tw-bg-gray-200 tw-text-gray-800 tw-p-3 tw-rounded-r-md tw-transition-colors hover:tw-bg-gray-300"
                >
                  <PlusIcon class="tw-w-4 tw-h-4" />
                </button>
              </div>
              <div
                class="tw-bg-gray-100 tw-border-l-4 tw-border-gray-800 tw-p-3 tw-rounded-md tw-max-w-2xl tw-shadow-sm tw-mb-4"
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
              <div class="tw-flex tw-flex-col sm:tw-flex-row tw-gap-4">
                <button
                  @click="addToCart"
                  :disabled="isAddedToCart"
                  class="tw-bg-gray-200 tw-uppercase tw-text-gray-800 tw-py-3 tw-px-6 tw-rounded-md tw-shadow-sm tw-font-medium tw-transition-all tw-duration-300 hover:tw-bg-gray-300 focus:tw-outline-none focus:tw-ring-2 focus:tw-ring-offset-2 focus:tw-ring-gray-500 tw-flex-1 tw-flex tw-justify-center tw-items-center"
                  :class="{
                    'tw-bg-green-100 hover:tw-bg-green-100': isAddedToCart,
                  }"
                >
                  <span v-if="!isAddedToCart">{{
                    $t("product.add_to_cart")
                  }}</span>
                  <MarkIcon v-else class="tw-w-6 tw-h-6 tw-text-green-600" />
                </button>

                <button
                  @click="buyNow"
                  class="tw-bg-gray-800 tw-uppercase tw-text-white tw-py-3 tw-px-6 tw-rounded-md tw-shadow-sm tw-font-medium tw-transition-colors tw-duration-300 hover:tw-bg-gray-700 focus:tw-outline-none focus:tw-ring-2 focus:tw-ring-offset-2 focus:tw-ring-gray-500 tw-flex-1"
                >
                  {{ $t("product.buy_now") }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Похожие товары -->
        <div v-if="!isLoading && !relatedError && relatedProducts.length > 0">
          <RelatedProductsSlider
            :products="relatedProducts"
            :title="$t('product.similar_products')"
          />
        </div>
      </div>
    </div>

    <!-- Модальное окно с формой заказа -->
    <div
      v-if="isOrderFormVisible && product"
      class="tw-fixed tw-inset-0 tw-bg-black/50 tw-flex tw-items-center tw-justify-center tw-z-50 tw-p-4"
    >
      <div class="tw-bg-white tw-rounded-lg tw-shadow-xl tw-max-w-md tw-w-full">
        <div class="tw-p-8 tw-pb-0">
          <div class="tw-flex tw-justify-between tw-items-center tw-mb-4">
            <h3 class="tw-text-xl tw-font-semibold tw-text-gray-800">
              {{ $t("product.order_title") }}
            </h3>
            <button
              @click="isOrderFormVisible = false"
              class="tw-text-gray-500 hover:tw-text-gray-700 focus:tw-outline-none"
            >
              <CrossIcon class="tw-w-5 tw-h-5" />
            </button>
          </div>

          <div class="tw-mb-4">
            <p class="tw-text-gray-600">{{ product.name }}</p>
            <p class="tw-text-gray-800 tw-font-medium">
              {{ product.price }} {{ currencyMap[product.currency || "RUB"] }} ×
              {{ quantity }} = {{ product.price * quantity }}
              {{ currencyMap[product.currency || "RUB"] }}
            </p>
          </div>
        </div>

        <!-- Форма обратной связи с передачей информации о товаре -->
        <OrderFormHandler
          class="tw-p-8 tw-pt-0"
          :productToOrder="currentProductToOrder"
          @success="handleOrderSuccess"
        />
      </div>
    </div>

    <!-- Модальный просмотрщик изображений -->
    <ImageViewer
      v-if="product && product.images"
      :images="product.images"
      :is-open="isImageViewerOpen"
      :initial-index="activeImageIndex"
      @close="closeImageViewer"
    />
  </div>
</template>

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

/* Стили для скрытия стрелок в input number */
input[type="number"]::-webkit-inner-spin-button,
input[type="number"]::-webkit-outer-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
input[type="number"] {
  -moz-appearance: textfield;
}
</style>
