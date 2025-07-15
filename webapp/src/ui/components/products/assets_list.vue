<template>
  <div class="-my-2 overflow-x-auto min-w-full">
    <div class="py-2 align-middle inline-block min-w-full">
      <div class="overflow-hidden border border-gray-300 sm:rounded-lg">
        <table class="table min-w-full divide-y divide-gray-200">
          <thead class="table-header-group bg-gray-50">
            <tr class="max-w-0">
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Path
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Size
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="min-w-full bg-white divide-y divide-gray-200">
            <tr v-for="asset in assets" :key="asset.id">
              <td class="px-6 py-4 whitespace-nowrap max-w-0">
                assets/{{ asset.name }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap max-w-0">
                {{ fileSize(asset.size) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap max-w-0 w-1/5">
                <div class="flex flex-row space-x-3">
                  <a :href="$mdninja.generateAssetIdUrl(website, asset.id, true)" download>
                    <sl-button variant="neutral" circle>
                      <ArrowDownCircleIcon class="h-5 w-5" aria-hidden="true" />
                    </sl-button>
                  </a>

                  <sl-button variant="neutral" circle @click="onDeleteClicked(asset)">
                    <TrashIcon class="h-5 w-5" aria-hidden="true" />
                  </sl-button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { type PropType } from 'vue';
import { type Asset, type Website } from '@/api/model';
import { TrashIcon, ArrowDownCircleIcon } from '@heroicons/vue/24/outline'
import { useMdninja } from '@/api/mdninja';
import SlButton from '@shoelace-style/shoelace/dist/components/button/button.js';
import fileSize from '@/filters/filesize';

// props
defineProps({
  assets: {
    type: Array as PropType<Asset[]>,
    required: true,
  },
  website: {
    type: Object as PropType<Website>,
    required: true,
  },
})

// events
const $emit = defineEmits(['delete', 'update']);

// composables
const $mdninja = useMdninja();

// lifecycle

// variables

// computed

// watch

// functions
function onDeleteClicked(asset: Asset) {
  $emit('delete', asset);
}
</script>
