<template>
  <div class="flex flex-col mb-10">
    <div class="px-4 sm:px-6 md:px-0">
      <h1 class="text-3xl font-extrabold text-gray-900">Billing</h1>
    </div>

    <div class="flex my-3">
      <p>
        Any questions? Send us a message at
        <a :href="`mailto:${$store.contact_email}`" class="text-(--primary-color) underline">
          {{ $store.contact_email }}
        </a>
      </p>
    </div>

    <div class="rounded-md bg-red-50 p-4 my-3" v-if="error">
      <div class="flex">
        <div class="ml-3">
          <p class="text-sm text-red-700">
            {{ error }}
          </p>
        </div>
      </div>
    </div>

    <div v-if="organization" class="flex flex-col">
      <div class="px-4 sm:px-6 md:px-0 mb-3">
        <h2 class="text-2xl font-bold text-gray-900">Subscription</h2>

        <p class="font-normal text-sm text-(--primary-color)">
              <a href="/pricing" target="_blank" rel="noopener" class="underline">
                Learn more about the different plans
              </a>
            </p>
      </div>

      <div class="rounded-md bg-red-50 p-4 my-3" v-if="organization.payment_due">
        <sl-button variant="danger" :loading="loading" @click="gotToStripeCustomerPortal()">
          Update payment method & Pay due invoice
        </sl-button>
      </div>

      <div class="flex flex-col">
        <div v-if="editingSubscription" class="flex flex-col">
          <SelectPlan v-model="plan" />

          <div v-if="plan !== 'free'" class="flex flex-col">
            <div class="flex flex-col my-2">
              <sl-input label="Additional Slots"
                :value="extraSlots" @input="extraSlots = parseInt($event.target.value, 10)" min="0" type="number"
                placeholder="0"
              />
            </div>
          </div>

          <div class="flex mt-4">
            <p class="text-lg font-medium text-gray-900">
              Total: {{ subscriptionTotalPrice }}€ / month
              <span v-if="subscriptionTotalPrice !== 0">(billed yearly)</span>
            </p>
          </div>

          <div class="flex mt-4 justify-between">
            <div class="flex">
              <sl-button outline @click="editingSubscription = false">
                Close
              </sl-button>
              <sl-button variant="primary" @click="updateSubscription()" :loading="loading" class="ml-3">
                Update Subscription
              </sl-button>
            </div>

            <!-- <div class="flex">
              <sl-button variant="danger" :loading="loading" @click="cancelSubscription()">
                Cancel Subscription
              </sl-button>
            </div> -->
          </div>
        </div>

        <div v-else class="flex border border-gray-300 w-fit rounded-md px-3 py-5">
          <div class="flex flex-col mx-2">
            <p class="text-base h-full content-center font-medium text-gray-900 ">
              Plan: {{ organization.plan }}
            </p>
            <p v-if="organization.plan !== 'free'" class="text-sm h-full content-center text-gray-500">
              Extra slots: {{ organization.extra_slots }}
            </p>
          </div>

          <div class="flex ml-5" v-if="canUpdateSubscription">
            <sl-button @click="editingSubscription = true" variant="text">
              Update
            </sl-button>
          </div>
        </div>

      </div>



      <div class="flex flex-col mt-5 space-y-3">
        <h2 class="text-2xl font-bold text-gray-900">Usage</h2>

        <BillingUsage :billing-usage="billingUsage!" />
      </div>


      <div class="flex flex-col">
        <div class="flex my-5">
          <h2 class="text-2xl font-bold">Billing information</h2>
        </div>

        <div v-if="editingBillingInformation" class="flex flex-col">
          <div class="mt-4 flex">
            <BillingInformationForm v-model="billingInformation" />
          </div>

          <div class="mt-5 flex">
            <sl-button outline @click="editingBillingInformation = false">
              Close
            </sl-button>
            <sl-button variant="primary" @click="updateBillingInformation()" :loading="loading" class="ml-3">
              Update
            </sl-button>
          </div>
        </div>
        <div v-else class="flex w-fit border border-gray-300 rounded-md max-w-4xl px-3 py-2">
          <div class="ml-3 flex flex-col">
            <div class="block text-base font-medium">
              {{ billingInformation.name }}
            </div>
            <div class="block text-sm text-gray-500">
              {{ billingInformation.email }}
            </div>
            <div class="block text-sm text-gray-500">
              {{ billingInformation.address_line1 }}
            </div>
            <div class="block text-sm text-gray-500">
              {{ billingInformation.postal_code }}, {{ billingInformation.city }}, {{ billingInformation.country_code }}
            </div>
            <div v-if="billingInformation.tax_id" class="block text-sm text-gray-500">
              VAT: {{ billingInformation.tax_id }}
            </div>
          </div>

          <div class="flex ml-20">
            <sl-button @click="editingBillingInformation = true" variant="text" class="flex">
              Edit
            </sl-button>
          </div>
        </div>

      </div>




      <div v-if="organization.stripe_customer" class="flex flex-col px-4 sm:px-6 md:px-0 mt-5">
        <div class="flex">
          <h2 class="text-2xl font-bold">Payment Methods & Invoices</h2>
        </div>

        <div class="flex mt-5">
          <sl-button variant="primary" @click="gotToStripeCustomerPortal()" :loading="loading">
            View / Update
          </sl-button>
        </div>
      </div>


    </div>
  </div>
</template>

<script lang="ts" setup>
import { useMdninja } from '@/api/mdninja';
import type { BillingInformation, GetOrganizationInput, Organization, OrganizationBillingUsage, OrganizationGetStripeCustomerPortalUrlInput, OrganizationUpdateSubscriptionInput, UpdateOrganizationInput } from '@/api/model';
import { computed, onBeforeMount, type Ref, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import BillingInformationForm from '@/ui/components/organizations/billing_information_form.vue';
import { useStore } from '@/app/store';
import SelectPlan from '@/ui/components/organizations/select_plan.vue';
import SlButton from '@shoelace-style/shoelace/dist/components/button/button.js';
import SlInput from '@shoelace-style/shoelace/dist/components/input/input.js';
import BillingUsage from '@/ui/components/organizations/billing_usage.vue';

// props

// events

// composables
const $route = useRoute();
const $mdninja = useMdninja();
const $store = useStore();
const $router = useRouter();

// lifecycle
onBeforeMount(() => fetchData());


// variables
const organizationId = $route.params.organization_id as string;

let organization: Ref<Organization | null> = ref(null);
let loading = ref(false);
let error = ref('');
let editingBillingInformation = ref(false);
let editingSubscription = ref(false);

let billingInformation: Ref<BillingInformation> = ref({
    name: '',
    email: '',
    address_line1: '',
    address_line2: '',
    postal_code: '',
    city: '',
    state: '',
    country_code: '',
    tax_id: '',
});
let billingUsage: Ref<OrganizationBillingUsage | null> = ref(null);
let plan = ref('');
let extraSlots = ref(0);

// computed
const subscriptionTotalPrice = computed(() => {
  let total = 0;
  switch (plan.value) {
    case 'free':
      return 0;
    case 'pro':
      total += 5;
      break;
  }

  total += Math.abs(extraSlots.value) * 5;

  return total;
});

const canUpdateSubscription = computed(() => {
  return organization.value?.plan !== 'enterprise';
});

// watch

// functions
function resetValues() {
  if (organization.value) {
    billingInformation.value = organization.value.billing_information;
    plan.value = organization.value.plan;
    extraSlots.value = organization.value.extra_slots;
  }

  billingInformation.value.tax_id = billingInformation.value.tax_id ?? '';

  if (billingInformation.value.country_code === 'XX' || billingInformation.value.country_code === '') {
    billingInformation.value.country_code = $store.country;
  }
}

async function fetchData() {
  loading.value = true;
  error.value = '';
  const input: GetOrganizationInput = {
    id: organizationId,
  };

  try {
    const [org, resBillingUsage] = await Promise.all([
      $mdninja.getOrganization(input),
      $mdninja.getorganizationBillingUsage(organizationId),
    ])
    $store.addOrUpdateOrganization(org);
    organization.value = org;
    billingUsage.value = resBillingUsage;
    resetValues();
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}


async function updateBillingInformation() {
  loading.value = true;
  error.value = '';
  const input: UpdateOrganizationInput = {
    id: organizationId,
    billing_information: billingInformation.value,
  };

  try {
    organization.value = await $mdninja.updateOrganization(input);
    resetValues();
    editingBillingInformation.value = false;
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function updateSubscription() {
  error.value = '';

  if (plan.value === 'free') {
    if (!confirm("Do you really want to cancel your subscription? All your benefits will be removed now.")) {
      return;
    }
    extraSlots.value = 0;
  }

  loading.value = true;

  const input: OrganizationUpdateSubscriptionInput = {
    organization_id: organizationId,
    plan: plan.value,
    extra_slots: extraSlots.value,
  };

  try {
    const res = await $mdninja.organizationUpdateSubscription(input);
    editingSubscription.value = false;
    if (res.stripe_checkout_session_url) {
      location.href = res.stripe_checkout_session_url;
      return;
    }
    $router.push({ path: `/organizations/${organizationId}/billing/checkout/complete`, query: { plan: plan.value, redirect_to: $route.path }});
    return;
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function gotToStripeCustomerPortal() {
  loading.value = true;
  error.value = '';
  const input: OrganizationGetStripeCustomerPortalUrlInput = {
    organization_id: organizationId,
  };

  try {
    const res = await $mdninja.organizationGetStripeCustomerPortal(input);
    location.href = res.stripe_customer_portal_url;
    return;
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

// function cancelSubscription() {
//   if (!confirm("Do you really want to cancel your subscription now? All your benefits will be removed now.")) {
//     return;
//   }

//   plan.value = 'free';
//   updateSubscription();
// }
</script>
