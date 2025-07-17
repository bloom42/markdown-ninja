<template>
  <Disclosure as="nav" class="w-full shadow-sm bg-white" v-slot="{ open, close }">

  <header v-if="showAppNavbar" class="shrink-0 bg-white fixed w-full top-0 left-0 z-50 border-b border-gray-200 shadow-2xs">
    <div class="flex h-14 justify-between px-2 sm:px-6 w-full">
      <div class="flex">
        <div class="flex sm:hidden items-center mr-2" v-if="showSidebarButton">
          <sl-tooltip content="Menu">
            <sl-button variant="text" @click="$emit('sidebarButtonClicked')" circle
              class="rounded-full hover:bg-neutral-100">
              <span class="sr-only">Open menu</span>
              <Bars3Icon class="h-6 w-6" aria-hidden="true" />
            </sl-button>
          </sl-tooltip>
        </div>
        <RouterLink to="/organizations" class="flex flex-shrink-0 items-center">
          <img class="h-9 w-auto hidden sm:flex" src="/webapp/markdown_ninja_logo.svg" alt="Markdown Ninja logo" />
          <h1 class="text-base xs:text-xl font-medium  sm:ml-3 inline">
            Markdown Ninja
          </h1>
        </RouterLink>

        <sl-breadcrumb v-if="!$route.path.startsWith('/admin/')" class="hidden xs:inline-flex ml-5">
          <sl-breadcrumb-item v-if="organizationLink" class="">
            <RouterLink  :to="organizationLink.url"
              class="flex flex-shrink-0 items-center text-base font-normal text-(--primary-color) hover:bg-neutral-100 rounded-md px-2.5 py-2">
              {{  organizationLink.name  }}
            </RouterLink>
          </sl-breadcrumb-item>
          <sl-breadcrumb-item v-if="websiteLink">
            <a :href="websiteLink.url" target="_blank" rel="noopener"
              class="flex flex-shrink-0 items-center text-base font-normal text-(--primary-color) hover:bg-neutral-100 rounded-md px-2.5 py-2">
              {{  websiteLink.name  }}
            </a>
          </sl-breadcrumb-item>
        </sl-breadcrumb>

      </div>
      <div class="flex items-center gap-x-2 sm:gap-x-6 font-normal">

        <RouterLink v-if="$store.isAdmin" to="/admin" >
          <sl-button variant="text">
            Admin
          </sl-button>
        </RouterLink>

        <sl-tooltip content="Docs & Help" placement="bottom">
          <a :href="$config.docsUrl" target="_blank" rel="noopener">
            <sl-button variant="text" circle
              class="rounded-full hover:bg-neutral-100 h-10 w-10">
              <QuestionMarkCircleIcon class="flex-auto h-6 w-6" />
            </sl-button>
          </a>
        </sl-tooltip>
        <sl-tooltip content="Account" placement="bottom">
          <a :href="accountUrl">
            <div class="h-10 w-10 relative flex cursor-pointer max-w-xs items-center rounded-full text-center hover:bg-neutral-100">
              <UserIcon class="flex-auto h-6 w-auto"  />
            </div>
          </a>
        </sl-tooltip>
      </div>
    </div>
  </header>

  <header v-else :class="'bg-white shrink-0 fixed w-full top-0 left-0 z-50 h-14 border-b border-gray-200 shadow-2xs'">
    <div class="flex justify-between px-2 sm:px-6 w-full place-content-center h-full items-center">
      <div class="flex">
        <div class="flex sm:hidden items-center" v-if="showSidebarButton">
          <sl-button variant="text" @click="$emit('sidebarButtonClicked'); close()" circle
            class="hover:bg-neutral-200 hover:rounded-full">
            <span class="sr-only">Open menu</span>
            <Bars3Icon class="h-6 w-6 text-gray-900" aria-hidden="true" />
          </sl-button>
        </div>
        <RouterLink to="/" class="flex flex-shrink-0 items-center">
          <img :class="['h-9 w-auto', showSidebarButton ? 'hidden sm:flex': 'flex' ]"
            src="/webapp/markdown_ninja_logo.svg" alt="Markdown Ninja logo" />
          <h1 class="text-base xs:text-xl font-medium ml-1 sm:ml-3">
            Markdown Ninja
          </h1>
        </RouterLink>
      </div>
      <div class="hidden sm:flex gap-x-8 md:gap-x-12 items-center -ml-16">
        <!-- render a RouterLink if URL starts with '/' and a normal link otherwise -->
        <component :is="item.url.startsWith('http') ? 'a' : 'RouterLink'"
          v-for="item in navigation" :key="item.url"
          :to="item.url" :href="item.url" :target="item.url.startsWith('/') ? undefined : '_blank'"
          class="text-sm font-semibold leading-6 hover:underline"
        >
          {{ item.name }}
        </component>
      </div>

      <div class="items-center flex flex-row space-x-2">
        <RouterLink v-if="$store.userId" to="/organizations">
          <sl-button variant="primary" >
            Dashboard
          </sl-button>
        </RouterLink>
        <sl-button v-else variant="primary" @click="signIn()">
          Sign up / Log in
        </sl-button>
        <div class="flex sm:hidden items-center">
          <DisclosureButton as="sl-button" variant="text" circle size="medium"
            class="hover:bg-neutral-200 hover:rounded-full">
            <span class="sr-only">Open menu</span>
            <EllipsisVerticalIcon v-if="!open" class="block h-6 w-6" aria-hidden="true" />
            <XMarkIcon v-else class="block h-6 w-6" aria-hidden="true" />
          </DisclosureButton>
        </div>
      </div>

    </div>
  </header>

  <!-- mobile menu -->
    <DisclosurePanel class="md:hidden mt-14 -mb-15" v-slot="{ close }" >
      <div class="space-y-1 py-2 px-2 rounded-md bg-white">
        <RouterLink v-for="item in navigation" :to="item.url" @click="close()"
          :class="[routeMatchUrl(item.url) ? 'border-[var(--primary-color)]  bg-white brightness-[0.95]' : 'border-transparent text-gray-500', 'hover:bg-white hover:brightness-[0.95] block border-l-4 py-2 pl-3 pr-4 text-base font-medium sm:pl-5 sm:pr-6']">
          {{ item.name }}
      </RouterLink>
      </div>
    </DisclosurePanel>

  </Disclosure>
</template>

<script lang="ts" setup>
import { computed, onBeforeMount, type PropType } from 'vue'
import { Bars3Icon, UserIcon, QuestionMarkCircleIcon, EllipsisVerticalIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { usePingoo } from '@/pingoo/pingoo';
import SlButton from '@shoelace-style/shoelace/dist/components/button/button.js';
import { useStore } from '@/app/store';
import { useRoute } from 'vue-router';
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue'
import SlBreadcrumb from '@shoelace-style/shoelace/dist/components/breadcrumb/breadcrumb.js';
import SlBreadcrumbItem from '@shoelace-style/shoelace/dist/components/breadcrumb-item/breadcrumb-item.js';
import { useMdninja } from '@/api/mdninja';
import SlTooltip from '@shoelace-style/shoelace/dist/components/tooltip/tooltip.js';
import { useConfig } from '@/app/config';


type NavigationLink = {
  name: string,
  url: string,
}

const navigation: NavigationLink[] = [
  { name: 'Home', url: '/' },
  { name: 'Docs', url: 'https://docs.markdown.ninja' },
  { name: 'Pricing', url: '/pricing' },
];




// props
defineProps({
  showSidebarButton: {
    type: Boolean as PropType<boolean>,
    required: false,
    default: false,
  },
})

// events
const emit = defineEmits(['sidebarButtonClicked']);

// composables
const $pingoo = usePingoo();
const $store = useStore();
const $route = useRoute();
const $mdninja = useMdninja();
const $config = useConfig();

// lifecycle
onBeforeMount(() => {
  // If the website is not in store, then we fetch it to be able to display the website link.
  // We wait a little bit for the situation where another component or the current page is already
  // fetching the website.
  if ($route.params.website_id) {
    setTimeout(() => {
      if (!websiteLink.value) {
        $mdninja.getWebsite({ id: $route.params.website_id as string });
      }
    }, 350);
  }
})

// variables
const accountUrl = $pingoo.accountUrl();

// computed
const showAppNavbar = computed(() => {
  return $store.userId && $route.path !== '/'
    && $route.path !== '/pricing' && $route.path !== '/about' && $route.path !== '/terms'
    && $route.path !== '/privacy'   && $route.path !== '/contact'
});

const organizationLink = computed((): NavigationLink | null => {
  if ($route.params.organization_id) {
    const org = $store.organizations.find((org) => org.id === $route.params.organization_id);
    return org ? { name: org.name, url: `/organizations/${org.id}` } : null;
  }
  if ($route.params.website_id) {
    const website = $store.websites.find((site) => site.id === $route.params.website_id);
    if (!website) {
      return null;
    }
    const org = $store.organizations.find((org) => org.id === website.organization_id);
    return org ? { name: org.name, url: `/organizations/${org.id}` } : null;
  }

  return null;
});

const websiteLink = computed((): NavigationLink | null => {
  if ($route.params.website_id) {
    const website = $store.websites.find((website) => website.id === $route.params.website_id);
    return website ? { name: website.primary_domain, url: $mdninja.generateWebsiteUrl(website) } : null;
  }

  return null;
})


// watch

// functions
function signIn() {
  $pingoo.login();
}

function routeMatchUrl(url: string): boolean {
  if (url === "/") {
    return $route.path === "/";
  }

  return $route.path.startsWith(url);
}
</script>
