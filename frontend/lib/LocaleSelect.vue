<template>
  <div class="tw-relative tw-inline-block">
    <!-- Текущая локаль -->
    <button
      @click="toggleDropdown"
      class="tw-flex tw-items-center tw-gap-2 tw-px-3 tw-py-1 tw-border tw-rounded-md tw-bg-white tw-cursor-pointer"
    >
      <div class="tw-w-5 tw-h-5">
        <FlagIcon :name="props.modelValue" />
      </div>
      <DownArrowIcon class="tw-w-4 tw-h-4" />
    </button>

    <!-- Выпадающий список -->
    <ul
      v-if="isOpen"
      class="tw-absolute tw-mt-1 tw-bg-white tw-border tw-rounded-md tw-shadow-lg tw-z-50 tw-w-full"
    >
      <li
        v-for="loc in props.locales"
        :key="loc.code"
        @click="selectLanguage(loc.code)"
        class="tw-flex tw-items-center tw-gap-2 tw-py-1 tw-px-3 tw-cursor-pointer hover:tw-bg-gray-100"
      >
        <FlagIcon :name="loc.code" class="tw-w-5 tw-h-5" />
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { onClickOutside } from "@vueuse/core";
import FlagIcon from "~/components/icons/FlagIcon.vue";
import { ref } from "vue";
import type { LangType } from "~/components";
import DownArrowIcon from "~/components/icons/DownArrowIcon.vue";

const props = defineProps<{
  modelValue: LangType;
  locales: { code: LangType }[];
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: LangType): void;
}>();

const isOpen = ref(false);

function toggleDropdown() {
  isOpen.value = !isOpen.value;
}

function selectLanguage(code: LangType) {
  emit("update:modelValue", code);
  isOpen.value = false;
}

// Закрыть по клику вне
const el = ref();
onClickOutside(el, () => {
  isOpen.value = false;
});
</script>
