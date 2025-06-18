<template>
  <div class="flex-1">

    <div class="mt-5" >
      <PContact :website-id="websiteId" />
    </div>

  </div>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router';
import PContact from '@/ui/components/contacts/contact.vue';

// props

// events

// composables
const $route = useRoute();

// lifecycle

// variables
const websiteId = $route.params.website_id as string;

// computed

// watch

// functions

</script>
