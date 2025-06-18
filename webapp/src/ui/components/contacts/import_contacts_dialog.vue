<template>
  <sl-dialog :open="model" @sl-request-close="model = false" label="Import Contacts">
    <div class="rounded-md bg-red-50 p-4 mb-3" v-if="error">
      <div class="flex">
        <div class="ml-3">
          <p class="text-sm text-red-700">
            {{ error }}
          </p>
        </div>
      </div>
    </div>

    <div class="flex-col mt-2">
      <sl-textarea label="Contacts" :value="contacts" @input="contacts = $event.target.value" rows="10"
        :placeholder="contactsPlaceholder" />
    </div>


    <div slot="footer" class="mt-6 flex flex-row space-x-3 place-content-end">
      <sl-button outline @click="close()">
        Cancel
      </sl-button>
      <sl-button variant="primary" :loading="loading" @click="onImportContactsClicked()">
        Import Contacts
      </sl-button>
    </div>

  </sl-dialog>
</template>

<script lang="ts" setup>
import { ref, type PropType } from 'vue';
import { type ImportContactsInput } from '@/api/model';
import { importContacts } from '@/api/mdninja';
import SlButton from '@shoelace-style/shoelace/dist/components/button/button.js';
import SlTextarea from '@shoelace-style/shoelace/dist/components/textarea/textarea.js';
import SlDialog from '@shoelace-style/shoelace/dist/components/dialog/dialog.js';

// props
const model = defineModel({
  type: Boolean as PropType<boolean>,
  required: true,
});

const props = defineProps({
  websiteId: {
    type: String as PropType<string>,
    required: true,
  },
});

// events
const $emit = defineEmits(['update:modelValue', 'imported']);

// composables

// lifecycle

// variables
const contactsPlaceholder = `email,name,country,subscribed_at
email1@email.com,name,XX,2023-01-01T01:01:01Z
email2@email.com,Name2,FR,
email3@email.com,Name3,US,2023-01-03T01:01:01Z
...
`

let error = ref('');
let loading = ref(false);
let contacts = ref('');

// computed

// watch

// functions
function close() {
  model.value = false;
  resetValues();
}

function resetValues() {
  contacts.value = '';
}

async function onImportContactsClicked() {
  loading.value = true;
  error.value = '';
  const input: ImportContactsInput = {
    contacts: contacts.value,
    website_id: props.websiteId,
  };

  try {
    await importContacts(input);
    $emit('imported');
    close();
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}
</script>

