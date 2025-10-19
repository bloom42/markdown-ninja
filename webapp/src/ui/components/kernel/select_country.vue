<template>
  <sl-select label="Country" class="w-full"
    :value="selected" @sl-change="selected = $event.target.value" :disabled="readonly">
    <sl-option v-for="country in countries" :value="country.code">{{ country.name }}</sl-option>
  </sl-select>
</template>

<script lang="ts" setup>
import { computed, type PropType } from 'vue';
import countries from 'mdninja-js/src/data/countries';
import SlSelect from '@shoelace-style/shoelace/dist/components/select/select.js';
import SlOption from '@shoelace-style/shoelace/dist/components/option/option.js';

// props
const props = defineProps({
  modelValue: {
    type: String as PropType<string>,
    required: false,
  },
  label: {
    type: String as PropType<string>,
    required: false,
    default: 'Country',
  },
  readonly: {
    type: Boolean,
    required: false,
    default: false,
  },
});


// events
const $emit = defineEmits(['update:modelValue']);

// composables

// lifecycle

// variables

// computed
const selected = computed({
  get(): string | undefined {
    return props.modelValue;
  },
  set(value: string | undefined) {
    $emit('update:modelValue', value);
  }
});

// watch

// functions
</script>
