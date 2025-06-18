<template>
  <div class="flex flex-col">
    <div class="flex flex-row justify-between items-center">
      <div class="flex">
        <div class="flex">
          <RouterLink :to="backRoute">
            <sl-button outline>
                Back
            </sl-button>
          </RouterLink>
        </div>

        <div class="flex ml-5">
          <sl-button v-if="contact" variant="primary" @click="updateContact()" :loading="loading">
              Update
          </sl-button>
          <sl-button v-else variant="primary" @click="createContact()" :loading="loading">
              Create
          </sl-button>
        </div>
      </div>

      <div v-if="contact" class="flex ml-5">
          <Menu as="div" class="relative inline-block text-left">
            <div>
              <MenuButton class="inline-flex w-full justify-center gap-x-1.5 rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-xs ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
                <!-- Options -->
                <EllipsisVerticalIcon class="h-5 w-5 text-gray-700" aria-hidden="true" />
              </MenuButton>
            </div>

            <transition enter-active-class="transition ease-out duration-100" enter-from-class="transform opacity-0 scale-95" enter-to-class="transform opacity-100 scale-100" leave-active-class="transition ease-in duration-75" leave-from-class="transform opacity-100 scale-100" leave-to-class="transform opacity-0 scale-95">
              <MenuItems class="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-gray-300 focus:outline-none">
                <div class="py-1">
                  <MenuItem v-if="blocked" v-slot="{ active }">
                    <span @click="blockOrUnblockContact()"
                      :class="[active ? 'bg-neutral-100 text-gray-900' : 'text-gray-700', 'cursor-pointer block px-4 py-2 text-sm']">
                      Unblock Contact
                    </span>
                  </MenuItem>
                  <MenuItem v-else v-slot="{ active }">
                    <span @click="blockOrUnblockContact()"
                      :class="[active ? 'bg-neutral-100 text-gray-900' : 'text-gray-700', 'cursor-pointer block px-4 py-2 text-sm']">
                      Block Contact
                    </span>
                  </MenuItem>
                  <MenuItem v-slot="{ active }">
                    <span @click="openDeleteContactDialog"
                      :class="[active ? 'bg-neutral-100 text-gray-900' : 'text-gray-700', 'cursor-pointer block px-4 py-2 text-sm']">
                      Delete Contact
                    </span>
                  </MenuItem>
                </div>
              </MenuItems>
            </transition>
          </Menu>
        </div>
    </div>

    <div class="rounded-md bg-red-50 p-4 my-5" v-if="error">
      <div class="flex">
        <div class="ml-3">
          <p class="text-sm text-red-700">
            {{ error }}
          </p>
        </div>
      </div>
    </div>

    <div class="flex flex-col mt-5 w-full">

      <sl-input :value="name" @input="name = $event.target.value" type="text"
          :readonly="loading" :disabled="loading" placeholder="Name" label="Name" />

      <sl-input :value="email" @input="email = $event.target.value" type="email"
        :readonly="loading" :disabled="loading" placeholder="Email" label="Email"
        class="mt-5" />

      <div class="flex mt-5" v-if="contact">
        <b>Country</b>: {{ countryName(contact.country) }} ({{ contact.country }})
      </div>

      <div v-if="blocked" class="flex mt-5">
        <div>
          <span  class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-red-100 text-red-800">
            Blocked
          </span>
          <span class="ml-2">{{ date(contact!.blocked_at!) }}</span>
        </div>
      </div>

    </div>



    <div  v-if="contact" class="flex flex-col mt-8 w-full">
      <div class="flex flex-col">
        <div class="flex mb-5">
          <h1 class="text-xl font-extrabold text-gray-900">Marketing</h1>
        </div>

        <sl-switch :checked="subscribedToNewsletter" @sl-change="subscribedToNewsletter = $event.target.checked">
          Subscribed to newsletter
        </sl-switch>
      </div>
      <!-- End of marketing -->


      <div class="flex flex-col mt-10">
        <div class="flex">
          <h1 class="text-xl font-extrabold text-gray-900">Products</h1>
        </div>

        <div class="flex flex-col">
          <div class="overflow-x-auto min-w-full">
            <div class="py-2 align-middle inline-block min-w-full">
              <div class="overflow-hidden border border-gray-300 sm:rounded-lg">
                <table class="min-w-full divide-y divide-gray-200">
                  <thead class="bg-gray-50">
                    <tr class="max-w-0">
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Name
                      </th>
                    </tr>
                  </thead>
                  <tbody class="min-w-full bg-white divide-y divide-gray-200">
                    <RouterLink :to="productUrl(product)" v-for="product in products" :key="product.id"
                        class="table-row cursor-pointer min-w-full">
                      <td class="px-6 py-4 whitespace-nowrap max-w-0 w-2/5">
                        <div class="text-md font-medium text-gray-900 truncate">
                          {{ product.name }}
                        </div>
                      </td>
                    </RouterLink>
                  </tbody>
                </table>
              </div>
            </div>
          </div>


        </div>
      </div>
      <!-- End of products -->


      <div class="flex flex-col mt-10 space-y-4">
        <div class="flex">
          <h1 class="text-xl font-extrabold text-gray-900">Billing</h1>
        </div>

        <sl-input :value="stripeCustomerId" @input="stripeCustomerId = $event.target.value" type="text"
          readonly disabled placeholder="None" label="Stripe CustomerID" />
      </div>
      <!-- End of billing -->

      <div class="flex flex-col mt-10">
        <div class="flex">
          <h1 class="text-xl font-extrabold text-gray-900">Orders</h1>
        </div>

        <div class="flex">
          <OrdersList :orders="orders" />
        </div>
      </div>
      <!-- End of orders -->

    </div>

  </div>

  <DeleteDialog v-model="showDeleteContactDialog" :error="deleteContactDialogError"
    :title="deleteContactDialogTitle" :message="deleteContactDialogMessage" :loading="deleteContactDialogLoading"
    @delete="deleteContact"
  />
</template>

<script lang="ts" setup>
import type { BlockContactInput, Contact, CreateContactInput, Order, Product, UnblockContactInput, UpdateContactInput } from '@/api/model';
import { ref, type PropType, watch, onBeforeMount, type Ref, computed } from 'vue';
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue';
import { EllipsisVerticalIcon } from '@heroicons/vue/24/outline';
import DeleteDialog from '@/ui/components/mdninja/delete_dialog.vue';
import { useRoute, useRouter } from 'vue-router';
import { useMdninja } from '@/api/mdninja';
import OrdersList from '@/ui/components/products/orders_list.vue';
import date from 'mdninja-js/src/libs/date';
import SlSwitch from '@shoelace-style/shoelace/dist/components/switch/switch.js';
import { oneRouteUp } from '@/libs/router_utils';
import SlInput from '@shoelace-style/shoelace/dist/components/input/input.js';
import SlButton from '@shoelace-style/shoelace/dist/components/button/button.js';
import { countryName } from '@/libs/countries';

// props
const props = defineProps({
  websiteId: {
    type: String as PropType<string>,
    required: true,
  },
  contact: {
    type: Object as PropType<Contact | null>,
    required: false,
    default: null,
  },
});

// events
const $emit = defineEmits(['updated', 'deleted']);

// composables
const $router = useRouter();
const $mdninja = useMdninja();
const $route = useRoute();

// lifecycle
onBeforeMount(() => resetValues(props.contact));

// variables
const backRoute = oneRouteUp($route.path);
const deleteContactDialogTitle = 'Delete Contact';
const deleteContactDialogMessage = 'Are you sure you want to delete this Contact? This action cannot be undone.';


let error = ref('');
let loading = ref(false);

let showDeleteContactDialog = ref(false);
let deleteContactDialogError = ref('');
let deleteContactDialogLoading = ref(false);

let email = ref('');
let name = ref('');

let stripeCustomerId = ref('');
let subscribedToNewsletter = ref(false);

let products: Ref<Product[]> = ref([]);
let orders: Ref<Order[]> = ref([]);

// computed
const blocked = computed(() => props.contact?.blocked_at ? true : false);

// watch
watch(() => props.contact, (to) => resetValues(to), { deep: true });

// functions
function productUrl(product: Product): string {
  return `/websites/${props.websiteId}/products/${product.id}`;
}

function resetValues(contact: Contact | null) {
  if (contact) {
    email.value = contact.email;
    name.value = contact.name;
    subscribedToNewsletter.value = contact.subscribed_to_newsletter_at ? true : false;
    stripeCustomerId.value = contact.stripe_customer_id ?? '';
    products.value = contact.products ?? products.value;
    orders.value = contact.orders ?? orders.value;
  } else {
    email.value = '';
    name.value = '';
    subscribedToNewsletter.value = false;
    stripeCustomerId.value = '';
    products.value = [];
    orders.value = [];
  }
}

async function createContact() {
  loading.value = true;
  error.value = '';
  const input: CreateContactInput = {
    website_id: props.websiteId,
    email: email.value,
    name: name.value,
  };

  try {
    const newContact = await $mdninja.createContact(input);
    $router.push(`/websites/${props.websiteId}/contacts/${newContact.id}`);
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function updateContact() {
  loading.value = true;
  error.value = '';

  const input: UpdateContactInput = {
    id: props.contact!.id!,
    email: email.value,
    name: name.value,
    subscribed_to_newsletter: subscribedToNewsletter.value,
  };

  try {
    const updatedContact = await $mdninja.updateContact(input);
    $emit('updated', updatedContact);
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function blockOrUnblockContact() {
  loading.value = true;
  error.value = '';

  const input: BlockContactInput | UnblockContactInput = {
    id: props.contact!.id!,
  };

  try {
    const updatedContact = props.contact?.blocked_at ?
      await $mdninja.unblockContact(input) :
      await $mdninja.blockContact(input);
    $emit('updated', updatedContact);
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function deleteContact() {
  deleteContactDialogLoading.value = true;
  deleteContactDialogError.value = '';

  try {
    await $mdninja.deleteContact(props.contact!.id);
    showDeleteContactDialog.value = false;
    $router.push(backRoute)
  } catch (err: any) {
    deleteContactDialogError.value = err.message;
  } finally {
    deleteContactDialogLoading.value = false;
  }
}

function openDeleteContactDialog() {
  showDeleteContactDialog.value = true;
}
</script>
