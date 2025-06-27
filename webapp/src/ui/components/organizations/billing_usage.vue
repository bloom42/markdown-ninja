<template>
  <div class="flex">
    <!-- <h3 class="text-base font-semibold text-gray-900">Last 30 days</h3> -->
    <dl class="grid grid-cols-1 divide-y divide-gray-200 overflow-hidden rounded-lg bg-white border border-gray-300 md:grid-cols-4 md:divide-x md:divide-y-0">
      <div v-for="metric in usageMetrics" :key="metric.name" class="px-4 py-5 sm:p-6">
        <dt class="text-base font-normal text-gray-900">{{ metric.name }}</dt>
        <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
          <div class="flex items-baseline text-2xl font-semibold text-gray-800">
            {{ +parseFloat((metric.used).toFixed(4)) }}
            <span class="ml-2 text-sm font-medium text-gray-500">/ {{ metric.allowed }}</span>
          </div>

          <!-- <div :class="[item.changeType === 'increase' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800', 'inline-flex items-baseline rounded-full px-2.5 py-0.5 text-sm font-medium md:mt-2 lg:mt-0']">
            <ArrowUpIcon v-if="item.changeType === 'increase'" class="-ml-1 mr-0.5 h-5 w-5 flex-shrink-0 self-center text-green-500" aria-hidden="true" />
            <ArrowDownIcon v-else class="-ml-1 mr-0.5 h-5 w-5 flex-shrink-0 self-center text-red-500" aria-hidden="true" />
            <span class="sr-only"> {{ item.changeType === 'increase' ? 'Increased' : 'Decreased' }} by </span>
            {{ item.change }}
          </div> -->
        </dd>
      </div>
    </dl>
  </div>
</template>

<script lang="ts" setup>
import type { OrganizationBillingUsage } from '@/api/model';
import { computed, type PropType } from 'vue';


interface UsageMetric {
  name: string;
  allowed: number;
  used: number;
}

// props
const props = defineProps({
  billingUsage: {
    type: Object as PropType<OrganizationBillingUsage>,
    required: true,
  }
});

// events

// composables

// lifecycle

// variables

// computed
const usageMetrics = computed((): UsageMetric[] => {
  return [
    {
      name: 'Websites',
      allowed: props.billingUsage.allowed_websites,
      used: props.billingUsage.used_websites,
    },
    {
      name: 'Storage (GB)',
      allowed: props.billingUsage.allowed_storage / 1_000_000_000,
      used: props.billingUsage.used_storage / 1_000_000_000,
    },
    {
      name: 'Staffs',
      allowed: props.billingUsage.allowed_staffs,
      used: props.billingUsage.used_staffs,
    },
    {
      name: 'Emails (Max, 1€ / 1000)',
      allowed: props.billingUsage.allowed_emails,
      used: props.billingUsage.used_emails,
    },
  ];
});

// watch

// functions
</script>
