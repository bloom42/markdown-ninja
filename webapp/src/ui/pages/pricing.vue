<template>
  <div class="bg-white py-10">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">

      <div class="mx-auto max-w-4xl text-center">
        <!-- <p class="mt-2 text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl">Simple and Transparent Pricing</p> -->
        <h2 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl">
          Simple and transparent pricing
        </h2>
      </div>
      <div class="mx-auto mt-6 max-w-2xl text-center text-lg leading-8 text-gray-900">
        <p>
          <!-- Building a Website and living off your passion has never been so simple. <br /> -->
          <!-- Save big with the all-in-one platform. <br /> -->
           <!-- Made for Humans, not bean counters. <br /> -->
          With Markdown Ninja you can focus on what matters, <br />
          <!-- Your audience.<br /> -->
          we take care of all the technical details.
        </p>
      </div>

      <div class="mx-auto max-w-6xl">
        <PlansTable  />
      </div>

      <!--
      Any questions? Talk to us at hello@markdown.ninja
      -->


      <!-- <div class="mx-auto max-w-6xl mt-10">
        <div class="relative isolate overflow-hidden px-6 py-10 text-center sm:rounded-3xl sm:px-16 border">
          <h2 class="mx-auto max-w-2xl text-4xl font-bold tracking-tight text-gray-900">Enterprise</h2>
          <p class="mx-auto mt-6 max-w-xl text-lg leading-8 text-gray-700">
            Do you have custom needs? Open Source hosting assistance, Huge volume, Priority support,
            Dedicated onboarding, Personalized SLAs, Custom limits,
            Custom features, Unique pricing, Training...
          </p>
          <div class="mt-8 flex items-center justify-center gap-x-6">
            <RouterLink to="/contact"
              class="block rounded-md bg-(--primary-color) py-2 px-3 text-center text-sm font-semibold leading-6 text-white shadow-xs hover:bg-(--primary-color) focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-(--primary-color)">
              Contact Us
            </RouterLink>
          </div>
        </div>
      </div> -->

      <!-- <PricingFaq class="mt-20 mb-10" /> -->

    </div>
  </div>
</template>

<script lang="ts" setup>
import PlansTable from '@/ui/components/about/plans_table.vue';
// import PricingFaq from '@/ui/components/about/pricing_faq.vue';

// props

// events

// composables

// lifecycle
document.title = 'Markdown Ninja - Pricing'

// variables

// computed

// watch

// functions
</script>
