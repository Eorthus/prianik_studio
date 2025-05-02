<template>
  <!-- Модальное окно успешной отправки или ошибки -->
  <transition name="modal">
    <div
      v-if="props.isOpen"
      class="tw-fixed tw-inset-0 tw-bg-black/50 tw-flex tw-items-center tw-justify-center tw-z-50"
    >
      <div
        class="tw-bg-white tw-rounded-lg tw-shadow-xl tw-max-w-md tw-w-full tw-p-12 tw-text-center"
      >
        <div
          v-if="!props.errorMessage"
          class="tw-flex tw-justify-center tw-mb-4"
        >
          <div class="tw-bg-green-100 tw-rounded-full tw-p-3">
            <SuccessIcon class="tw-w-8 tw-h-8 tw-text-green-500" />
          </div>
        </div>

        <div v-else class="tw-flex tw-justify-center tw-mb-4">
          <div class="tw-bg-red-100 tw-rounded-full tw-p-3">
            <RejectIcon class="tw-w-8 tw-h-8 tw-text-red-500" />
          </div>
        </div>

        <h3 class="tw-text-xl tw-font-semibold tw-text-gray-800 tw-mb-2">
          {{
            !props.errorMessage
              ? t("form.success_title")
              : t("form.error_title")
          }}
        </h3>
        <p class="tw-text-gray-600">
          {{
            !props.errorMessage
              ? t("form.success_message")
              : props.errorMessage
          }}
        </p>
      </div>
    </div>
  </transition>
</template>

<script lang="ts" setup>
import { useI18n } from "vue-i18n";
import SuccessIcon from "../icons/SuccessIcon.vue";
import RejectIcon from "../icons/RejectIcon.vue";

const { t } = useI18n();

const props = defineProps<{
  isOpen: boolean;
  errorMessage?: string;
}>();
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
