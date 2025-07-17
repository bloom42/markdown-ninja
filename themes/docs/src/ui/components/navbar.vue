<template>
  <!-- <nav> -->
    <Disclosure as="nav" class="w-full shadow-sm bg-[var(--mdninja-background)]" v-slot="{ open }">
      <div :class="[!hasSidebar && 'max-w-7xl', 'mx-autopx-2 xs:px-4 sm:px-8']">
      <div class="flex h-16 justify-between w-full">
        <div class="flex">
          <div v-if="hasSidebar && website.navigation.secondary.length > 0"class="flex items-center md:hidden">
          <!-- Mobile sidebar button -->
          <button @click="$store.toggleSidebar(true)"
            class="mr-1 sm:mr-2 relative inline-flex items-center justify-center rounded-md p-2 text-gray-400 bg-[var(--mdninja-background)] hover:brightness-[0.95] shadow-none">
            <span class="absolute -inset-0.5" />
            <span class="sr-only">Open sidebar</span>
            <Bars3Icon class="block h-6 w-6" aria-hidden="true" />
          </button>
        </div>

          <!-- Logo and site name -->
          <div class="flex items-center flex-row">
            <PLink href="/" class="flex no-underline items-center">
              <img v-if="$store.website?.logo" :src="$store.website.logo" alt="Website logo" class="inline-flex h-8 w-8" />
              <h1 class="inline-flex ml-1.5 sm:ml-3 font-normal text-base xs:text-lg sm:text-xl text-[var(--mdninja-text)]">
                {{ website.name }}
              </h1>
            </PLink>
          </div>

          <!-- desktop menu -->
          <div class="hidden md:flex md:ml-6 md:space-x-7">
            <template v-for="nav in website.navigation.primary">
              <PLink v-if="nav.url" :href="nav.url"
                :class="[routeMatchUrl(nav.url) ? 'border-[var(--mdninja-accent)]' : 'border-transparent hover:border-gray-300', 'px-1 inline-flex items-center border-b-2 pt-1 text-sm font-medium text-[var(--mdninja-text)] no-underline']">
                {{ nav.label }}
              </PLink>
              <Popover class="relative items-center" v-else>
                <div class="flex items-center h-full">
                  <PopoverButton class="p-0 px-1 inline-flex items-center font-medium text-[var(--mdninja-text)] bg-transparent shadow-none outline-none">
                    <span>{{ nav.label }}</span>
                    <ChevronDownIcon class="size-4.5 pt-0.5" aria-hidden="true" />
                  </PopoverButton>
                </div>


                <transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0 translate-y-1" enter-to-class="opacity-100 translate-y-0" leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100 translate-y-0" leave-to-class="opacity-0 translate-y-1">
                  <PopoverPanel class="absolute left-1/2 z-10 mt-2 flex w-screen max-w-min -translate-x-1/2 px-4">
                    <div class="w-56 shrink rounded-xl bg-(--mdninja-background) p-4 text-sm/6 font-medium shadow-lg ring-1 ring-gray-900/5">
                      <PLink v-for="child in nav.children" :href="child.url!" class="text-[var(--mdninja-text)] no-underline hover:underline block p-2">
                        {{ child.label }}
                      </PLink>
                    </div>
                  </PopoverPanel>
                </transition>
              </Popover>
            </template>
          </div>

        </div>

        <div class="flex items-center content-center space-x-3 sm:space-x-4">
          <button type="button" @click="openSearchbar()"
            class="p-2.5 inline-flex rounded-full bg-[var(--mdninja-background)] text-gray-500 shadow-none hover:bg-[var(--mdninja-background)] hover:brightness-[0.95]">
            <!-- <span class="absolute -inset-1.5" /> -->
            <span class="sr-only">Search</span>
            <MagnifyingGlassIcon class="size-5" aria-hidden="true" />
          </button>


          <PLink v-if="subscribeButton && $store.contact" href="/account">
            <PButton type="button" class="inline-flex items-center gap-x-1.5 rounded-md px-3 py-2 text-sm font-semibold shadow-xs">
              Account
            </PButton>
          </PLink>
          <PLink v-else-if="subscribeButton && !$store.contact" href="/subscribe">
            <PButton type="button" class="inline-flex items-center gap-x-1.5 rounded-md px-3 py-2 text-sm font-semibold shadow-xs">
              Subscribe
            </PButton>
          </PLink>

          <div v-if="website.navigation.primary.length > 0"class="flex items-center md:hidden">
            <!-- Mobile menu button -->
            <DisclosureButton class="relative inline-flex items-center justify-center rounded-md p-2 text-gray-500 bg-[var(--mdninja-background)] hover:brightness-[0.95] shadow-none">
              <!-- focus:outline-none focus:ring-2 focus:ring-inset focus:ring-sky-500"> -->
              <span class="absolute -inset-0.5" />
              <span class="sr-only">Open menu</span>
              <EllipsisVerticalIcon v-if="!open" class="block h-6 w-6" aria-hidden="true" />
              <XMarkIcon v-else class="block h-6 w-6" aria-hidden="true" />
            </DisclosureButton>
          </div>
        </div>

      </div>
    </div>

    <!-- mobile menu -->
    <DisclosurePanel class="md:hidden" v-slot="{ close }">
      <div class="space-y-1 pb-3 pt-2 bg-[var(--mdninja-background)]">
        <!-- DisclosureButton -->
        <template v-for="item in website.navigation.primary">
          <PLink v-if="item.url" :href="item.url" @click="close()"
              :class="[routeMatchUrl(item.url) ? 'border-[var(--mdninja-accent)]  bg-[var(--mdninja-background)] brightness-[0.95]' : 'border-transparent', ' text-(--mdninja-text) hover:bg-[var(--mdninja-background)] hover:brightness-[0.95] block border-l-4 py-2 pl-3 pr-4 text-base font-medium sm:pl-5 sm:pr-6 no-underline']">
              {{ item.label }}
          </PLink>
          <div v-else class="block py-2 pl-3 pr-4 text-base font-medium  text-(--mdninja-text)">
            <span>{{ item.label }}</span>
            <ChevronDownIcon class="size-4.5 inline-flex ml-3" />
            <div class="space-y-1 mt-3 px-3">
              <PLink v-for="subItem in item.children" :href="subItem.url!" @click="close()"
                :class="[routeMatchUrl(subItem.url!) ? 'border-[var(--mdninja-accent)]  bg-[var(--mdninja-background)] brightness-[0.95]' : 'border-transparent', 'sm:px-8 text-(--mdninja-text) hover:bg-[var(--mdninja-background)] hover:brightness-[0.95] block border-l-4 py-2 pl-3 pr-4 text-base font-medium sm:pl-5 sm:pr-6 no-underline']">
                {{ subItem.label }}
              </PLink>
            </div>
          </div>
        </template>
      </div>
    </DisclosurePanel>
  </Disclosure>
</template>

<script lang="ts" setup>
import { type PropType } from 'vue';
import type { Website } from "@/app/model";
import PLink from '@/ui/components/p_link.vue';
import { useRoute } from 'vue-router';
import { useStore } from '@/app/store';
import { Disclosure, DisclosureButton, DisclosurePanel, Popover, PopoverButton, PopoverPanel } from '@headlessui/vue'
import { Bars3Icon, XMarkIcon, EllipsisVerticalIcon, MagnifyingGlassIcon, ChevronDownIcon } from '@heroicons/vue/24/outline'
import PButton from '@/ui/components/p_button.vue';


// props
defineProps({
  website: {
    type: Object as PropType<Website>,
    required: true,
  },
  hasSidebar: {
    type: Boolean as PropType<boolean>,
    required: false,
    default: false,
  },
  subscribeButton: {
    type: Boolean as PropType<boolean>,
    required: false,
    default: false,
  }
});

// events

// composables
const $route = useRoute();
const $store = useStore();

// lifecycle

// variables

// computed

// watch

// functions
function routeMatchUrl(url: string): boolean {
  if (url === "/") {
    return $route.path === "/";
  }

  return $route.path.startsWith(url);
}

function openSearchbar() {
  $store.setShowSearchbar(true);
}
</script>
