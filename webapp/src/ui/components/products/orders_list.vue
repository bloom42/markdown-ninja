<template>
  <div class="overflow-x-auto min-w-full">
    <div class="py-2 align-middle inline-block min-w-full">
      <div class="overflow-hidden border border-gray-300 sm:rounded-lg">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr class="max-w-0">
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Order ID
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Date
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Country
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Total
              </th>
            </tr>
          </thead>
          <tbody class="min-w-full bg-white divide-y divide-gray-200">
            <!-- <tr v-for="contact in contacts" :key="contact.id"> -->
            <RouterLink :to="orderUrl(order)" v-for="order in orders" :key="order.id"
                class="table-row cursor-pointer min-w-full">
              <td class="px-6 py-4 whitespace-nowrap w-1/4">
                <div class="text-md font-medium text-gray-900 truncate">
                  {{ order.id }}
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap w-1/4">
                <div class="text-md font-medium text-gray-900 truncate">
                  {{ date(order.created_at) }}
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap w-1/8">
                <div class="text-md font-medium text-gray-900 truncate">
                  {{ order.country }}
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap w-1/8">
                <POrderStatus :status="order.status" />
              </td>
              <td class="px-6 py-4 whitespace-nowrap w-1/8">
                <div class="text-md font-medium text-gray-900 truncate">
                  {{ order.total_amount }} {{ order.currency }}
                </div>
              </td>
            </RouterLink>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { type OrderMetadata } from '@/api/model'
import { type PropType } from 'vue'
import { useRoute } from 'vue-router';
import date from 'mdninja-js/src/libs/date';
import POrderStatus from '@/ui/components/products/order_status.vue';

// props
defineProps({
  orders: {
    type: Array as PropType<OrderMetadata[]>,
    required: true,
  },
});

// events

// composables
const $route = useRoute();

// lifecycle

// variables
const websiteId = $route.params.website_id as string;

// computed

// watch

// functions
function orderUrl(order: OrderMetadata): string {
  return `/websites/${websiteId}/orders/${order.id}`;
}
</script>
