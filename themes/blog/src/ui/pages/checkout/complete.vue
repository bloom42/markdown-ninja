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
        Completing your order. Please do not reload or change page.
      </div>

    </div>

    <div v-if="success">
      <p>
        Thank you! <br />
        You can now access your products by logging into your account. <br />

        <RouterLink to="/account">
          <span class="text-xl">Go to my account</span>
        </RouterLink>
      </p>

    </div>
  </div>
</template>

<script lang="ts" setup>
import { useStore } from '@/app/store';
import type { CompleteOrderInput } from '@/app/model';
import { onBeforeMount, ref } from 'vue';
import { useRoute } from 'vue-router';
import { completeOrder } from '@/app/mdninja';

// props

// events

// composables
const $route = useRoute();
const $store = useStore();

// lifecycle
onBeforeMount(() => {
  $store.setLoading(false);
  callCompleteOrder();
});

// variables
let success = ref(false);
let error = ref('');
let loading = ref(false);

// computed

// watch

// functions
async function callCompleteOrder() {
  loading.value = true;
  error.value = '';
  const completeOrderInput: CompleteOrderInput = {
    order_id: $route.params.order_id as string,
  };

  try {
    await completeOrder(completeOrderInput);
    success.value = true;
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

</script>
