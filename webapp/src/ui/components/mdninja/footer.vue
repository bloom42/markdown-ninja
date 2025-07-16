<template>
  <footer class="bg-white py-10 w-full">
    <div class="mx-auto max-w-7xl overflow-hidden px-6 lg:px-8">
      <nav class="columns-2 sm:flex sm:justify-center space-x-4 sm:space-x-12" aria-label="Footer">
        <span v-for="item in navigation" :key="item.name">
          <a v-if="item.external" :href="item.to" target="_blank" rel="noopener" class="text-sm leading-6 hover:underline">
            {{ item.name }}
          </a>
          <RouterLink v-else :to="item.to" class="text-sm leading-6 hover:underline">
            {{ item.name }}
          </RouterLink>
        </span>
      </nav>
      <!-- <p class="mt-5 text-center text-xs leading-5 text-gray-500">&copy; {{ year }} Markdown Ninja
         <!-- All rights reserved.
      </p> -->
    </div>
  </footer>
</template>

<script lang="ts" setup>
import { useConfig } from '@/app/config';


type NavigationItem = {
  name: string;
  to: string;
  external?: boolean;
}
// props

// events

// composables
const $config = useConfig();

// lifecycle

// variables
// const year = new Date().getFullYear();
const navigation: NavigationItem[] = [
  { name: 'About', to: '/about' },
  { name: 'Terms of Service', to: '/terms' },
  { name: 'Privacy Policy', to: '/privacy' },
  { name: 'Contact', to: '/contact' },
  { name: 'GitHub', to: $config.githubRepository, external: true },
];

// computed

// watch

// functions
</script>
