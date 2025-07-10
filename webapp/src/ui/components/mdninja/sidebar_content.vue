<template>
  <nav class="flex flex-1 flex-col">
    <ul role="list" class="flex flex-1 flex-col gap-y-7">
      <li>
        <ul role="list" class="-mx-2 space-y-1">
          <li v-for="item in navigation" :key="item.name">
            <RouterLink v-if="!item.children" :to="item.to!"
              :class="[isCurrentPage(item.to!) ? 'bg-gray-100 text-(--primary-color)' : 'text-neutral-600 hover:text-(--primary-color) hover:bg-gray-50', 'cursor-pointer group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-medium']">
              <component v-if="item.icon" :is="item.icon" :class="[isCurrentPage(item.to) ? 'text-(--primary-color)' : 'text-gray-400 group-hover:text-(--primary-color)', 'h-6 w-6 shrink-0']" aria-hidden="true" />
              {{ item.name }}
            </RouterLink>

            <Disclosure as="div" v-else v-slot="{ open }">
              <DisclosureButton :class="[isCurrentPage(item.to) ? 'bg-gray-100' : 'hover:bg-gray-50', 'hover:text-(--primary-color) flex w-full items-center gap-x-3 rounded-md p-2 text-left text-sm/6 font-medium text-gray-600']">
                <component v-if="item.icon" :is="item.icon" :class="[isCurrentPage(item.to) ? 'text-(--primary-color)' : 'text-gray-400 group-hover:text-(--primary-color)', 'h-6 w-6 shrink-0']" />
                {{ item.name }}
                <ChevronRightIcon :class="[open ? 'rotate-90 text-gray-500' : 'text-gray-400', 'group-hover:text-(--primary-color) ml-auto size-5 shrink-0']" aria-hidden="true" />
              </DisclosureButton>

              <DisclosurePanel as="ul" class="mt-1 px-2">
                <li v-for="subItem in item.children" :key="subItem.name">
                  <RouterLink :to="subItem.to!">
                    <span :class="[isCurrentPage(subItem.to) ? 'bg-gray-100 text-(--primary-color)' : 'hover:bg-gray-50', 'block rounded-md py-2 pl-6 pr-2 text-sm/6 text-gray-600']">
                      <component v-if="subItem.icon" :is="subItem.icon" :class="[isCurrentPage(subItem.to) ? 'text-(--primary-color)' : 'text-gray-400 group-hover:text-(--primary-color)', 'h-6 w-6 shrink-0 inline mr-2']" aria-hidden="true" />
                      {{ subItem.name }}
                    </span>
                  </RouterLink>
                </li>
              </DisclosurePanel>
            </Disclosure>
          </li>

        </ul>
      </li>
    </ul>
  </nav>
</template>

<script lang="ts" setup>
import { onBeforeMount, ref, watch, type Ref, markRaw } from 'vue';
import { useRoute } from 'vue-router';
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue'
import {
  DocumentTextIcon,
  PhotoIcon,
  TagIcon,
  CodeBracketIcon,
  ArrowsRightLeftIcon,
  UserGroupIcon,
  ListBulletIcon,
  ReceiptPercentIcon,
  ShoppingCartIcon,
  EnvelopeIcon,
  Cog6ToothIcon,
  MapIcon,
  GlobeAltIcon,
  CreditCardIcon,
  ClockIcon,
  HomeIcon,
  ArrowUturnLeftIcon,
  PresentationChartLineIcon,
  SparklesIcon,
  ShieldCheckIcon,
} from '@heroicons/vue/24/outline';
import { ChevronRightIcon } from '@heroicons/vue/20/solid'
import FeatherIcon from '@/ui/icons/feather.vue';
import SendIcon from '@/ui/icons/send.vue';
import LettersLowercaseIcon from '@/ui/icons/letters_lowercase.vue';
import KeyIcon from '@/ui/icons/key.vue';
import SettingsIcon from '@/ui/icons/settings.vue';
import BookUserIcon from '@/ui/icons/book_user.vue';
import PaletteIcon from '@/ui/icons/palette.vue';
import { useStore } from '@/app/store';

type NavigationItem = {
  name: string;
  icon?: any;
  to?: string;
  children?: NavigationItem[],
};


// props
defineExpose({
  open,
});

// events

// composables
const $route = useRoute();
const $store = useStore();

// lifecycle
onBeforeMount(() => onRouteChanged());


// variables
let navigation: Ref<NavigationItem[]> = ref([]);

// computed

// watch
watch($route, () => onRouteChanged(), { deep: true });

// functions
function onRouteChanged() {
  const websiteId = $route.params.website_id;
  if ($route.path.startsWith('/admin')) {
    navigation.value = [
      { name: 'Admin Home', to: '/admin', icon: HomeIcon },
      { name: 'Organizations', to: '/admin/organizations', icon: UserGroupIcon },
      { name: 'Websites', to: '/admin/websites', icon: GlobeAltIcon },
      { name: 'Background Jobs', to: '/admin/queue', icon: ClockIcon },
      { name: 'TLS Certificates', to: '/admin/tls', icon: ShieldCheckIcon },
    ];
  } else if ($route.path.startsWith('/websites')) {
    const nav = [
      { name: 'Home', to: `/websites/${websiteId}`, icon: PresentationChartLineIcon },
      { name: 'Posts', to: `/websites/${websiteId}/posts`, icon: markRaw(FeatherIcon) },
      { name: 'Pages', to: `/websites/${websiteId}/pages`, icon: DocumentTextIcon },
      { name: 'Assets & Media', to: `/websites/${websiteId}/assets`, icon: PhotoIcon },
      { name: 'Newsletters', to: `/websites/${websiteId}/newsletters`, icon: markRaw(SendIcon) },
      { name: 'Contacts', to: `/websites/${websiteId}/contacts`, icon: markRaw(BookUserIcon) },
      {
        name: 'Settings',
        icon: Cog6ToothIcon,
        children: [
          { name: 'General', to: `/websites/${websiteId}/settings`, icon: Cog6ToothIcon },
          { name: 'Design & Branding', to: `/websites/${websiteId}/settings/design`, icon: markRaw(PaletteIcon) },
          { name: 'Emails', to: `/websites/${websiteId}/settings/emails`, icon: EnvelopeIcon },
          { name: 'Code', to: `/websites/${websiteId}/settings/code`, icon: CodeBracketIcon },
          { name: 'Tags', to: `/websites/${websiteId}/tags`, icon: TagIcon },
          { name: 'Redirects', to: `/websites/${websiteId}/redirects`, icon: ArrowsRightLeftIcon },
          { name: 'Navigation', to: `/websites/${websiteId}/navigation`, icon: MapIcon },
          { name: 'Domains', to: `/websites/${websiteId}/settings/domains`, icon: markRaw(LettersLowercaseIcon) },
        ],
      },
    ];

    if ($store.isAdmin) {
      nav.push({
        name: 'Beta',
        icon: SparklesIcon,
        children: [
          { name: 'Products', to: `/websites/${websiteId}/products`, icon: ShoppingCartIcon },
          { name: 'Coupons', to: `/websites/${websiteId}/coupons`, icon: ReceiptPercentIcon },
          { name: 'Orders', to: `/websites/${websiteId}/orders`, icon: ListBulletIcon },
          { name: 'Refunds', to: `/websites/${websiteId}/refunds`, icon: ArrowUturnLeftIcon },
          { name: 'Snippets', to: `/websites/${websiteId}/snippets`, icon: CodeBracketIcon }
        ],
      });
    }

    navigation.value = nav;
  } else if ($route.params.organization_id) {
    const organizationId = $route.params.organization_id;
    navigation.value = [
      { name: 'Websites', to: `/organizations/${organizationId}/websites`, icon: GlobeAltIcon },
      { name: 'Staffs', to: `/organizations/${organizationId}/staffs`, icon: UserGroupIcon },
      { name: 'API Keys', to: `/organizations/${organizationId}/api`, icon: markRaw(KeyIcon) },
      { name: 'Billing', to: `/organizations/${organizationId}/billing`, icon: CreditCardIcon },
      { name: 'Settings', to: `/organizations/${organizationId}/settings`, icon: markRaw(SettingsIcon) },
    ];
  } else if ($route.path.startsWith('/organizations')) {
    navigation.value = [
      { name: 'Organizations', to: '/organizations', icon: UserGroupIcon },
      { name: 'Invitations', to: '/organizations/invitations', icon: EnvelopeIcon },
    ];
  } else if ($route.path.startsWith('/docs')) {
    navigation.value = [
      // { name: 'Introduction', to: '/docs' },
      { name: 'CLI & Git integration', to: '/docs/cli' },
      { name: 'API & Headless CMS', to: '/docs/api' },
      { name: 'Feedback & Open Source', to: '/docs/feedback' },
      // { name: 'Newsletter', to: '/docs/newsletter' },
    ];
  }  else {
    navigation.value = [];
  }
}

function isCurrentPage(path?: string): boolean {
  return $route.path === path;
}
</script>
