<template>
  <div class="bg-white">
    <div class="relative isolate overflow-hidden bg-white">
    <svg class="absolute inset-0 -z-10 size-full stroke-gray-200 [mask-image:radial-gradient(100%_100%_at_top_right,white,transparent)]" aria-hidden="true">
      <defs>
        <pattern id="0787a7c5-978c-4f66-83c7-11c213f99cb7" width="200" height="200" x="50%" y="-1" patternUnits="userSpaceOnUse">
          <path d="M.5 200V.5H200" fill="none" />
        </pattern>
      </defs>
      <rect width="100%" height="100%" stroke-width="0" fill="url(#0787a7c5-978c-4f66-83c7-11c213f99cb7)" />
    </svg>
    <div class="mx-auto max-w-7xl px-6 pb-14 lg:flex lg:px-8 mt-10">
      <div class="mx-auto max-w-2xl lg:mx-0 lg:shrink-0">
      <!--<img class="h-11" src="/m.svg" alt="Your Company" />
        <div class="mt-24 sm:mt-32 lg:mt-16">
          <a href="#" class="inline-flex space-x-6">
            <span class="rounded-full bg-(--primary-color)/10 px-3 py-1 text-sm/6 font-semibold text-(--primary-color) ring-1 ring-inset ring-(--primary-color)/10">What's new</span>
            <span class="inline-flex items-center space-x-2 text-sm/6 font-medium text-gray-600">
              <span>Just shipped v1.0</span>
              <ChevronRightIcon class="size-5 text-gray-400" aria-hidden="true" />
            </span>
          </a>
        </div> -->
        <h1 class="mt-10 lg:mt-15 text-pretty text-5xl sm:text-7xl font-semibold tracking-tight">
          Publishing <Typewriter />
          has never been easier
        </h1>
        <p class="mt-6 text-pretty text-lg text-neutral-500 sm:text-xl/8">
          Markdown-first CMS for bloggers, minimalists and startups.
        </p>
        <div class="mt-6 flex items-center gap-x-2 sm:gap-x-6">
          <a href="https://kerkour.com" target="_blank" rel="noopener">
            <sl-button outline class="items-center">
              View example blog
              <img class="inline ml-2 w-2.5 pb-0.5"
                src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAQElEQVR42qXKwQkAIAxDUUdxtO6/RBQkQZvSi8I/pL4BoGw/XPkh4XigPmsUgh0626AjRsgxHTkUThsG2T/sIlzdTsp52kSS1wAAAABJRU5ErkJggg==" />
            </sl-button>
          </a>

          <sl-button variant="primary" @click="$pingoo.signup()">
            Get started for free &nbsp; <span aria-hidden="true">→</span>
          </sl-button>
        </div>
      </div>
      <div class="mx-auto mt-16 flex max-w-2xl sm:mt-24 lg:ml-10 lg:mr-0 lg:mt-0 lg:max-w-none lg:flex-none">
        <div class="max-w-3xl flex-none sm:max-w-5xl lg:max-w-none">
          <img src="/webapp/about/app_screenshot.png" alt="App screenshot" class="w-[65rem] m-2 rounded-md shadow-2xl ring-1 ring-gray-900/10" />
        </div>
      </div>
    </div>
  </div>
    <!-- <div class="py-10 sm:py-20">
      <LogosCloud />
    </div> -->


    <Features class="py-12" />


    <div class="py-14">
      <Testimonials />
    </div>

    <div class="flex items-center justify-center py-2">
      <sl-button variant="primary" @click="$pingoo.signup()">
        Join <s>the cult</s> Markdown Ninja
      </sl-button>
    </div>



    <div class="py-12">
      <BeforeAfter />
    </div>


    <div class="py-6 px-6 sm:px-6 lg:px-8 mx-auto max-w-3xl">
      <h2 class="text-5xl font-semibold tracking-tight text-center">
        Ready to simplify your workflow?
      </h2>
      <div class="bg-white py-10 px-10">
        <Terminal prompt="$" command="markdown-ninja publish" :showCopyToClipboard="false" />
      </div>
      <div class="flex items-center justify-center gap-x-6">
        <sl-button variant="primary" @click="$pingoo.signup()">
          Try Markdown Ninja &nbsp; <span aria-hidden="true">→</span>
        </sl-button>
      </div>
    </div>

  </div>
</template>

<script lang="ts" setup>
import { usePingoo } from '@/pingoo/pingoo';
import BeforeAfter from '@/ui/components/about/before_after.vue';
import Features from '@/ui/components/about/features.vue';
import Testimonials from '@/ui/components/about/testimonials.vue';
import SlButton from '@shoelace-style/shoelace/dist/components/button/button.js';
import Typewriter from '../mdninja/typewriter.vue';
import Terminal from '../kernel/terminal.vue';

// props

// events

// composables
const $pingoo = usePingoo();

// lifecycle

// variables

// computed

// watch

// functions
</script>
