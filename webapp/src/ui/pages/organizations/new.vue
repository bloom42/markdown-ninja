<template>
  <div class="min-h-full flex flex-col justify-center sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
        New Organization
      </h2>
    </div>

    <div class="mt-5 sm:mx-auto sm:w-full sm:max-w-xl">
      <div class="bg-white py-8 px-4 sm:rounded-lg sm:px-10 space-y-6">

        <div class="rounded-md bg-red-50 p-4 mb-3" v-if="error">
          <div class="flex">
            <div class="ml-3">
              <p class="text-sm text-red-700">
                {{ error }}
              </p>
            </div>
          </div>
        </div>

        <sl-input :value="name" @input="name = $event.target.value.trim()"
          :disabled="loading" placeholder="My Organization" label="Name"
        />

        <div class="flex">
          <fieldset>
            <legend class="font-medium leading-6 text-gray-900 p-0 m-0">Select a plan</legend>
            <p class="font-normal mt-1 mb-3 text-sm text-(--primary-color)">
              <a href="/pricing" target="_blank" rel="noopener" class="hover:underline">
                Learn more about pricing
              </a>
            </p>
            <RadioGroup v-model="selectedPlan" class="grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-4">
              <RadioGroupOption as="template" v-for="plan in plans" :key="plan.id" :value="plan" :aria-label="plan.name" v-slot="{ active, checked }">
                <div :class="[active ? 'border-(--primary-color) ring-1 ring-(--primary-color)' : 'border-gray-300', 'relative flex cursor-pointer rounded-lg border bg-white p-4 shadow-xs focus:outline-none']">
                  <span class="flex flex-1">
                    <span class="flex flex-col">
                      <span class="block text-md font-medium text-gray-900">{{ plan.name }}</span>
                      <!-- <span class="mt-1 flex items-center text-sm text-gray-500">{{ plan.description }}</span> -->
                      <span class="mt-6 text-sm text-gray-900">{{ plan.price }}</span>
                    </span>
                  </span>
                  <CheckCircleIcon :class="[!checked ? 'invisible' : '', 'h-5 w-5 text-(--primary-color)']" aria-hidden="true" />
                  <span :class="[active ? 'border' : 'border-2', checked ? 'border-(--primary-color)' : 'border-transparent', 'pointer-events-none absolute -inset-px rounded-lg']" aria-hidden="true" />
                </div>
              </RadioGroupOption>
            </RadioGroup>
          </fieldset>
        </div>

        <div class="flex">
          <sl-input :value="billingEmail" @input="billingEmail = $event.target.value.trim()" type="email"
            :disabled="loading" placeholder="my@email.com" label="Billing Email"
          />
        </div>


        <div class="mt-5 flex justify-between">
          <RouterLink :to="backRoute">
            <sl-button outline>
              Back
            </sl-button>
          </RouterLink>

          <sl-button variant="primary" :loading="loading" @click="createOrganization()">
            Create Organization
          </sl-button>
        </div>

      </div>
    </div>

  </div>
</template>

<script lang="ts" setup>
import type { CreateOrganizationInput } from '@/api/model';
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router';
import { useMdninja } from '@/api/mdninja';
import { useStore } from '@/app/store';
import { RadioGroup, RadioGroupOption } from '@headlessui/vue'
import { CheckCircleIcon } from '@heroicons/vue/20/solid';
import SlButton from '@shoelace-style/shoelace/dist/components/button/button.js';
import { oneRouteUp } from '@/libs/router_utils';
import SlInput from '@shoelace-style/shoelace/dist/components/input/input.js';


const plans = [
  {
    id: 'free',
    name: 'Free',
    price: 'Try Markdown Ninja for free',
  },
  {
    id: 'pro',
    name: 'Pro',
    price: '10€ / month',
  }
];

// props

// events

// composables
const $mdninja = useMdninja();
const $router = useRouter();
const $store = useStore();
const $route = useRoute();

// lifecycle

// variables
const backRoute = oneRouteUp($route.path);

let name = ref('');
let error = ref('');
let loading = ref(false);
let selectedPlan = ref(plans[1]);
let billingEmail = ref($store.userEmail ?? '');

// computed

// watch

// functions
async function createOrganization() {
  loading.value = true;
  error.value = '';
  const input: CreateOrganizationInput = {
    name: name.value.trim(),
    plan: selectedPlan.value.id,
    billing_email: billingEmail.value.trim(),
  };

  try {
    const res = await $mdninja.createOrganization(input);

    // if stripe_checkout_session_url is provided, we redirect to it
    if (res.stripe_checkout_session_url) {
      location.href = res.stripe_checkout_session_url;
      return;
    }

    $store.addOrUpdateOrganization(res.organization);
    $router.push(`/organizations/${res.organization.id}`);
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}
</script>
