<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { ref, computed, watch, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import DownArrowIcon from "~/components/icons/DownArrowIcon.vue";
import FilterIcon from "~/components/icons/FilterIcon.vue";
import LeftArrowIcon from "~/components/icons/LeftArrowIcon.vue";
import RightArrowIcon from "~/components/icons/RightArrowIcon.vue";
import { useApiService } from "~/services/api";
import { useCategories } from "~/shared/useCategories";
import { currencyMap } from "~/components";
import LoaderView from "~/components/LoaderView.vue";
import type { Product } from "~/components";

const { t, locale } = useI18n();
const route = useRoute();
const router = useRouter();

// API сервис
const { getProducts, isLoading } = useApiService();
const { categories, getCategory, getSubcategory } = useCategories();

// Мобильное представление
const isMobile = ref(false);

// Фильтры
const searchQuery = ref("");
const sortPrice = ref<"asc" | "desc" | null>(null);

// Выбранная категория и подкатегория
const selectedCategory = ref<number | null>(null);
const selectedSubcategory = ref<number | null>(null);

// Мобильное меню категорий
const isMobileCategoryMenuOpen = ref(false);

// Параметры пагинации
const pageSize = ref(10); // Количество товаров на странице
const currentPage = ref(1); // Текущая страница
const totalItems = ref(0); // Общее количество товаров (будет приходить с бэкенда)

// Данные продуктов
const products = ref<Product[]>([]);
const loading = computed(() => isLoading.value);
const error = ref<string | null>(null);

const toggleMobileCategoryMenu = () => {
  isMobileCategoryMenuOpen.value = !isMobileCategoryMenuOpen.value;
};

const loadProducts = async () => {
  error.value = null;

  try {
    const response = await getProducts(
      currentPage.value,
      pageSize.value,
      selectedCategory.value,
      selectedSubcategory.value,
      searchQuery.value,
      sortPrice.value,
      locale.value
    );

    if (response.success && response.data) {
      products.value = response.data.items;
      // Убедимся, что totalItems устанавливается корректно
      totalItems.value = response.data.total_items;
    } else {
      error.value = response.error || "Не удалось загрузить товары";
      console.error("Ошибка загрузки товаров:", response.error);
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Неизвестная ошибка";
    console.error("Ошибка при загрузке товаров:", err);
  }
};

// Общее количество страниц
const totalPages = computed(() => {
  return Math.ceil(totalItems.value / pageSize.value);
});

// Массив номеров страниц для пагинации
const pageNumbers = computed(() => {
  const pages = [];
  const maxVisiblePages = 5; // Максимальное количество отображаемых страниц

  if (totalPages.value <= maxVisiblePages) {
    // Если страниц меньше максимального количества, показываем все
    for (let i = 1; i <= totalPages.value; i++) {
      pages.push(i);
    }
  } else {
    // Иначе показываем страницы вокруг текущей
    let startPage = Math.max(
      currentPage.value - Math.floor(maxVisiblePages / 2),
      1
    );
    let endPage = startPage + maxVisiblePages - 1;

    if (endPage > totalPages.value) {
      endPage = totalPages.value;
      startPage = Math.max(endPage - maxVisiblePages + 1, 1);
    }

    for (let i = startPage; i <= endPage; i++) {
      pages.push(i);
    }

    // Добавляем многоточие и первую/последнюю страницы
    if (startPage > 1) {
      pages.unshift("...");
      pages.unshift(1);
    }

    if (endPage < totalPages.value) {
      pages.push("...");
      pages.push(totalPages.value);
    }
  }

  return pages;
});

// Переключение ценовой сортировки
const togglePriceSort = () => {
  if (sortPrice.value === null) {
    sortPrice.value = "asc";
  } else if (sortPrice.value === "asc") {
    sortPrice.value = "desc";
  } else {
    sortPrice.value = null;
  }

  // Немедленно обновляем URL и загружаем данные с сервера
  updateUrlQuery();
};

// Сброс фильтра и пагинации
const resetFilters = () => {
  searchQuery.value = "";
  sortPrice.value = null;
  selectedCategory.value = null;
  selectedSubcategory.value = null;
  currentPage.value = 1;

  // Обновляем URL
  router.push({ query: {} });

  // Загружаем продукты без фильтров
  loadProducts();
};

// Переход на страницу
const goToPage = (page) => {
  if (typeof page === "number" && page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
    scrollToTop();
    updateUrlQuery();
  }
};

// Прокрутка страницы наверх при смене страницы
const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: "smooth",
  });
};

// Обработка выбора категории
const selectCategory = (categoryId: number | null) => {
  if (selectedCategory.value === categoryId) {
    // Если выбрана та же категория - сбрасываем выбор
    selectedCategory.value = null;
    selectedSubcategory.value = null;
  } else {
    selectedCategory.value = categoryId;
    selectedSubcategory.value = null;
  }

  // Сбрасываем пагинацию при смене категории
  currentPage.value = 1;

  // Закрываем мобильное меню при выборе
  if (isMobileCategoryMenuOpen.value) {
    isMobileCategoryMenuOpen.value = false;
  }

  // Обновляем URL и загружаем товары
  updateUrlQuery();
};

// Обработка выбора подкатегории
const selectSubcategory = (subcategory) => {
  if (selectedSubcategory.value === subcategory.id || !subcategory) {
    // Если выбрана та же категория - сбрасываем выбор
    selectedSubcategory.value = null;
  } else {
    selectedSubcategory.value = subcategory?.id;
    selectedCategory.value = subcategory?.parent_id;
  }
  // Сбрасываем пагинацию при смене подкатегории
  currentPage.value = 1;

  // Закрываем мобильное меню при выборе
  if (isMobileCategoryMenuOpen.value) {
    isMobileCategoryMenuOpen.value = false;
  }

  // Обновляем URL и загружаем товары
  updateUrlQuery();
};

// Обновление URL при изменении параметров
const updateUrlQuery = () => {
  const query: Record<string, string> = {};

  if (selectedCategory.value) {
    query.category = selectedCategory.value.toString();
  }

  if (selectedSubcategory.value) {
    query.subcategory = selectedSubcategory.value.toString();
  }

  if (currentPage.value > 1) {
    query.page = currentPage.value.toString();
  }

  if (pageSize.value !== 10) {
    // если отличается от значения по умолчанию
    query.pageSize = pageSize.value.toString();
  }

  if (searchQuery.value) {
    query.search = searchQuery.value;
  }

  if (sortPrice.value) {
    query.sort = sortPrice.value;
  }

  router.push({ query });

  // Загружаем продукты с новыми параметрами
  loadProducts();
};

// Обновление заголовка страницы в зависимости от выбранной категории/подкатегории
const pageTitle = computed(() => {
  if (selectedCategory.value && selectedSubcategory.value) {
    const category = categories.value?.find(
      (cat) => cat.id === selectedCategory.value
    );
    if (category) {
      const subcategory = category.subcategories?.find(
        (subcat) => subcat.id === selectedSubcategory.value
      );
      if (subcategory) {
        return `${category.name} - ${subcategory.name}`;
      }
    }
  }

  if (selectedCategory.value) {
    const category = categories.value?.find(
      (cat) => cat.id === selectedCategory.value
    );
    if (category) {
      return category.name;
    }
  }

  return t("catalog.title");
});

// Наблюдение за изменением URL (для обработки прямых переходов и возвратов по истории)
watch(
  () => route.query,
  (newQuery) => {
    if (newQuery.category) {
      selectedCategory.value =
        parseInt(newQuery.category as string, 10) || null;
    } else {
      selectedCategory.value = null;
    }

    if (newQuery.subcategory) {
      selectedSubcategory.value =
        parseInt(newQuery.subcategory as string, 10) || null;
    } else {
      selectedSubcategory.value = null;
    }

    if (newQuery.page) {
      currentPage.value = parseInt(newQuery.page as string, 10) || 1;
    } else {
      currentPage.value = 1;
    }

    if (newQuery.pageSize) {
      pageSize.value = parseInt(newQuery.pageSize as string, 10) || 10;
    }

    if (newQuery.search) {
      searchQuery.value = newQuery.search as string;
    } else {
      searchQuery.value = "";
    }

    if (newQuery.sort) {
      sortPrice.value =
        (newQuery.sort as string) === "asc" ||
        (newQuery.sort as string) === "desc"
          ? (newQuery.sort as "asc" | "desc")
          : null;
    } else {
      sortPrice.value = null;
    }

    // Загружаем товары при изменении URL
    loadProducts();
  },
  { immediate: true }
);

// Выполнение поиска при нажатии на Enter
const handleSearch = (event) => {
  if (event.key === "Enter") {
    currentPage.value = 1;
    updateUrlQuery();
  }
};

// Наблюдение за изменением поисковой строки для сброса пагинации
watch(searchQuery, () => {
  currentPage.value = 1;
});

// Проверка размера экрана
const checkScreenSize = () => {
  isMobile.value = window?.innerWidth < 768;
};

onMounted(() => {
  checkScreenSize();
  window.addEventListener("resize", checkScreenSize);

  // Инициализация начальных значений из URL
  if (route.query.category) {
    selectedCategory.value =
      parseInt(route.query.category as string, 10) || null;
  }

  if (route.query.subcategory) {
    selectedSubcategory.value =
      parseInt(route.query.subcategory as string, 10) || null;
  }

  if (route.query.page) {
    currentPage.value = parseInt(route.query.page as string, 10) || 1;
  }

  if (route.query.pageSize) {
    pageSize.value = parseInt(route.query.pageSize as string, 10) || 10;
  }

  if (route.query.search) {
    searchQuery.value = route.query.search as string;
  }

  if (route.query.sort) {
    sortPrice.value =
      (route.query.sort as string) === "asc" ||
      (route.query.sort as string) === "desc"
        ? (route.query.sort as "asc" | "desc")
        : null;
  }

  // Загружаем товары с учетом параметров URL
  loadProducts();
});
</script>

<template>
  <div class="tw-py-6">
    <div class="tw-container tw-mx-auto tw-px-4">
      <h1
        class="tw-text-3xl tw-font-bold tw-text-gray-800 tw-mb-8 tw-text-center"
      >
        {{ pageTitle }}
      </h1>

      <div class="tw-flex tw-flex-col md:tw-flex-row tw-gap-8">
        <!-- Мобильная кнопка открытия меню категорий -->
        <div class="md:tw-hidden">
          <button
            @click="toggleMobileCategoryMenu"
            class="tw-w-full tw-flex tw-items-center tw-justify-between tw-bg-gray-100 tw-text-gray-800 tw-px-4 tw-py-3 tw-rounded-md tw-shadow-sm"
          >
            <span>{{
              isMobileCategoryMenuOpen
                ? $t("catalog.hide_categories")
                : $t("catalog.show_categories")
            }}</span>

            <DownArrowIcon
              class="tw-w-5 tw-h-5 tw-transition-transform tw-duration-300"
              :class="{ 'tw-rotate-180': isMobileCategoryMenuOpen }"
            />
          </button>
        </div>

        <!-- Боковое меню категорий (на десктопе всегда видимо, на мобильных устройствах скрыто/раскрыто) -->
        <div
          class="md:tw-w-1/4 tw-transition-all tw-duration-300"
          :class="[
            isMobileCategoryMenuOpen || !isMobile ? 'tw-block' : 'tw-hidden',
          ]"
        >
          <div class="tw-bg-gray-50 tw-shadow-md tw-p-6">
            <nav class="tw-space-y-2">
              <!-- Категории и подкатегории -->
              <div
                v-for="category in categories"
                :key="category.id"
                class="tw-space-y-1"
              >
                <button
                  @click="selectCategory(category.id)"
                  class="tw-w-full tw-text-left tw-py-2 tw-px-3 tw-rounded-md tw-transition-colors tw-duration-300 tw-font-medium hover:tw-opacity-70"
                  :class="{
                    'tw-pointer-events-none tw-opacity-70':
                      selectedCategory === category.id,
                  }"
                >
                  {{ category.name }}
                </button>

                <!-- Подкатегории - всегда отображаем, если они есть -->
                <div
                  v-if="category.subcategories?.length"
                  class="tw-pl-4 tw-space-y-1"
                >
                  <button
                    v-for="subcategory in category.subcategories"
                    :key="subcategory.id"
                    @click="selectSubcategory(subcategory)"
                    class="tw-w-full tw-text-left tw-py-1 tw-px-3 tw-rounded-md tw-transition-colors tw-duration-300 hover:tw-opacity-70"
                    :class="{
                      'tw-pointer-events-none tw-opacity-70':
                        selectedSubcategory === subcategory.id,
                    }"
                  >
                    {{ subcategory.name }}
                  </button>
                </div>
              </div>
            </nav>
          </div>
        </div>

        <!-- Список товаров -->
        <div class="md:tw-w-3/4">
          <!-- Информация о количестве найденных товаров и активных фильтрах -->
          <div class="tw-mb-6">
            <div
              class="tw-flex tw-justify-between tw-items-center tw-mb-4 tw-gap-6"
            >
              <p class="tw-text-gray-600">
                {{ $t("catalog.found_items") }} {{ totalItems }}
              </p>
              <div class="tw-flex tw-items-center tw-gap-6">
                <input
                  type="text"
                  id="search"
                  v-model="searchQuery"
                  class="tw-block tw-w-full tw-max-w-[250px] tw-border tw-border-gray-300 tw-rounded-md tw-shadow-sm tw-py-2 tw-px-3 focus:tw-outline-none focus:tw-ring-gray-500 focus:tw-border-gray-500"
                  :placeholder="$t('catalog.search_placeholder')"
                  @keyup="handleSearch"
                />
                <button
                  class="tw-text-sm tw-text-gray-800 tw-flex tw-items-center hover:tw-text-gray-600 tw-transition-colors tw-w-10 tw-h-10"
                  @click="togglePriceSort"
                  :title="$t('catalog.sort_by_price')"
                >
                  <FilterIcon
                    :class="{
                      'tw-rotate-180': sortPrice === 'desc',
                    }"
                  />
                </button>
              </div>
            </div>

            <!-- Активные фильтры -->
            <div
              v-if="
                selectedCategory ||
                selectedSubcategory ||
                searchQuery ||
                sortPrice
              "
              class="tw-flex tw-flex-wrap tw-gap-2 tw-mb-4"
            >
              <div
                v-if="selectedCategory"
                class="tw-bg-gray-200 tw-text-gray-800 tw-px-3 tw-py-1 tw-rounded-full tw-text-sm tw-flex tw-items-center"
              >
                <span class="tw-mr-1">{{ getCategory(selectedCategory) }}</span>
                <button
                  @click="selectCategory(null)"
                  class="tw-text-gray-500 hover:tw-text-gray-700"
                >
                  ×
                </button>
              </div>
              <div
                v-if="selectedSubcategory"
                class="tw-bg-gray-200 tw-text-gray-800 tw-px-3 tw-py-1 tw-rounded-full tw-text-sm tw-flex tw-items-center"
              >
                <span class="tw-mr-1">
                  {{ getSubcategory(selectedSubcategory, selectedCategory) }}
                </span>
                <button
                  @click="selectSubcategory(null)"
                  class="tw-text-gray-500 hover:tw-text-gray-700"
                >
                  ×
                </button>
              </div>
              <div
                v-if="searchQuery"
                class="tw-bg-gray-200 tw-text-gray-800 tw-px-3 tw-py-1 tw-rounded-full tw-text-sm tw-flex tw-items-center"
              >
                <span class="tw-mr-1"
                  >{{ $t("catalog.search") }}: {{ searchQuery }}</span
                >
                <button
                  @click="
                    searchQuery = '';
                    updateUrlQuery();
                  "
                  class="tw-text-gray-500 hover:tw-text-gray-700"
                >
                  ×
                </button>
              </div>
              <div
                v-if="sortPrice"
                class="tw-bg-gray-200 tw-text-gray-800 tw-px-3 tw-py-1 tw-rounded-full tw-text-sm tw-flex tw-items-center"
              >
                <span class="tw-mr-1"
                  >{{ $t("catalog.price") }}:
                  {{ sortPrice === "asc" ? "↑" : "↓" }}</span
                >
                <button
                  @click="
                    sortPrice = null;
                    updateUrlQuery();
                  "
                  class="tw-text-gray-500 hover:tw-text-gray-700"
                >
                  ×
                </button>
              </div>
              <button
                @click="resetFilters"
                class="tw-text-gray-800 hover:tw-opacity-70 tw-px-3 tw-py-1 tw-rounded-full tw-text-sm"
              >
                {{ $t("catalog.reset_filters") }}
              </button>
            </div>
          </div>

          <!-- Состояние загрузки -->
          <LoaderView v-if="loading" />

          <!-- Сообщение об ошибке -->
          <div
            v-else-if="error"
            class="tw-bg-red-50 tw-rounded-lg tw-shadow-md tw-p-8 tw-text-center"
          >
            <p class="tw-text-red-600">{{ error }}</p>
            <button
              @click="loadProducts"
              class="tw-mt-4 tw-bg-gray-800 tw-text-white tw-py-2 tw-px-4 tw-rounded-md tw-shadow-sm tw-text-sm tw-font-medium tw-transition-colors tw-duration-300 hover:tw-bg-gray-700 focus:tw-outline-none focus:tw-ring-2 focus:tw-ring-offset-2 focus:tw-ring-gray-500"
            >
              {{ $t("catalog.try_again") }}
            </button>
          </div>

          <!-- Сетка товаров -->
          <div
            v-else-if="products?.length > 0"
            class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 lg:tw-grid-cols-3 tw-gap-4"
          >
            <div
              v-for="product in products"
              :key="product.id"
              class="tw-bg-white tw-shadow-md tw-overflow-hidden tw-transition-all tw-duration-300 hover:tw-shadow-lg"
            >
              <NuxtLink :to="`/catalog/${product.id}`" class="tw-block">
                <div class="tw-aspect-w-1 tw-aspect-h-1 tw-overflow-hidden">
                  <img
                    :src="
                      product.images && product.images?.length > 0
                        ? product.images[0]
                        : 'https://s3.stroi-news.ru/img/masterskaya-kartinki-1.jpg'
                    "
                    :alt="product.name"
                    class="tw-w-full tw-h-full tw-object-cover tw-transition-transform tw-duration-500 hover:tw-scale-110"
                  />
                </div>
                <div class="tw-p-4">
                  <h3
                    class="tw-text-lg tw-font-medium tw-text-gray-800 tw-mb-2 tw-line-clamp-2"
                  >
                    {{ product.name }}
                  </h3>
                  <p class="tw-text-lg tw-font-semibold tw-text-gray-800">
                    {{ product.price }}
                    {{ currencyMap[product.currency || "RUB"] }}
                  </p>
                </div>
              </NuxtLink>
            </div>
          </div>

          <!-- Пагинация -->
          <div
            v-if="totalPages > 1"
            class="tw-flex tw-justify-center tw-items-center tw-mt-8 tw-space-x-2"
          >
            <!-- Кнопка "Назад" -->
            <button
              @click="goToPage(currentPage - 1)"
              :disabled="currentPage === 1"
              class="tw-px-3 tw-py-2 tw-rounded tw-bg-white tw-border tw-border-gray-300 tw-text-gray-800 hover:tw-bg-gray-100 tw-text-sm"
              :class="{
                'tw-opacity-50 tw-pointer-events-none': currentPage === 1,
              }"
            >
              <LeftArrowIcon class="tw-w-5 tw-h-5" />
            </button>

            <!-- Номера страниц -->
            <template v-for="(page, index) in pageNumbers" :key="index">
              <button
                v-if="typeof page === 'number'"
                @click="goToPage(page)"
                class="tw-px-3 tw-py-2 tw-rounded tw-font-medium tw-text-sm tw-border"
                :class="
                  page === currentPage
                    ? 'tw-bg-gray-800 tw-text-white tw-border-gray-800'
                    : 'tw-bg-white tw-text-gray-800 tw-border-gray-300 hover:tw-bg-gray-100'
                "
              >
                {{ page }}
              </button>
              <span v-else class="tw-px-3 tw-py-2 tw-text-gray-800 tw-text-sm">
                {{ page }}
              </span>
            </template>

            <!-- Кнопка "Вперед" -->
            <button
              @click="goToPage(currentPage + 1)"
              :disabled="currentPage === totalPages"
              class="tw-px-3 tw-py-2 tw-rounded tw-bg-white tw-border tw-border-gray-300 tw-text-gray-800 hover:tw-bg-gray-100 tw-text-sm"
              :class="{
                'tw-opacity-50 tw-cursor-not-allowed':
                  currentPage === totalPages,
              }"
            >
              <RightArrowIcon class="tw-w-5 tw-h-5" />
            </button>
          </div>

          <!-- Сообщение, если товары не найдены -->
          <div
            v-if="!loading && !error && products?.length === 0"
            class="tw-bg-gray-50 tw-shadow-md tw-p-8 tw-text-center"
          >
            <p class="tw-text-gray-600">
              {{ $t("catalog.no_results") }}
            </p>
            <button
              @click="resetFilters"
              class="tw-mt-4 tw-bg-gray-800 tw-text-white tw-py-2 tw-px-4 tw-rounded-md tw-shadow-sm tw-text-sm tw-font-medium tw-transition-colors tw-duration-300 hover:tw-bg-gray-700 focus:tw-outline-none focus:tw-ring-2 focus:tw-ring-offset-2 focus:tw-ring-gray-500"
            >
              {{ $t("catalog.reset_filters") }}
            </button>
          </div>
        </div>
      </div>
    </div>
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
</style>
