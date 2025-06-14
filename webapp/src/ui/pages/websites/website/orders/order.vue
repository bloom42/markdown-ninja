<template>
  <div class="flex-1">
    <div class="px-4 sm:px-6 md:px-0 mb-4">
      <h1 class="text-2xl font-bold text-gray-900">Order {{ orderId }}</h1>
    </div>

    <div class="rounded-md bg-red-50 p-4 my-5" v-if="error">
      <div class="flex">
        <div class="ml-3">
          <p class="text-sm text-red-700">
            {{ error }}
          </p>
        </div>
      </div>
    </div>

    <div v-if="order" class="flex flex-col">
      <div class="flex">
        <b>Date</b>: {{ order.created_at }}
      </div>
      <div class="flex">
        <b>Country</b>: {{ countryName(order.country) }} ({{ order.country }})
      </div>
      <div class="flex">
        <b>Status</b>: <POrderStatus :status="order.status" class="ml-2" />
      </div>
      <div class="flex">
        <b>Total</b>: {{ order.total_amount }} {{ order.currency }}
      </div>
      <div class="flex">
        <b>Contact</b>:&nbsp;
        <RouterLink :to="contactUrl(order.contact_id)" class="text-(--primary-color) hover:underline">
          {{ order.contact_id }}
        </RouterLink>
      </div>
      <div class="flex">
        <b>Email</b>: {{ order.email }}
      </div>
      <div class="flex">
        <b>Stripe Checkout Session ID</b>: {{ order.stripe_checkout_session_id }}
      </div>

      <div class="flex">
        <b>Stripe Payment Intent</b>:&nbsp;
        <a v-if="order.stripe_payment_intent_id" :href="`https://dashboard.stripe.com/payments/${order.stripe_payment_intent_id}`" target="_blank" rel="noopener noreferrer" class="text-(--primary-color) hover:underline">
          {{ order.stripe_payment_intent_id }}
        </a>
        <span v-else>-</span>
      </div>

      <div class="flex">
        <b>Stripe Invoice</b>:&nbsp;
        <a v-if="order.stripe_invoice_id" :href="`https://dashboard.stripe.com/invoices/${order.stripe_invoice_id}`" target="_blank" rel="noopener noreferrer" class="text-(--primary-color) hover:underline">
          {{ order.stripe_invoice_id }}
        </a>
        <span v-else>-</span>
      </div>

      <div class="flex max-w-2xl flex flex-col">
        <b>Additional invoice information:</b> <br />
        <sl-textarea readonly placeholder="(empty)"
          :value="order.additional_invoice_information" />
      </div>

      <div class="flex flex-col mt-5 space-y-5">
        <div class="flex">
          <h1 class="text-xl font-extrabold text-gray-900">Line Items</h1>
        </div>
        <POrderLineItemsList :website-id="websiteId" :line-items="order.line_items!" />
      </div>

      <div class="flex flex-col mt-5 space-y-2">
        <div class="flex">
          <h1 class="text-xl font-extrabold text-gray-900">Refunds</h1>
        </div>

        <div class="flex">
          <sl-button variant="primary" @click="onCreateRefundClicked" :disabled="order.status !== OrderStatus.Completed">
            Create Refund
          </sl-button>
        </div>

        <div class="flex">
          <RefundsList :refunds="order.refunds!" />
        </div>
      </div>
    </div>
  </div>

  <RefundDialog v-model="showRefundDialog" :refund="refundToShow" :order-id="orderId"
    @created="onRefundCreated"
  />
</template>

<script lang="ts" setup>
import { OrderStatus, type Order, type Refund } from '@/api/model';
import { useMdninja } from '@/api/mdninja';
import { onBeforeMount, ref, type Ref } from 'vue';
import { useRoute } from 'vue-router';
import POrderLineItemsList from '@/ui/components/products/order_line_items_list.vue';
import POrderStatus from '@/ui/components/products/order_status.vue';
import RefundsList from '@/ui/components/products/refunds_list.vue';
import RefundDialog from '@/ui/components/products/refund_dialog.vue';
import SlButton from '@shoelace-style/shoelace/dist/components/button/button.js';
import SlTextarea from '@shoelace-style/shoelace/dist/components/textarea/textarea.js';
import { countryName } from '@/libs/countries';

// props

// events

// composables
const $route = useRoute();
const $mdninja = useMdninja();

// lifecycle
onBeforeMount(() => fetchData());

// variables
const orderId = $route.params.order_id as string;
const websiteId = $route.params.website_id as string;

let loading = ref(false);
let error = ref('');
let order: Ref<Order | null> = ref(null);
let showRefundDialog = ref(false);
let refundToShow: Ref<Refund | null> = ref(null);

// computed

// watch

// functions
async function fetchData() {
  loading.value = true;
  error.value = '';

  try {
    order.value = await $mdninja.getOrder(orderId);
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

function contactUrl(contactId: string): string {
  return `/websites/${websiteId}/contacts/${contactId}`;
}

function onCreateRefundClicked() {
  showRefundDialog.value = true;
}

function onRefundCreated(refund: Refund) {
  order.value!.refunds!.push(refund);
}
</script>
