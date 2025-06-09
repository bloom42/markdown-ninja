<template>
  <div class="flex-1">
    <div class="rounded-md bg-red-50 p-4" v-if="error">
      <div class="flex">
        <div class="ml-3">
          <p class="text-sm text-red-700">
            {{ error }}
          </p>
        </div>
      </div>
    </div>

    <div v-if="website" class="flex flex-col">
      <!-- <div class="flex flex-row justify-between">
        <div class="flex flex-row justify-left">
          <h1 class="text-2xl font-medium text-gray-900 hover:underline content-center">
            <a :href="websiteUrl" target="_blank" rel="noopener">
              {{ website.primary_domain }}
            </a>
          </h1>

          <span class="ml-3 cursor-default content-center rounded-md bg-white px-3 py-2 text-sm text-gray-900 ring-1 ring-inset ring-gray-300">
            Last 30 Days
          </span>
          <span v-if="$store.isAdmin" class="text-xl ml-5 font-medium content-center">
            {{ (website?.revenue ?? 0) + ' ' + website?.currency }}
          </span>
        </div>
        <div class="flex">
          <RouterLink :to="newPostUrl">
            <sl-button variant="primary" class="ml-5">
              <PlusIcon class="-ml-1 mr-2 h-5 w-5 inline" aria-hidden="true" />
              New Post
            </sl-button>
          </RouterLink>
        </div>
      </div> -->

      <!-- <div class="border-b border-b-gray-900/10 lg:border-t lg:border-t-gray-900/5"> -->
      <div class="mt-2 rounded-md border border-gray-900/10">
        <dl class="mx-auto grid max-w-7xl grid-cols-1 sm:grid-cols-3 lg:px-2 xl:px-0">
          <div v-for="(stat, statIdx) in stats" :key="stat.name" :class="[statIdx % 2 === 1 ? 'sm:border-l' : statIdx === 2 ? 'lg:border-l' : '', 'flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1 border-t border-gray-900/5 px-4 py-5 sm:px-6 lg:border-t-0 xl:px-8']">
            <dt class="text-md">{{ stat.name }}</dt>
            <dd v-if="stat.change !== null" :class="[stat.change < 0 ? 'text-red-600' : 'text-green-600', 'text-sm']">
              {{ stat.change < 0 ? '-' : '+' }}{{ stat.change.toLocaleString('en-US') }}
            </dd>
            <dd class="w-full flex-none text-3xl font-medium tracking-tight">
              {{ stat.value }}
            </dd>
          </div>
        </dl>
      </div>


      <!-- <div class="flex w-full  mt-5">
        <dl class="flex w-full grid grid-cols-2 gap-5 sm:grid-cols-4">
          <div class="overflow-hidden rounded-lg bg-white px-4 py-5 border sm:p-6">
            <dt class="truncate text-sm font-medium text-gray-500">Revenue</dt>
            <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">
              {{ (website.revenue ?? 0) + ' ' + website.currency }}
            </dd>
          </div>
          <div class="overflow-hidden rounded-lg bg-white px-4 py-5 border sm:p-6">
            <dt class="truncate text-sm font-medium text-gray-500">Subscribers</dt>
            <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">{{ website.subscribers }}
              <span class="text-green-600 ml-2 align-center text-sm font-medium">
              + {{ analyticsData?.new_subscribers }}
              </span>
            </dd>
          </div>
          <div class="overflow-hidden rounded-lg bg-white px-4 py-5 border sm:p-6">
            <dt class="truncate text-sm font-medium text-gray-500">Page Views</dt>
            <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">{{ analyticsData!.total_page_views }}</dd>
          </div>
          <div class="overflow-hidden rounded-lg bg-white px-4 py-5 border sm:p-6">
            <dt class="truncate text-sm font-medium text-gray-500">Visitors</dt>
            <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">{{ analyticsData!.total_visitors }}</dd>
          </div>
        </dl>
      </div> -->

      <div class="flex h-72 my-8">
        <AnalyticsChart v-if="analyticsData" :data="analyticsData" />
      </div>

      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">

        <div class="flex flex-col">
          <div class="flex text-lg font-bold">
            Top Pages
          </div>
          <div class="flex">
            <div class="overflow-x-auto w-full">
              <div class="inline-block min-w-full align-middle">
                <table class="min-w-full divide-y divide-gray-300">
                  <thead>
                    <tr>
                      <th scope="col" class="py-3.5 pl-4 pr-3 text-left font-medium sm:pl-0">Page</th>
                      <th scope="col" class="px-3 py-3.5 text-left font-medium">Visitors</th>
                    </tr>
                  </thead>
                  <tbody class="">
                    <tr v-for="page in analyticsData!.pages" :key="page.label">
                      <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-0">{{ page.label }}</td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ page.count.toLocaleString('en-US') }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

        <div class="flex flex-col">
          <div class="flex text-lg font-bold">
            Top Referrers
          </div>
          <div class="flex">
            <div class="overflow-x-auto w-full">
              <div class="inline-block min-w-full align-middle">
                <table class="min-w-full divide-y divide-gray-300">
                  <thead>
                    <tr>
                      <th scope="col" class="py-3.5 pl-4 pr-3 text-left font-medium sm:pl-0">Referrer</th>
                      <th scope="col" class="px-3 py-3.5 text-left font-medium">Visitors</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="referrer in analyticsData!.referrers" :key="referrer.label">
                      <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-0">{{ referrer.label }}</td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ referrer.count.toLocaleString('en-US') }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

        <div class="flex flex-col">
          <div class="flex text-lg font-bold">
            Top Countries
          </div>
          <div class="flex">
            <div class="overflow-x-auto w-full">
              <div class="inline-block min-w-full align-middle">
                <table class="min-w-full divide-y divide-gray-300">
                  <thead>
                    <tr>
                      <th scope="col" class="py-3.5 pl-4 pr-3 text-left font-medium sm:pl-0">Country</th>
                      <th scope="col" class="px-3 py-3.5 text-left font-medium">Visitors</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="country in analyticsData!.countries" :key="country.label">
                      <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-0">{{ country.label }}</td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ country.count.toLocaleString('en-US') }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

        <div class="flex flex-col">
          <div class="flex text-lg font-bold">
            Top Browsers
          </div>
          <div class="flex">
            <div class="overflow-x-auto w-full">
              <div class="inline-block min-w-full align-middle">
                <table class="min-w-full divide-y divide-gray-300">
                  <thead>
                    <tr>
                      <th scope="col" class="py-3.5 pl-4 pr-3 text-left font-medium sm:pl-0">Browser</th>
                      <th scope="col" class="px-3 py-3.5 text-left font-medium">Visitors</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="browser in analyticsData!.browsers" :key="browser.label">
                      <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-0">{{ browser.label }}</td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ browser.count.toLocaleString('en-US') }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

      </div>
    </div>

  </div>
</template>

<script lang="ts" setup>
import type { AnalyticsData, GetAnalyticsDataInput } from '@/api/model';
import type { GetWebsiteInput, Website } from '@/api/model';
import { getCountryFlag } from 'mdninja-js/src/libs/flag';
import { computed, onBeforeMount, ref, type Ref } from 'vue';
import { useRoute } from 'vue-router';
import { useMdninja } from '@/api/mdninja';
// import AnalyticsChart from '@/ui/components/websites/analytics_chart.vue';
import { defineAsyncComponent } from 'vue'
const AnalyticsChart = defineAsyncComponent(() =>
  import('@/ui/components/websites/analytics_chart.vue')
);

type Stat = {
  name: string;
  value: string;
  change: number | null,
}

// props

// events

// composables
const $route = useRoute();
const $mdninja = useMdninja();

// lifecycle
onBeforeMount(() => fetchData());

// variables
const websiteId = $route.params.website_id as string;
// const newPostUrl = computed(() => `/websites/${websiteId}/posts/new`);

let loading = ref(false);
let error = ref('');
let analyticsData: Ref<AnalyticsData | null> = ref(null);
let website: Ref<Website | null> = ref(null);


// computed
// const websiteUrl = computed((): string => {
//   if (website.value) {
//     return $mdninja.generateWebsiteUrl(website.value);
//   }
//   return '';
// });

const stats = computed((): Stat[] => {
  return [
    {
      name: 'Subscribers',
      value: (website.value?.subscribers ?? 0).toLocaleString('en-US'),
      change: analyticsData.value?.new_subscribers ?? 0,
    },
    {
      name: 'Page Views',
      value: (analyticsData.value?.total_page_views ?? 0).toLocaleString('en-US'),
      change: null,
    },
    {
      name: 'Visitors',
      value: (analyticsData.value?.total_visitors ?? 0).toLocaleString('en-US'),
      change: null,
    },
  ]
})

// watch

// functions
async function fetchData() {
  loading.value = true;
  error.value = '';
  const fetchAnalyticsInput: GetAnalyticsDataInput = {
    website_id: websiteId,
  };
  const fetchWebsiteInput: GetWebsiteInput = {
    id: websiteId,
  }

  try {
    const [analyticsDataApi, websiteApi] = await Promise.all([
      $mdninja.getAnalyticsData(fetchAnalyticsInput),
      $mdninja.getWebsite(fetchWebsiteInput),
    ]);
    analyticsData.value = analyticsDataApi;
    website.value = websiteApi;
    analyticsData.value.countries = analyticsData.value.countries.map((item) => {
      item.label = `${getCountryFlag(item.label)} ${item.label}`;
      return item;
    });
    // nextTick(() => initChart());
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}
</script>
