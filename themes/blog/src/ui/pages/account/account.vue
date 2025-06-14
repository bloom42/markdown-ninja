<template>
  <div class="flex flex-col">

    <div class="rounded-md bg-red-50 p-4 mb-3" v-if="error">
      <div class="flex">
        <div class="ml-3">
          <p class="text-sm text-red-700">
            {{ error }}
          </p>
        </div>
      </div>
    </div>

    <div class="rounded-md bg-green-50 p-4 mb-5 mt-4" v-if="emailUpdated">
      <div class="flex">
        <div class="ml-3">
          <p class="text-sm text-green-700">
            Almost finished... We need to confirm your email address. Please click the link in the email we just sent you.
          </p>
        </div>
      </div>
    </div>

    <div class="rounded-md bg-green-50 p-4 mb-5 mt-4" v-if="newEmailVerified">
      <div class="flex">
        <div class="ml-3">
          <p class="text-sm text-green-700">
            Success! Your new email address is now verified.
          </p>
        </div>
      </div>
    </div>

    <div class="flex flex-col" v-if="$store.contact">

      <div class="flex flex-row justify-between items-center">
        <div class="flex">
          <h1>Account</h1>
        </div>
        <div class="flex">
          <button @click="onLogoutClicked()" :loading="loading"
            class="hover:brightness-[0.95] text-[var(--mdninja-text)] h-fit mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 px-4 py-2 bg-white text-base font-medium sm:mt-0 sm:w-auto sm:text-sm" ref="cancelButtonRef">
            Logout
          </button>
        </div>
      </div>

      <div class="border rounded-md border-gray-200 relative p-4 flex">
        <div v-if="editingNameAndEmail" class="flex flex-grow flex-col space-y-5">
          <div class="flex flex-col space-y-4">

            <div class="flex-1">
              <label for="email" class="block text-sm font-medium">
                Email
              </label>
              <div class="mt-1">
                <input id="email" name="email" type="email" required placeholder="you@email.com"
                  v-model="email" :disabled="loading"
                  class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-xs placeholder-gray-400 focus:outline-hidden focus:ring-[var(--mdninja-accent)] focus:border-[var(--mdninja-accent)] sm:text-sm"
                />
              </div>
            </div>

            <div class="flex-1">
              <label for="name" class="block text-sm font-medium">
                Name
              </label>
              <div class="mt-1">
                <input id="name" name="name" type="text" required placeholder="Your Name"
                  v-model="name" :disabled="loading"
                  class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-xs placeholder-gray-400 focus:outline-hidden focus:ring-[var(--mdninja-accent)] focus:border-[var(--mdninja-accent)] sm:text-sm"
                />
              </div>
            </div>

          </div>
          <div class="flex flex-row space-x-4">
            <button @click="cancelEditNameAndEmailClicked()" :loading="loading"
              class="hover:brightness-[0.95] text-[var(--mdninja-text)] h-fit mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 px-4 py-2 bg-white text-base font-medium sm:mt-0 sm:w-auto sm:text-sm" ref="cancelButtonRef">
              Cancel
            </button>
            <PButton @click="updateNameAndEmail()" :loading="loading">
              Save
            </PButton>
          </div>
        </div>

        <div v-else class="flex flex-grow">
          <div class="ml-3 flex flex-col flex-grow">
            <div class="block text-base font-medium">
              {{ $store.contact!.name }}
            </div>
            <div class="block text-sm opacity-60">
              {{ $store.contact!.email }}
            </div>
          </div>

          <div class="flex">
            <button @click="editNameAndEmailClicked"
              class="bg-[var(--mdninja-background)] hover:brightness-[0.95] text-[var(--mdninja-accent)] rounded-md font-medium  shadow-none">
              Edit
            </button>
          </div>
        </div>

      </div>

      <div class="mt-5 border rounded-md border-gray-200 relative p-4 flex">
        <div class="flex flex-grow">
          <div class="ml-3 flex flex-col flex-grow">
            <div class="block text-base font-medium">
              Email Newsletter
            </div>
            <div class="block text-sm opacity-60">
              <span v-if="subscribedToNewsletter">
                Subscribed
              </span>
              <span v-else>
                Unsubscribed
              </span>
            </div>
          </div>

          <div class="flex items-center">
            <Switch v-model="subscribedToNewsletter" as="div" @click="subscribeOrUnsubscribe()"
              :class="[subscribedToNewsletter ? 'bg-[var(--mdninja-accent)]' : 'bg-gray-200', 'ml-10 relative inline-flex flex-shrink-0 h-6 w-11 border-2 border-transparent rounded-full cursor-pointer transition-colors ease-in-out duration-200']">
              <span aria-hidden="true"
                :class="[subscribedToNewsletter ? 'translate-x-5' : 'translate-x-0', 'pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow-xs transform ring-0 transition ease-in-out duration-200']" />
            </Switch>
          </div>
        </div>
      </div>


      <div class="flex">
        <h2>Products</h2>
      </div>

      <div class="flex">
        <ProductsList :products="products" />
      </div>


      <div class="flex">
        <h2>Orders</h2>
      </div>

      <div class="flex">
        <OrdersList :orders="orders" />
      </div>


      <div class="flex flex-col my-20">
        <h2 class="text-2xl font-medium text-red-500 my-0">Danger Zone</h2>
        <p class="text-sm text-red-500 my-2">Irreversible and destructive actions.</p>

        <div class="mt-5 flex">
          <PButton :loading="loading" @click="onDeleteMyAccountClicked()" class="bg-red-500">
            Delete My Account
          </PButton>
        </div>
      </div>
    </div>


  </div>
</template>

<script lang="ts" setup>
import { useStore } from '@/app/store';
import { deleteMyAccount, fetchMe, listMyOrders, listMyProducts, logout, trackPage, updateMyAccount, verifyEmail } from '@/app/mdninja';
import { onBeforeMount, onBeforeUpdate, ref, type Ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Switch } from '@headlessui/vue';
import type { Order, Product, UpdateMyAccountInput, VerifyEmailInput, VerifyEmailJwt } from '@/app/model';
import PButton from '@/ui/components/p_button.vue';
import OrdersList from '@/ui/components/orders_list.vue';
import ProductsList from '@/ui/components/products_list.vue';
import { jwtDecode } from '@/libs/jwt';

const SUCCESS_MESSAGE_TIMEOUT = 10_000;

// props

// events

// composables
const $store = useStore();
const $router = useRouter();
const $route = useRoute();

// lifecycle
onBeforeMount(() => {
  document.title = `${website.name} - Account`;

  trackPage();
  fetchData();

  if (updateEmailToken) {
    callVerifyEmail();
  }
});

// this is an ugly hack for when we reload the page without the query parameters after callVerifyEmail
// without that, the page will how the app's loader
// TODO: improve
onBeforeUpdate(() => {
  $store.setLoading(false);
});

// variables
const website = $store.website!;

let error = ref('');
let loading = ref(false);

let subscribedToNewsletter = ref($store.contact?.subscribed_to_newsletter ?? true);
let name = ref($store.contact?.name ?? '');
let editingNameAndEmail = ref(false);

let email = ref($store.contact?.email ?? '');
let emailUpdated = ref(false);
let updateEmailToken = $route.query['update-email-token'] as string | undefined;
let newEmailVerified = ref(false);

let orders: Ref<Order[]> = ref([]);
let products: Ref<Product[]> = ref([]);

// computed

// watch

// functions
function resetValues() {
  if ($store.contact) {
    name.value = $store.contact.name;
    email.value = $store.contact.email;
    subscribedToNewsletter.value = $store.contact.subscribed_to_newsletter;
  } else {
    name.value = '';
    email.value = '';
    subscribedToNewsletter.value = false;
  }
}

async function fetchData() {
  loading.value = true;
  error.value = '';

  try {
    const [me, ordersApiRes, productsApiRes] = await Promise.all([
      fetchMe(),
      listMyOrders(),
      listMyProducts(),
    ]);
    $store.setContact(me!);
    resetValues();
    orders.value = ordersApiRes.data;
    products.value = productsApiRes.data;
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
    $store.setLoading(false);
  }
}

function editNameAndEmailClicked() {
  editingNameAndEmail.value = true;
}

function cancelEditNameAndEmailClicked() {
  editingNameAndEmail.value = false;
  name.value = $store.contact!.name;
  email.value = $store.contact!.email;
}

async function onLogoutClicked() {
  loading.value = true;
  error.value = '';

  try {
    await logout();
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function subscribeOrUnsubscribe() {
  loading.value = true;
  error.value = '';
  const input: UpdateMyAccountInput = {
    subscribed_to_newsletter: !subscribedToNewsletter.value,
  };

  try {
    const contact = await updateMyAccount(input);
    $store.setContact(contact);
    subscribedToNewsletter.value = contact.subscribed_to_newsletter;
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function updateNameAndEmail() {
  loading.value = true;
  error.value = '';
  const newEmail = email.value !== $store.contact?.email ? email.value : undefined;

  const input: UpdateMyAccountInput = {
    name: name.value,
    email: newEmail,
  };

  try {
    const contact = await updateMyAccount(input);
    $store.setContact(contact);
    editingNameAndEmail.value = false;
    if (newEmail) {
      emailUpdated.value = true;
      setTimeout(() => {
        emailUpdated.value = false;
      }, SUCCESS_MESSAGE_TIMEOUT);
    }
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function callVerifyEmail() {
  loading.value = true;
  error.value = '';
  const input: VerifyEmailInput = {
    token: updateEmailToken!,
  };

  try {
    const decodedJwt: VerifyEmailJwt = jwtDecode(updateEmailToken!);
    await verifyEmail(input);

    $store.updateContact({ email: decodedJwt.email! });

    $router.push({ query: {} });
    newEmailVerified.value = true;
    setTimeout(() => {
      newEmailVerified.value = false;
    }, SUCCESS_MESSAGE_TIMEOUT);
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function onDeleteMyAccountClicked() {
  if (!confirm("Do you really want to delete your account and all the associated data? This action is irreversible and you will lose access to your orders and invoices history.")) {
    return;
  }

  loading.value = true;
  error.value = '';

  try {
    await deleteMyAccount();
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}
</script>
