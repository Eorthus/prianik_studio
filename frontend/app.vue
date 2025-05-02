<template>
  <NuxtLayout name="default">
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { useHead, useSeoMeta } from "nuxt/app";
import { getSeoMetadata } from "~/shared/useSeo";
import { ref, watch, onMounted, computed } from "vue";
import { useApiService } from "~/services/api";
import { useCategories } from "~/shared/useCategories";

const { getCategories, isLoading } = useApiService();

const { categories } = useCategories();

const { t, locale } = useI18n();

const loadCategories = async () => {
  try {
    const response = await getCategories(locale.value);
    if (response.success && response.data) {
      categories.value = response.data;
    } else {
      console.error("Ошибка загрузке категорий:", response.error);
    }
  } catch (err) {
    console.error("Ошибка при загрузке категорий:", err);
  }
};
onMounted(loadCategories);

// Получаем данные из composables
const route = useRoute();
const router = useRouter();

// Создаем реактивную переменную для SEO данных
const seoData = ref(getSeoMetadata(route, t, locale.value));

// Функция обновления SEO метаданных
function updateSeoMetadata() {
  seoData.value = getSeoMetadata(route, t, locale.value);
}

// Устанавливаем SEO-метатеги с использованием реактивного объекта
useSeoMeta({
  title: computed(() => seoData.value.title),
  description: computed(() => seoData.value.description),
  ogTitle: computed(() => seoData.value.title),
  ogDescription: computed(() => seoData.value.description),
  ogImage: computed(() => seoData.value.image),
  ogUrl: computed(() => seoData.value.pageUrl),

  //@ts-expect-error
  ogType: computed(() => seoData.value.type),
  ogLocale: computed(() => seoData.value.locale),
  twitterTitle: computed(() => seoData.value.title),
  twitterDescription: computed(() => seoData.value.description),
  twitterImage: computed(() => seoData.value.image),
  //@ts-expect-error
  twitterCard: computed(() => seoData.value.twitterCard),
});

// Устанавливаем дополнительные head-теги
useHead({
  htmlAttrs: {
    lang: computed(() => seoData.value.locale),
  },
  link: computed(() => [
    { rel: "canonical", href: seoData.value.pageUrl },
    ...seoData.value.alternateLinks,
  ]),
});

// Отслеживаем изменения локали
watch(locale, () => {
  updateSeoMetadata();
  loadCategories();
});

// Используем хук маршрутизации для обновления метаданных
onMounted(() => {
  updateSeoMetadata();

  router.afterEach(() => {
    updateSeoMetadata();
  });
});
</script>
