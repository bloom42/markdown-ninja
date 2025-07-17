<template>
  <div class="flex">

    <div class="flex w-full" v-if="adminStatistics.length !== 0">
      <div class="w-full mt-2 rounded-md border border-gray-900/10">
        <dl class="mx-auto grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:px-2 xl:px-0">
          <div v-for="stat in adminStatistics" :key="stat.name" class="sm:border-l flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1 border-t border-gray-900/5 px-4 py-5 sm:px-6 lg:border-t-0 xl:px-8">
            <dt class="text-md">{{ stat.name }}</dt>
            <dd class="w-full flex-none text-3xl font-medium tracking-tight">
              {{ stat.value }}
            </dd>
          </div>
        </dl>
      </div>
    </div>

  </div>
</template>

<script lang="ts" setup>
import { useMdninja } from '@/api/mdninja';
import { onBeforeMount, ref, type Ref } from 'vue';

type Stat = {
  name: string;
  value: number;
}

// props

// events

// composables
const $mdninja = useMdninja();

// lifecycle
onBeforeMount(() => fetchData());


// variables
let adminStatistics: Ref<Stat[]> = ref([]);
let loading = ref(false);
let error = ref('');
// computed

// watch

// functions
async function fetchData() {
  loading.value = true;
  error.value = '';

  try {
    const [organizationsStatistics, websitesStatistics] = await Promise.all([
      $mdninja.getOrganizationsAdminStatistics(),
      $mdninja.getWebsitessAdminStatistics(),
    ]);
    adminStatistics.value = [
      {
        name: 'Websites',
        value: websitesStatistics.websites,
      },
      {
        name: 'Organizations',
        value: organizationsStatistics.organizations,
      },
      {
        name: 'Paying Organizations',
        value: organizationsStatistics.paying_organizations,
      },
    ]
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}
</script>
