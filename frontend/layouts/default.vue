<script setup lang="ts">
import MainFooter from "~/components/MainFooter.vue";
import MainHeader from "~/components/MainHeader.vue";
import { ref, computed } from "vue";
import { useRoute } from "vue-router";
import CrossIcon from "~/components/icons/CrossIcon.vue";

// Для уведомления-баннера
const showNotification = ref(true);
const closeNotification = () => {
  showNotification.value = false;
};

// Проверяем нужно ли показывать уведомление (на главной, каталоге и карточке товара)
const route = useRoute();
const shouldShowNotification = computed(() => {
  return (
    route.path.includes("/contacts") === false &&
    route.path.includes("/cart") === false
  );
});
</script>

<template>
  <div class="tw-flex tw-flex-col tw-min-h-screen">
    <MainHeader />

    <!-- Основной контент с отступом под фиксированный хедер -->
    <main class="tw-flex-grow tw-pt-20">
      <!-- Уведомление -->
      <div
        v-if="showNotification && shouldShowNotification"
        class="tw-bg-gray-100 tw-bg-opacity-90 tw-py-3 tw-px-4 tw-text-center tw-w-full tw-relative"
      >
        <div class="tw-container tw-mx-auto">
          <p class="tw-text-gray-800">
            {{ $t("notification.custom_idea") }}
            <NuxtLink
              to="/contacts"
              class="tw-underline tw-text-gray-600 hover:tw-text-gray-500 tw-transition-colors tw-duration-300"
            >
              {{ $t("contacts.contact_us") }}
            </NuxtLink>
          </p>
          <button
            @click="closeNotification"
            class="tw-absolute tw-right-4 tw-top-1/2 tw-transform -tw-translate-y-1/2 tw-text-gray-600 hover:tw-text-gray-500 tw-transition-colors tw-duration-300"
          >
            <span class="tw-sr-only">{{ $t("notification.close") }}</span>
            <CrossIcon class="tw-w-5 tw-h-5" />
          </button>
        </div>
      </div>
      <slot />
    </main>

    <MainFooter />
  </div>
</template>
