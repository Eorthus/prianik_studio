<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from "vue";
import LogoFull from "./icons/LogoFull.vue";
import CartIcon from "./icons/CartIcon.vue";
import { useCart } from "~/shared/useCart";
import LocaleSelect from "~/lib/LocaleSelect.vue";
import type { LangType } from "./types";
import CrossIcon from "./icons/CrossIcon.vue";
// Используем встроенный композабл для i18n
//@ts-expect-error no type for usei18n
const { locale, locales, t, loadLocaleMessages, setLocale } = useI18n();

const { cartCounter } = useCart();

const localeReady = ref(false);

const changeLang = (code: LangType) => {
  setLocale(code);
};

onMounted(async () => {
  await loadLocaleMessages(locale.value);
  localeReady.value = true;
});

const isMobileMenuOpen = ref(false);
const isMobile = ref(false);

const menuMap = computed(() => {
  const _ = locale.value;

  const arr = [
    { name: t("header.catalog"), path: "/catalog" },
    { name: t("header.gallery"), path: "/gallery" },
    { name: t("header.contacts"), path: "/contacts" },
  ];

  if (cartCounter.value) {
    arr.push({ name: t("header.cart"), path: "/cart" });
    return arr;
  }
  return arr;
});

const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value;
  // Блокировка скролла при открытом меню
  if (isMobileMenuOpen.value) {
    document.body.style.overflow = "hidden";
  } else {
    document.body.style.overflow = "";
  }
};

const closeMobileMenu = () => {
  isMobileMenuOpen.value = false;
  document.body.style.overflow = "";
};

// Проверка размера экрана
const checkScreenSize = () => {
  isMobile.value = window?.innerWidth < 768;
  if (!isMobile.value && isMobileMenuOpen.value) {
    closeMobileMenu();
  }
};

onMounted(() => {
  checkScreenSize();
  window.addEventListener("resize", checkScreenSize);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", checkScreenSize);
});
</script>

<template>
  <header
    v-if="localeReady"
    class="tw-fixed tw-bg-white tw-top-0 tw-left-0 tw-right-0 tw-z-50 tw-transition-all tw-duration-300 tw-h-[80px]"
  >
    <div class="tw-container tw-mx-auto tw-px-4">
      <div class="tw-flex tw-justify-between tw-items-center">
        <!-- Логотип -->
        <div class="tw-flex tw-items-center">
          <NuxtLink
            to="/"
            class="tw-transition-colors tw-duration-300 hover:tw-opacity-70"
          >
            <div class="tw-text-2xl tw-font-bold tw-text-gray-800 tw-w-[200px]">
              <LogoFull />
            </div>
          </NuxtLink>
        </div>

        <!-- Десктопное меню -->
        <div class="tw-hidden md:tw-flex tw-items-center tw-gap-8 tw-py-6">
          <NuxtLink
            v-for="(menu, i) in menuMap"
            :key="menu.name"
            :to="menu.path"
            class="tw-text-gray-800 tw-transition-colors tw-duration-300 hover:tw-text-gray-500"
          >
            <div
              v-if="i === 3"
              class="after:tw-content-[''] after:tw-bg-red-500 tw-relative after:tw-w-[10px] after:tw-h-[10px] after:tw-rounded-full after:tw-top-[2px] after:tw-right-0 after:tw-absolute"
            >
              <CartIcon class="tw-w-6 tw-h-6" />
            </div>
            <template v-else>
              {{ menu.name }}
            </template>
          </NuxtLink>
          <!-- Языковой селектор -->
          <LocaleSelect
            :model-value="locale"
            :locales="locales"
            @update:modelValue="changeLang"
          />
        </div>

        <!-- Бургер кнопка -->
        <div class="md:tw-hidden tw-py-6">
          <button
            @click="toggleMobileMenu"
            class="tw-text-gray-800 tw-p-2 tw-transition-colors tw-duration-300 hover:tw-text-gray-500"
            aria-label="Меню"
          >
            <div v-if="!isMobileMenuOpen" class="tw-w-6 tw-h-5 tw-relative">
              <span
                class="tw-absolute tw-h-0.5 tw-w-full tw-bg-gray-800 tw-top-0 tw-left-0 tw-transition-all"
              ></span>
              <span
                class="tw-absolute tw-h-0.5 tw-w-full tw-bg-gray-800 tw-top-2 tw-left-0 tw-transition-all"
              ></span>
              <span
                class="tw-absolute tw-h-0.5 tw-w-full tw-bg-gray-800 tw-bottom-0 tw-left-0 tw-transition-all"
              ></span>
            </div>
            <div v-else class="tw-w-6 tw-h-5 tw-relative">
              <span
                class="tw-absolute tw-h-0.5 tw-w-full tw-bg-gray-800 tw-top-2 tw-left-0 tw-transform tw-rotate-45 tw-transition-all"
              ></span>
              <span
                class="tw-absolute tw-h-0.5 tw-w-full tw-bg-gray-800 tw-top-2 tw-left-0 tw-transform -tw-rotate-45 tw-transition-all"
              ></span>
            </div>
          </button>
        </div>
      </div>
    </div>

    <!-- Мобильное меню -->
    <div
      v-if="isMobileMenuOpen"
      class="tw-fixed tw-inset-0 tw-bg-white tw-z-50 tw-flex tw-flex-col tw-justify-center tw-items-center tw-transition-all tw-duration-300 tw-py-6"
    >
      <!-- Кнопка закрытия (крестик) -->
      <button
        @click="closeMobileMenu"
        class="tw-absolute tw-top-6 tw-right-6 tw-text-gray-800 tw-p-2 tw-transition-colors tw-duration-300 hover:tw-text-gray-500"
        aria-label="Закрыть меню"
      >
        <CrossIcon class="tw-w-6 tw-h-6" />
      </button>
      <div class="tw-flex tw-flex-col tw-gap-8 tw-text-center tw-items-center">
        <NuxtLink
          v-for="(menu, i) in menuMap"
          :key="menu.name"
          :to="menu.path"
          @click="closeMobileMenu"
          class="tw-text-2xl tw-text-gray-800 tw-transition-colors tw-duration-300 hover:tw-text-gray-500"
        >
          <div
            v-if="i === 3"
            class="after:tw-content-[''] after:tw-bg-red-500 tw-relative after:tw-w-[10px] after:tw-h-[10px] after:tw-rounded-full after:tw-top-[2px] after:tw-right-0 after:tw-absolute"
          >
            <CartIcon class="tw-w-6 tw-h-6 tw-m-auto" />
          </div>
          <template v-else> {{ menu.name }} </template>
        </NuxtLink>
        <!-- Языковой селектор -->
        <LocaleSelect
          :model-value="locale"
          :locales="locales"
          @update:modelValue="changeLang"
        />
      </div>
    </div>
  </header>
</template>
