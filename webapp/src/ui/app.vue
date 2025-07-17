<template>
  <div class="border border-red-600 rounded-md bg-red-50 p-4 my-5 mx-5 h-full" v-if="$store.appLoadingError">
    <div class="flex">
      <div class="ml-3">
        <p class="text-sm text-red-700">
          {{ $store.appLoadingError }}
        </p>
      </div>
    </div>
  </div>


  <div v-if="organizationPaymentDue" class="flex z-[100] w-full fixed items-center gap-x-6 bg-red-600 px-6 py-2.5 sm:px-3.5 sm:before:flex-1">
    <p class="text-sm/6 text-center w-full text-white">
      <RouterLink :to="`/organizations/${organizationPaymentDue.id}/billing`">
        <strong class="font-semibold">
          A payment has failed for your organization {{ organizationPaymentDue.name }}. Click here to regularize the situation
          <span aria-hidden="true">&rarr;</span>
        </strong>
      </RouterLink>
    </p>
  </div>

  <div class="h-full" v-if="!$store.appLoadingError">
    <Navbar :show-sidebar-button="showSidebar" @sidebar-button-clicked="openSidebar()" />
    <Sidebar v-if="showSidebar"  ref="sidebar" >
      <SidebarContent />
    </Sidebar>

    <main class="mt-16 h-full">
      <div :class="['h-full', showSidebar ? 'sm:ml-64 px-4': '']">
        <RouterView />
      </div>
    </main>

    <Footer v-if="showFooter" />
  </div>
</template>

<script lang="ts" setup>
import { useStore } from '@/app/store';
// import Sidenav from '@/ui/components/mdninja/sidenav.vue';
import Sidebar from '@/ui/components/mdninja/sidebar.vue';
import Navbar from '@/ui/components/mdninja/navbar.vue';
import { computed, onBeforeMount, ref, useTemplateRef, watch, type Ref } from 'vue';
import { useRoute } from 'vue-router';
import Footer from '@/ui/components/mdninja/footer.vue';
import type { Organization } from '@/api/model';
import SidebarContent from './components/mdninja/sidebar_content.vue';

// props

// events

// composables
const $store = useStore();
const $route = useRoute();

// lifecycle
onBeforeMount(() => {
  setShowSidebar();
  checkPaymentDueForOrganization();
});

// variables
let showSidebar = ref(false);
let organizationPaymentDue: Ref<Organization | null> = ref(null);
  const sidebar = useTemplateRef('sidebar');

// computed
const showFooter = computed(() => {
  return $route.path === '/'
    || $route.path === '/pricing' || $route.path === '/privacy' || $route.path === '/terms'
    || $route.path === '/about' || $route.path === '/contact'
});


// watch
watch($route, () => setShowSidebar(), { deep: true });
watch(() => $store.organizations, () => checkPaymentDueForOrganization(), { deep: true });

// functions
function setShowSidebar() {
  if (($route.path.startsWith('/websites/') && $route.path !== '/websites/new')
    || $route.path.startsWith('/organizations')
    || $route.params.organization_id
    || $route.path.startsWith('/admin')
  ) {
      showSidebar.value = true;
  } else {
    showSidebar.value = false;
  }
}

function checkPaymentDueForOrganization(){
  for (let org of $store.organizations) {
    if (org.payment_due) {
      organizationPaymentDue.value = org;
      return;
    }
  }

  organizationPaymentDue.value = null;
}

function openSidebar() {
  sidebar.value?.open();
}
</script>
