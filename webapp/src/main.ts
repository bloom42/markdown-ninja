
import '@shoelace-style/shoelace/dist/themes/light.css';
import './index.css';

import App from '@/ui/app.vue';
import { createApp } from 'vue';
import { createPinia } from 'pinia';
import { newRouter } from '@/app/router';
import { createConfig } from '@/app/config';
import { createMdninja, getInitData } from '@/api/mdninja';
import { createPingooClient } from './pingoo/pingoo';
import { createMarkdownNinjaClient } from './sdk/markdown_ninja_sdk';
import { createLinkify } from './libs/linkify';
import { useStore } from './app/store';

async function main() {
  const config = createConfig();
  if (config.env !== 'production') {
    console.log(config);
  }

  createMarkdownNinjaClient(config.cmsBaseUrl);

  const router = newRouter(config);
  createLinkify(router);

  const app = createApp(App);

  const pinia = createPinia()
  app.use(pinia);

  await getInitData()
  const $store = useStore();

  const pingooClient = createPingooClient({
    endpoint: $store.pingooEndpoint,
    appId: $store.pingooAppId,
    redirectUri: config.oidcRedirectUri,
  });
  await pingooClient.init();

  const mdninjaService = createMdninja(config, router);
  await mdninjaService.init();

  // the router must be initialized after init for the beforeEach hook to works correctly
  app.use(router);
  await router.isReady();

  app.mount('#markdown-ninja-webapp')
}

main();
