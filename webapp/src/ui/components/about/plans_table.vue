<template>
  <!-- <div class="mt-16 flex justify-center">
    <RadioGroup v-model="frequency" class="grid grid-cols-2 gap-x-1 rounded-full p-1 text-center text-xs font-semibold leading-5 ring-1 ring-inset ring-gray-200">
      <RadioGroupLabel class="sr-only">Payment frequency</RadioGroupLabel>
      <RadioGroupOption as="template" v-for="option in frequencies" :key="option.value" :value="option" v-slot="{ checked }">
        <div :class="[checked ? 'bg-(--primary-color) text-white' : 'text-gray-500', 'cursor-pointer rounded-full px-2.5 py-1']">
          <span>{{ option.label }}</span>
        </div>
      </RadioGroupOption>
    </RadioGroup>
  </div> -->

  <div class="isolate mx-auto mt-10 grid max-w-md grid-cols-1 gap-8 md:mx-0 md:max-w-none md:grid-cols-3">
    <!-- <div v-for="tier in tiers" :key="tier.id" :class="[tier.mostPopular ? 'ring-2 ring-(--primary-color)' : 'ring-1 ring-gray-200', 'rounded-3xl p-8 xl:p-10']"> -->
    <div v-for="(plan, _index) in $store.pricing" :key="plan.id"
      :class="[plan.id === 'pro' ? 'ring-2 ring-(--primary-color)' : 'ring-1 ring-gray-200', 'rounded-3xl p-8 xl:p-10']">
      <div class="flex items-center justify-between gap-x-4">
        <!-- <h3 :id="tier.id" :class="[tier.mostPopular ? 'text-(--primary-color)' : 'text-gray-900', 'text-lg font-semibold leading-8']"> -->
        <h3 :id="plan.id" class="text-gray-900 text-4xl font-semibold leading-8">
          {{ plan.name }}
        </h3>
        <p v-if="plan.id === 'pro'" class="rounded-full bg-green-100 px-2.5 py-1 text-xs/5 font-medium text-green-800">
          Most Popular
        </p>
        <!-- <p v-if="plan.favorite" class="rounded-full bg-(--primary-color) px-2.5 py-1 text-xs font-semibold leading-5 text-(--primary-color)">
          Favorite
        </p> -->
      </div>
      <p class="mt-4 text-sm leading-6 text-gray-600 lg:h-8">
        {{ plan.description }}
      </p>
      <p v-if="plan.id === 'free'" class="mt-6 flex items-baseline gap-x-1">
        <span class="text-4xl font-medium tracking-tight">
          Free
        </span>
      </p>
      <p v-else-if="plan.id === 'pro'" class="mt-6 flex items-baseline gap-x-1">
        <span class="text-4xl font-medium tracking-tight">
          {{ plan.price }}€
        </span>
        <span class="text-sm font-medium leading-6 text-neutral-600 ml-2">
          per slot / month
        </span>
      </p>
      <p v-else class="mt-6 flex items-baseline gap-x-1">
        <span class="text-4xl font-medium tracking-tight">
          Custom
        </span>
      </p>
      <div class="flex w-full pt-6">
        <sl-button v-if="plan.id === 'pro'" variant="primary" @click="$pingoo.signup()">
          Get Started
          <!-- {{ plan.cta.text }} -->
        </sl-button>
        <sl-button v-else-if="plan.id === 'free'" variant="primary" @click="$pingoo.signup()">
          Try it free
        </sl-button>
        <RouterLink v-else to="/contact">
          <sl-button variant="primary">
            Contact Us
          </sl-button>
        </RouterLink>

      </div>
      <ul role="list" class="mt-6 space-y-3 text-sm leading-6 xl:mt-10">
        <!-- <li class="flex gap-x-3" v-if="$index !== 0">
          <ArrowLeftIcon class="h-6 w-5 flex-none text-(--primary-color)" aria-hidden="true" />
          Everything in {{ $store.pricing[$index-1].name }}
        </li> -->
        <li v-for="feature in plan.features" :key="feature" class="flex gap-x-3">
          <CheckIcon class="h-6 w-5 flex-none text-(--primary-color)" aria-hidden="true" />
          {{ feature }}
        </li>
      </ul>
    </div>

    <!-- <p class="text-neutral-500 -mt-5 text-sm">EU VAT may apply</p> -->
  </div>
</template>

<script lang="ts" setup>
import { CheckIcon } from '@heroicons/vue/20/solid'
import SlButton from '@shoelace-style/shoelace/dist/components/button/button.js';
import { usePingoo } from '@/pingoo/pingoo';
import { useStore } from '@/app/store';

// props

// events

// composables
const $pingoo = usePingoo();
const $store = useStore();

// lifecycle

// variables

// computed

// watch

// functions
</script>
