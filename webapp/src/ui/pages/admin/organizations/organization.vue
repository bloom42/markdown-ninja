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
        <h3 class="font-bold">Billing email:</h3> &nbsp; <span>{{ organization.billing_information.email }}</span>
      </div>

      <div class="flex">
        <h3 class="font-bold">Stripe Customer ID:</h3> &nbsp; <span>{{ organization.stripe_customer_id }}</span>
      </div>

      <div class="flex">
        <h3 class="font-bold">Stripe Subscription ID:</h3> &nbsp; <span>{{ organization.stripe_subscription_id }}</span>
      </div>


      <div class="flex">
        <h3 class="font-bold">Payment Due:</h3> &nbsp; <span>{{ organization.payment_due }}</span>
      </div>


      <div class="flex flex-col">
        <SelectPlan v-model="organization.plan" :all-plans="true" />
      </div>

      <div class="flex flex-col">
        <sl-input label="Extra Slots"
          :value="organization.extra_slots" @input="organization.extra_slots = parseInt($event.target.value, 10)" min="0" type="number"
          placeholder="0"
        />
      </div>

      <div class="flex">
        <sl-button variant="primary" @click="updateOrganization()" :loading="loading">
          Update Organization
        </sl-button>
      </div>

      <hr />

      <div class="flex">
        <sl-button variant="primary" @click="syncStripeData()" :loading="loading">
          Sync Stripe data
        </sl-button>
      </div>



      <div class="flex flex-col mt-5 space-y-2">
        <h2 class="text-2xl font-bold">Usage</h2>

        <BillingUsage :billing-usage="billingUsage!" />
      </div>


      <div class="flex flex-col mt-5 space-y-2">
        <h2 class="text-2xl font-bold">Staffs</h2>

        <div>
          <sl-button variant="primary" @click="addStaff" :loading="loading">
            Add Staff
          </sl-button>
        </div>

        <StaffsList :staffs="organization.staffs!" @remove="removeStaff" />
      </div>


    </div>

    <div class="flex flex-col mt-3 space-y-2">
      <h2 class="text-2xl font-bold">Websites</h2>

      <WebsitesList :websites="websites" />
    </div>

  </div>
</template>

<script lang="ts" setup>
import { useMdninja } from '@/api/mdninja';
import type { Organization, OrganizationBillingUsage, RemoveStaffInput, Staff, UpdateOrganizationInput, Website } from '@/api/model';
import { onBeforeMount, ref, type Ref } from 'vue';
import { useRoute } from 'vue-router';
import StaffsList from '@/ui/components/organizations/staffs_list.vue';
import SlButton from '@shoelace-style/shoelace/dist/components/button/button.js';
import WebsitesList from '@/ui/components/admin/websites_list.vue';
import BillingUsage from '@/ui/components/organizations/billing_usage.vue';
import SlInput from '@shoelace-style/shoelace/dist/components/input/input.js';
import SelectPlan from '@/ui/components/organizations/select_plan.vue';

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
let billingUsage: Ref<OrganizationBillingUsage | null> = ref(null);


// computed

// watch

// functions
async function fetchData() {
  loading.value = true;
  error.value = '';

  try {
    const [organizationRes, websitesRes, billingUsageRes] = await Promise.all([
      $mdninja.getOrganization({ id: organizationId, staffs: true }),
      $mdninja.listWebsites({ organization_id: organizationId }),
      $mdninja.getorganizationBillingUsage(organizationId),
    ]);
    organization.value = organizationRes;
    websites.value = websitesRes;
    billingUsage.value = billingUsageRes;
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
    await fetchData();
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

async function updateOrganization() {
  loading.value = true;
  error.value = '';
  const intput: UpdateOrganizationInput = {
    id: organizationId,
    plan: organization.value?.plan,
    extra_slots: organization.value?.extra_slots,
  };

  try {
    organization.value = await $mdninja.updateOrganization(intput);
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}
</script>
