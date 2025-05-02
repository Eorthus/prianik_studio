<template>
  <div>
    <label
      :for="props.name"
      class="tw-block tw-text-sm tw-font-medium tw-text-gray-700"
      >{{ $t(props.label) }}</label
    >
    <input
      :type="props.type"
      :id="props.name"
      v-model="model"
      class="tw-mt-1 tw-block tw-w-full tw-border tw-border-gray-300 tw-rounded-md tw-shadow-sm tw-py-2 tw-px-3 focus:tw-outline-none focus:tw-ring-gray-500 focus:tw-border-gray-500"
      :placeholder="$t(props.placeholder)"
      :required="required"
      :class="{ 'tw-border-red-500': error }"
    />
    <p v-if="error" class="tw-text-red-500 tw-text-sm tw-mt-1">
      {{ error }}
    </p>
  </div>
</template>

<script lang="ts" setup>
import { computed } from "vue";

const props = withDefaults(
  defineProps<{
    modelValue: string;
    error?: string;
    placeholder: string;
    label: string;
    required?: boolean;
    name:string
    type: string
  }>(),
  {
    required: true,
  }
);

const emit = defineEmits<{
  "update:model-value": [value: string];
}>();

const model = computed({
  get() {
    return props.modelValue;
  },
  set(obj) {
    emit("update:model-value", obj);
  },
});
</script>
