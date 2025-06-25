<template>
  <div class="flex-1">



    <div class="rounded-md bg-red-50 p-4 my-4" v-if="error">
      <div class="flex">
        <div class="ml-3">
          <p class="text-sm text-red-700">
            {{ error }}
          </p>
        </div>
      </div>
    </div>

    <div v-if="organization" class="flex flex-col space-y-5 mt-5">
      <div class="flex">
        <h3 class="font-bold">ID:</h3> &nbsp; <span>{{ organization.id }}</span>
      </div>

      <div class="flex">
        <h3 class="font-bold">Created At:</h3> &nbsp; <span>{{ organization.created_at }}</span>
      </div>

      <div class="flex">
        <h3 class="font-bold">Name:</h3> &nbsp; <span>{{ organization.name }}</span>
      </div>


      <div class="flex">
        <h3 class="font-bold">Plan:</h3> &nbsp; <span>{{ organization.plan }} ({{ organization.extra_slots }} extra slots)</span>
      </div>

      <div class="flex">
        <sl-button variant="primary" @click="syncStripeData()" :loading="loading">
          Sync Stripe data
        </sl-button>
      </div>

      <div class="flex">
        <sl-button variant="primary" @click="addStaff" :loading="loading">
          Add Staff
        </sl-button>
      </div>

      <StaffsList :staffs="organization.staffs!" @remove="removeStaff" />
    </div>

    <div class="flex flex-col mt-3">
      <h2 class="text-xl font-bold">Websites</h2>
      <WebsitesList :websites="websites" class="w-full mt-3" />
    </div>

  </div>
</template>

<script lang="ts" setup>
import { useMdninja } from '@/api/mdninja';
import type { Organization, RemoveStaffInput, Staff, Website } from '@/api/model';
import { onBeforeMount, ref, type Ref } from 'vue';
import { useRoute } from 'vue-router';
import StaffsList from '@/ui/components/organizations/staffs_list.vue';
import SlButton from '@shoelace-style/shoelace/dist/components/button/button.js';
import WebsitesList from '@/ui/components/admin/websites_list.vue';

// props

// events

// composables
const $route = useRoute();
const $mdninja = useMdninja();

// lifecycle
onBeforeMount(() => fetchData());


// variables
const organizationId = $route.params.organization_id as string;

let loading = ref(false);
let error = ref('');
let organization: Ref<Organization | null> = ref(null);
let websites: Ref<Website[]> = ref([]);


// computed

// watch

// functions
async function fetchData() {
  loading.value = true;
  error.value = '';

  try {
    const [organizationRes, websitesRes] = await Promise.all([
      $mdninja.getOrganization({ id: organizationId, staffs: true }),
      $mdninja.listWebsites({ organization_id: organizationId }),
    ]);
    organization.value = organizationRes;
    websites.value = websitesRes;
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function syncStripeData() {
  error.value = '';
  loading.value = true;

  try {
    await $mdninja.organizationSyncStripe(organizationId);
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function addStaff() {
  error.value = '';

  const userID = prompt("User ID:");
  if (!userID) {
    return
  }

  loading.value = true;

  try {
    const newStaffs = await $mdninja.addStaffs({ organization_id: organizationId, user_ids: [userID] });
    organization.value!.staffs?.push(...newStaffs);
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}

async function removeStaff(staff: Staff) {
  loading.value = true;
  error.value = '';
  const intput: RemoveStaffInput = {
    organization_id: organizationId,
    user_id: staff.user_id,
  }

  try {
    await $mdninja.removeStaff(intput);
    organization.value!.staffs = organization.value!.staffs!.filter((sta) => sta.user_id !== staff.user_id);
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}
</script>
