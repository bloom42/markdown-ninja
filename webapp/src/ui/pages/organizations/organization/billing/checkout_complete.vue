<template>
  <div>
    <div class="rounded-md bg-red-50 p-2 mb-3 mt-10" v-if="error">
      <div class="flex">
        <div class="ml-3">
          <p class="text-sm text-red-700">
            {{ error }}
          </p>
        </div>
      </div>
    </div>

    <div v-if="loading" class="flex flex-col items-center">
      <div class="flex">
        <svg class="animate-spin -ml-1 mr-3 h-12 w-12 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-50" cx="2" cy="2" r="2" stroke="currentColor" stroke-width="2"></circle>
          <path class="opacity-75" fill="#424242" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>

      <div class="flex mt-5">
        Finalizing your subscription. Please do not reload or change the page.
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useMdninja } from '@/api/mdninja';
import { onBeforeMount, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';


// This page waits for the stripe webhook to arrive in the backend so the plan is updated.
// This is done by polling the backend and checking that the plan has been updated.
// TODO: handle error and show a message if the user is waiting for too long.

// props

// events

// composables
const $mdninja = useMdninja();
const $route = useRoute();
const $router = useRouter();

// lifecycle
onBeforeMount(() => {
  syncStripe();
});

// variables
const organizationId = $route.params.organization_id as string;
// const plan = $route.query.plan as string | undefined;

let error = ref('');
let loading = ref(true);

// computed

// watch

// functions
async function syncStripe() {
  try {
    await $mdninja.organizationSyncStripe(organizationId);
    if ($route.query.redirect_to && ($route.query.redirect_to as string).startsWith('/')) {
      $router.push($route.query.redirect_to as string);
      return ;
    }

    $router.push(`/organizations/${organizationId}`);
    // const res = await $mdninja.getOrganization(input);
    // if ((plan && res.plan === plan)
    //   || (!plan && res.plan !== 'free')) {
    //   $router.push(`/organizations/${res.id}/billing`);
    //   return
    // }
  } catch (err: any) {
    // error.value = err.message;
    error.value = 'Internal Error. Please reload the page and contact support if the problem persists'
    loading.value = false;
    // clearInterval(intervalId);
  }
}
</script>
