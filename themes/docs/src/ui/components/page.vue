<template>
  <div>
    <div v-if="page.path !== '/'">
      <h1 class="font-semibold text-3xl">{{  page.title }}</h1>
      <hr />
    </div>

    <div v-html="page.body" />


    <div v-if="page.tags.length !== 0" class="my-5">
      <hr class="mb-5" />

      <strong class="mx-2">Tags:&nbsp;</strong>
      <!-- <span  class=""> -->
        <!-- <span v-if="$index !== 0">, </span> -->
        <RouterLink  v-for="(tag, _index) in page.tags" :key="tag.name" :to="tagUrl(tag)" class="mr-2">
          <span class="inline-flex items-center gap-x-1.5 rounded-md px-2 py-1 font-medium text-[var(--mdninja-accent)] ring-1 ring-inset ring-gray-200 hover:ring-[var(--mdninja-accent)]">
            {{ tag.name }}
          </span>
        </RouterLink>
      <!-- </span> -->
    </div>

  </div>
</template>

<script lang="ts" setup>
import type { Page, Tag } from '@/app/model';
import type { PropType } from 'vue';

// props
defineProps({
  page: {
    type: Object as PropType<Page>,
    required: true,
  },
});

// events

// composables

// lifecycle

// variables

// computed

// watch

// functions
function tagUrl(tag: Tag) {
  return `/tags/${tag.name}`;
}
</script>
