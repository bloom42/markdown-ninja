<template>
  <div class="flex flex-col">
    <TlsCertificatesList :tls-certificates="tlsCertificates" />
  </div>
</template>

<script lang="ts" setup>
import { useMdninja } from '@/api/mdninja';
import type { TlsCertificates } from '@/api/model';
import { onBeforeMount, ref, type Ref } from 'vue';
import TlsCertificatesList from '@/ui/components/kernel/tls_certificates_list.vue';

// props

// events

// composables
const $mdninja = useMdninja();

// lifecycle
onBeforeMount(() => fetchData());


// variables
let tlsCertificates: Ref<TlsCertificates[]> = ref([]);
let loading = ref(false);
let error = ref('');
// computed

// watch

// functions
async function fetchData() {
  loading.value = true;
  error.value = '';

  try {
    const res = await $mdninja.listTlsCertificates();
    tlsCertificates.value = res.data;
    console.log(tlsCertificates.value);
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}
</script>
