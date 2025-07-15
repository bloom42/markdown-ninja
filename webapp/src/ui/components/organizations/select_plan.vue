<template>
  <sl-select :value="selected" @sl-change="selected = $event.target.value" label="Plan">
    <sl-option v-for="plan in plans" :value="plan.id">
      {{ plan.name }}
      <span v-if="plan.price">- {{ plan.price }} € / month</span>
    </sl-option>
  </sl-select>
</template>

<script lang="ts" setup>
import { computed, type ModelRef, type PropType } from 'vue';
import SlSelect from '@shoelace-style/shoelace/dist/components/select/select.js';
import SlOption from '@shoelace-style/shoelace/dist/components/option/option.js';
import { useStore } from '@/app/store';

// props
const selected: ModelRef<string> = defineModel({ required: true });

const props = defineProps({
  allPlans: {
    type: Boolean as PropType<boolean>,
    required: false,
    default: false,
  }
})


// events

// composables
const $store = useStore();

// lifecycle

// variables

// computed
const plans = computed(() => {
  if (props.allPlans) {
    return $store.pricing;
  }

  return $store.pricing.filter((plan) => plan.id !== 'enterprise');
})

// watch

// functions
</script>
