<template>
  <div>
    <div class="flex flex-col space-y-4">

      <div class="flex w-full">
        <sl-input :value="modelValue.line1" @input="modelValue.line1 = $event.target.value" type="text"
          :readonly="readonly" placeholder="Line 1" label="Line 1" />
      </div>

      <div class="flex w-full">
        <sl-input :value="modelValue.line2" @input="modelValue.line2 = $event.target.value" type="text"
          :readonly="readonly" placeholder="Line 2" label="Line 2" />
      </div>



      <div class="flex flex-row space-x-4">
        <sl-input :value="modelValue.city" @input="modelValue.city = $event.target.value" type="text"
          :readonly="readonly" placeholder="City" label="City" />

        <sl-input :value="modelValue.postal_code" @input="modelValue.postal_code = $event.target.value" type="text"
          :readonly="readonly" placeholder="Postal Code" label="Postal Code" />
      </div>

      <div class="flex flex-row space-x-4">
        <sl-input :value="modelValue.state" @input="modelValue.state = $event.target.value" type="text"
          :readonly="readonly" placeholder="State, Province or Region (optional)" label="State, Province or Region (optional)" />

        <!-- @ts-ignore: ignore next line -->
        <SelectCountry v-model="modelValue.country_code" :readonly="readonly" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import type { Address } from '@/api/model';
import { type PropType } from 'vue';
import SelectCountry from './select_country.vue';
import SlInput from '@shoelace-style/shoelace/dist/components/input/input.js';

// props
defineProps({
  modelValue: {
    type: Object as PropType<Address>,
    required: true,
  },
  readonly: {
    type: Boolean,
    required: false,
    default: false,
  },
});

// events

// composables

// lifecycle

// variables

// computed

// watch

// functions
</script>
