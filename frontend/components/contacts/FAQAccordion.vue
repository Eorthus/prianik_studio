<script setup lang="ts">
import { ref } from "vue";

const props = defineProps<{
  faqs: {
    question: string;
    answer: string;
  }[];
}>();

// Состояние открытых/закрытых элементов
const openItems = ref({});

// Переключение состояния элемента
const toggleItem = (id) => {
  openItems.value = {
    ...openItems.value,
    [id]: !openItems.value[id],
  };
};

// Проверка, открыт ли элемент
const isOpen = (id) => {
  return !!openItems.value[id];
};
</script>

<template>
  <div class="tw-space-y-4 tw-container tw-mx-auto">
    <div
      v-for="(item, index) in props.faqs"
      :key="index"
      class="tw-border-b tw-border-gray-200"
    >
      <!-- Заголовок FAQ (всегда виден) -->
      <button
        @click="toggleItem(index)"
        class="tw-w-full tw-flex tw-justify-between tw-items-center tw-py-5 tw-text-left focus:tw-outline-none"
        :aria-expanded="isOpen(index)"
        :aria-controls="`faq-content-${index}`"
      >
        <h3 class="tw-text-lg tw-font-medium tw-text-gray-800">
          {{ item.question }}
        </h3>
      </button>

      <!-- Содержимое FAQ (скрывается/показывается) -->
      <div
        :id="`faq-content-${index}`"
        class="tw-overflow-hidden tw-transition-all tw-duration-300"
        :class="
          isOpen(index)
            ? 'tw-max-h-96 tw-opacity-100'
            : 'tw-max-h-0 tw-opacity-0'
        "
      >
        <div class="tw-pb-5">
          <p class="tw-text-gray-600">{{ item.answer }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
