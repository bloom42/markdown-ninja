<template>
  <div class="-my-2 overflow-x-auto min-w-full">
    <div class="py-2 align-middle inline-block min-w-full">
      <div class="overflow-hidden border border-gray-300 sm:rounded-lg">
        <table class="table min-w-full divide-y divide-gray-200">
          <thead class="table-header-group bg-gray-50">
            <tr class="max-w-0">
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Name
              </th>
              <th scope="col" class="px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">
                Size
              </th>
              <th scope="col" class="px-6 py-3 text-xs font-medium text-gray-500 uppercase">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="min-w-full bg-white divide-y divide-gray-200">
            <tr v-for="asset in assets" :key="asset.id">
              <td @dblclick="emitAssetDbClicked(asset)" class="px-6 py-4 whitespace-nowrap max-w-0 w-full cursor-pointer" >
                <div class="flex items-center">
                  <FolderIcon v-if="asset.type === AssetType.Folder" class="h-5 w-5" aria-hidden="true" />
                  <MusicalNoteIcon v-else-if="asset.type === AssetType.Audio" class="h-5 w-5" aria-hidden="true" />
                  <PhotoIcon v-else-if="asset.type === AssetType.Image" class="h-5 w-5" aria-hidden="true" />
                  <FilmIcon v-else-if="asset.type === AssetType.Video" class="h-5 w-5" aria-hidden="true" />
                  <DocumentIcon v-else class="h-5 w-5" aria-hidden="true" />
                  <span class="ml-2 text-md font-medium text-gray-900 truncate"></span>{{ asset.name }}
                </div>
              </td>
              <td class="px-4 py-4 whitespace-nowrap min-w-10">
                {{ fileSize(asset.size) }}
              </td>
              <td class="px-4 py-4 whitespace-nowrap space-x-4 flex flex-row">
                <sl-tooltip content="Copy path to clipboard" placement="bottom">
                  <sl-button v-if="asset.type !== AssetType.Folder" variant="neutral" @click="copyPathToClipboard(asset)" circle
                    class="mr-2">
                    <Square2StackIcon class="h-5 w-5" aria-hidden="true" />
                  </sl-button>
                </sl-tooltip>

                <sl-tooltip content="Delete" placement="bottom">
                  <sl-button variant="neutral" @click="onDeleteFileClicked(asset)" circle
                    :class="[asset.type === AssetType.Folder ? 'ml-12' : '']">
                    <TrashIcon class="h-5 w-5" aria-hidden="true" />
                  </sl-button>
                </sl-tooltip>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { AssetType, type Asset } from '@/api/model'
import { type PropType } from 'vue'
import { TrashIcon, Square2StackIcon } from '@heroicons/vue/24/outline'
import { FolderIcon, MusicalNoteIcon, DocumentIcon, PhotoIcon, FilmIcon } from '@heroicons/vue/24/outline'
import SlButton from '@shoelace-style/shoelace/dist/components/button/button.js';
import SlTooltip from '@shoelace-style/shoelace/dist/components/tooltip/tooltip.js';
import fileSize from '@/filters/filesize';


// props
defineProps({
  assets: {
    type: Array as PropType<Asset[]>,
    required: true,
  },
});

// events
const $emit = defineEmits(['delete', 'assetdbclicked']);

// composables

// lifecycle

// variables

// computed

// watch

// functions
function onDeleteFileClicked(asset: Asset) {
  $emit('delete', asset.id);
}

function emitAssetDbClicked(asset: Asset) {
  $emit('assetdbclicked', asset);
}

function copyPathToClipboard(asset: Asset) {
  navigator.clipboard.writeText(`${asset.folder}/${asset.name}`);
}
</script>
