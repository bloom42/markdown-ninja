import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue({
      template: {
        compilerOptions: {
          // treat all tags with a dash as custom elements
          isCustomElement: (tag) => tag.includes('-'),
        }
      },
    }),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:4080',
        // changeOrigin: true,
        secure: false,
        ws: true,
        configure: (proxy, _options) => {
          proxy.on('proxyReq', (proxyReq, req, _res) => {
            // forward the authorization header for API calls
            if (req.headers['authorization']) {
              proxyReq.setHeader('authorization', req.headers['authorization']);
            }
          });
        },
      },
      '/assets': {
        target: 'http://localhost:4080',
        // changeOrigin: true,
        secure: false,
        ws: true,
        configure: (proxy, _options) => {
          proxy.on('proxyReq', (proxyReq, req, _res) => {
            // forward the authorization header for API calls
            if (req.headers['authorization']) {
              proxyReq.setHeader('authorization', req.headers['authorization']);
            }
          });
        },
      }
    },
  },
  esbuild: {
    legalComments: 'none',
  },
  build:{
    // for debugging
    // minify: 'terser',
    // terserOptions: {
    //   compress: false,
    // //   mangle: false,
    //   parse: {
    //     html5_comments: false,
    //   },
    // },
    // modulePreload: false,
    manifest: true,
    cssCodeSplit: false,
    rollupOptions: {
      output: {
        // we use the full hashes to reduces the risk of collision with assets cached for long time
        assetFileNames(chunkInfo) {
          return `webapp/[name]-[hash:20][extname]`;
        },
        chunkFileNames(chunkInfo) {
          return `webapp/[name]-[hash:20].js`;
        },
        entryFileNames(chunkInfo) {
          return `webapp/[name]-[hash:20].js`;
        },
        // bundle as a single chunk
        // inlineDynamicImports: true,
        // manualChunks: {
        //   charts: ['chart.js'],
        // },
        // compact: true,
        // manualChunks: () => '',
        // manualChunks: (id, { getModuleIds, getModuleInfo }) => {
        //   if (id.includes('src/ui/components/')) {
        //     return '';
        //   }
        //   if (id.includes('chart') || id.includes('@kurkle') || id.includes('src/ui/pages/websites/')) {
        //     console.log(id);
        //     const info = getModuleInfo(id);
        //     console.log(info!.dynamicImporters)
        //     console.log(info!.dynamicallyImportedIds)
        //     console.log(info!.dynamicallyImportedIdResolutions)
        //     console.log(info!.importers)
        //     console.log(info!.importedIds)
        //     console.log(info?.importedIdResolutions);
        //     // console.log(getModuleInfo(id));
        //     return 'website';
        //   }
        //   // } else if (id.includes('codemirror') || id.includes('lezer') || id.includes('style-mod')
        //   //     || id.includes('crelt') || id.includes('@marijn/find-cluster-break') || id.includes('w3c-keyname')) {
        //   //   return 'codemirror';
        //   // }
        //   return '';
        // },

        // manualChunks: {
        //   // 'app': 'markdown_ninja_webapp',
        //   // codemirror: ['codemirror', '@codemirror/commands', '@codemirror/lang-markdown'],
        //   // charts: ['chart.js'],
        //   codemirror: [
        //     // Split CodeMirror code.
        //     'codemirror',
        //     '@codemirror/autocomplete',
        //     '@codemirror/commands',
        //     '@codemirror/language',
        //     '@codemirror/lint',
        //     '@codemirror/search',
        //     '@codemirror/state',
        //     '@codemirror/view',
        //   // ],
        //   // 'codemirror-lang': [
        //     // Add the following as needed.
        //     '@codemirror/lang-html',
        //     '@codemirror/lang-javascript',
        //     '@codemirror/lang-markdown',
        //     '@codemirror/lang-css',
        //   ],
        // }
        // manualChunks: {
        //   'admin': [
        //     './src/ui/pages/admin/admin.vue',
        //     './src/ui/pages/admin/organizations/organization.vue',
        //     './src/ui/pages/admin/organizations/organizations.vue',
        //     './src/ui/pages/admin/queue/queue.vue',
        //     './src/ui/pages/admin/websites/websites.vue',
        //     './src/ui/pages/admin/websites/website.vue',
        //   ],
        // }
          // manualChunks: (id: string, {getModuleInfo, getModuleIds}) => {
          // id = id.toLowerCase();
          // if (id.includes('node_modules') && !) {
          //   console.log(id)
          //   return 'vendor';
          // }
          // if (id.includes('/chart.js/')) {
          //   return 'chart.js';
          // }else if (id.includes('codemirror') || id.includes('lezer') || id.includes('style-mod')
          //     || id.includes('crelt') || id.includes('@marijn/find-cluster-break') || id.includes('w3c-keyname')) {
          //       console.log(id);
          //       console.log(getModuleInfo(id));
          //       return 'codemirror';
          //     }
          // if (id.includes('src/ui/pages/organizations')) {
          //   return 'organizations';
          // } else   if (id.includes('src/ui/pages/admin')) {
          //   return 'admin';
          // }
        //   return '';
        // },
        // manualChunks: {
        //   'organizations': [
        //     './src/ui/pages/organizations/*',
        //   ],
        //   'websites': [
        //     './src/ui/pages/websites/*',
        //   ],
        // },
        // manualChunks: (id: string) => {
        //   id = id.toLowerCase();
        //   if (id.includes('.css')) {
        //     console.log(id);
        //     return '';
        //   }
        //   if (id.includes('src/ui/pages/organizations')) {
        //     return 'organizations';
        //   } else   if (id.includes('src/ui/pages/admin')) {
        //     return 'admin';
        //   }

        //   return null;
        // },
        // manualChunks: (id: string) => {
        //   id = id.toLowerCase();
        //   if (id.includes('codemirror')) {
        //     console.log(id);
        //   }

        //   if (id.includes('src/ui/pages/admin')) {
        //     return 'admin';
        //   } else if (id.includes('src/ui/pages/organizations')) {
        //     return 'organizations';
        //   } else if (id.match(/src\/ui\/pages\/websites\/website\/(coupons|products|orders|refunds)/)) {
        //     return 'websites-store';
        //   } else if (id.includes('src/ui/pages/websites/website/settings')) {
        //     return 'websites-settings';
        //   } else if (id.includes('src/ui/pages/websites')) {
        //     return 'websites';
        //   // too small, it doesn't matter
        //   // } else if (id.match(/src\/ui\/pages\/(about|abuse|contact|pricing|privacy|terms)\.vue/)) {
        //   //   return 'about';
        //   } else if (id.includes('chart.js')) {
        //     return 'charts'
        //   // }
        //   // these are all the codemirror dependencies that need to be bundlded with it
        //   } else if (id.includes('codemirror') || id.includes('lezer') || id.includes('style-mod')
        //   || id.includes('crelt') || id.includes('@marijn/find-cluster-break') || id.includes('w3c-keyname')) {
        //     // console.log(id)
        //     return 'codemirror'
        //   }

        //   // all other modules go into the entry chunk (index-XXX)
        //   return '';
        // },
      },
    },
  },
//   optimizeDeps: {
//     exclude: [
//       "codemirror",
//       "@codemirror/state",
//       "@codemirror/lang-markdown",
//       "@codemirror/view",
//       "@codemirror/commands",
//       "@codemirror/language",
//       "@codemirror/search",
//       "@codemirror/autocomplete",
//       "@codemirror/lint",
//       "@codemirror/lang-html",
//       "@codemirror/lang-css",
//       "@codemirror/lang-javascript",
//       "@lezer/highlight",
//     ],
//  },
})

function includesAny(values: string[], str: string): boolean {
  for (const value of values) {
    if (value.includes(str)) {
      return true;
    }
  }

  return false;
}
