import { createRouter, createWebHistory, type Router } from 'vue-router';
import type { Config } from './config';
import { useStore } from './store';
import { usePingoo } from '@/pingoo/pingoo';


// IMPORTANT: When you add a new route, add it to the bundles in 'vite.config.js'

// account
// const Auth = () => import('@/ui/pages/auth.vue');
import Index from '@/ui/pages/index.vue';
import Page404 from '@/ui/pages/404.vue';
import Auth from '@/ui/pages/auth.vue';
import Content from '@/ui/pages/content.vue';
import Pricing from '@/ui/pages/pricing.vue';

// Organizations
import Organizations from '@/ui/pages/organizations/organizations.vue';
import NewOrganization from '@/ui/pages/organizations/new.vue';
import OrganizationStaffs from '@/ui/pages/organizations/organization/staffs.vue';
import OrganizationBilling from '@/ui/pages/organizations/organization/billing/billing.vue';
import OrganizationBillingCheckoutComplete from '@/ui/pages/organizations/organization/billing/checkout_complete.vue';
import OrganizationSettings from '@/ui/pages/organizations/organization/settings.vue';
import Websites from '@/ui/pages/organizations/organization/websites.vue';
import NewWebsite from '@/ui/pages//organizations/organization/new_website.vue';
import OrganizationApi from '@/ui/pages/organizations/organization/api.vue';
import OrganizationInvitations from '@/ui/pages/organizations/invitations.vue';

// Websites
import Website from '@/ui/pages/websites/website/website.vue';
import WebsitePosts from '@/ui/pages/websites/website/posts/posts.vue';
import WebsitePost from '@/ui/pages/websites/website/posts/post.vue';
import WebsiteNewPost from '@/ui/pages/websites/website/posts/new.vue';
import WebsitePages from '@/ui/pages/websites/website/pages/pages.vue';
import WebsitePage from '@/ui/pages/websites/website/pages/page.vue';
import WebsiteNewPage from '@/ui/pages/websites/website/pages/new.vue';
import WebsiteSnippets from '@/ui/pages/websites/website/settings/snippets.vue';
import WebsiteTags from '@/ui/pages/websites/website/settings/tags.vue';
import WebsiteAssets from '@/ui/pages/websites/website/assets.vue';
import WebsiteRedirects from '@/ui/pages/websites/website/settings/redirects.vue';
import WebsiteNavigation from '@/ui/pages/websites/website/settings/navigation.vue';

// Contacts
import WebsiteContacts from '@/ui/pages/websites/website/contacts/contacts.vue';
import WebsiteNewContact from '@/ui/pages/websites/website/contacts/new.vue';
import WebsiteContact from '@/ui/pages/websites/website/contacts/contact.vue';

// Emails
import WebsiteNewsletters from '@/ui/pages/websites/website/newsletters/newsletters.vue';
import WebsiteNewsletter from '@/ui/pages/websites/website/newsletters/newsletter.vue';
import WebsiteNewNewsletter from '@/ui/pages/websites/website/newsletters/new.vue';

// Store
import WebsiteProducts from '@/ui/pages/websites/website/products/products.vue';
import WebsiteProduct from '@/ui/pages/websites/website/products/product.vue';
import WebsiteCoupons from '@/ui/pages/websites/website/coupons/coupons.vue';
import WebsiteCoupon from '@/ui/pages/websites/website/coupons/coupon.vue';
import WebsiteNewCoupon from '@/ui/pages/websites/website/coupons/new.vue';
import WebsiteProductpage from '@/ui/pages/websites/website/products/product/pages/page.vue';
import WebsiteOrders from '@/ui/pages/websites/website/orders/orders.vue';
import WebsiteOrder from '@/ui/pages/websites/website/orders/order.vue';
import WebsiteRefunds from '@/ui/pages/websites/website/refunds/refunds.vue';

// Website Settings
import WebsiteSettings from '@/ui/pages/websites/website/settings/settings.vue';
import WebsiteSettingsCode from '@/ui/pages/websites/website/settings/code.vue';
import WebsiteSettingsDomains from '@/ui/pages/websites/website/settings/domains.vue';
import WebsiteSettingsEmails from '@/ui/pages/websites/website/settings/emails.vue';
import WebsiteSettingsDesign from '@/ui/pages/websites/website/settings/design.vue';

// Admin
import Admin from '@/ui/pages/admin/admin.vue';
import AdminOrganizations from '@/ui/pages/admin/organizations/organizations.vue';
import AdminOrganization from '@/ui/pages/admin/organizations/organization.vue';
import AdminWebsites from '@/ui/pages/admin/websites/websites.vue';
import AdminWebsite from '@/ui/pages/admin/websites/website.vue';
import AdminQueue from '@/ui/pages/admin/queue/queue.vue';
import AdminTls from '@/ui/pages/admin/tls.vue';

export function newRouter(config: Config): Router {
  const router = createRouter({
    history: createWebHistory(),
    // strict: true,
    routes: [
      { path: '/', component: Index, meta: { auth: false } },

      // about
      { path: '/pricing', component: Pricing, meta: { auth: false } },
      { path: '/docs', redirect: '/docs/cli' },
      { path: '/docs/:slug(.*)*', component: Content, meta: { auth: false } },
      { path: '/privacy', component: Content, meta: { auth: false } },
      { path: '/terms', component: Content, meta: { auth: false } },
      { path: '/about', component: Content, meta: { auth: false } },
      { path: '/contact', component: Content, meta: { auth: false } },
      // { path: '/privacy', component: PrivacyPolicy, meta: { auth: false } },
      // { path: '/terms', component: TermsOfService, meta: { auth: false } },
      // { path: '/about', component: About, meta: { auth: false } },
      // { path: '/contact', component: Contact, meta: { auth: false } },
      // { path: '/abuse', component: Abuse, meta: { auth: false } },

      // auth
      // the auth callback
      { path: '/auth', component: Auth, meta: { auth: false } },
      { path: '/login', redirect: '/' },
      { path: '/signup', redirect: '/' },

      // organizations
      { path: '/organizations', component: Organizations },
      { path: '/organizations/new', component: NewOrganization },
      { path: '/organizations/invitations', component: OrganizationInvitations },
      { path: '/organizations/:organization_id', redirect: (route) => `/organizations/${route.params.organization_id}/websites` },
      { path: '/organizations/:organization_id/websites', component: Websites },
      { path: '/organizations/:organization_id/websites/new', component: NewWebsite },
      { path: '/organizations/:organization_id/staffs', component: OrganizationStaffs },
      { path: '/organizations/:organization_id/billing', component: OrganizationBilling },
      { path: '/organizations/:organization_id/billing/checkout/complete', component: OrganizationBillingCheckoutComplete },
      { path: '/organizations/:organization_id/settings', component: OrganizationSettings },
      { path: '/organizations/:organization_id/api', component: OrganizationApi },

      // websites
      { path: '/websites/:website_id', component: Website, name: 'website_home' },
      { path: '/websites/:website_id/posts', component: WebsitePosts },
      { path: '/websites/:website_id/posts/:page_id', component: WebsitePost },
      { path: '/websites/:website_id/posts/new', component: WebsiteNewPost },
      { path: '/websites/:website_id/pages', component: WebsitePages },
      { path: '/websites/:website_id/pages/new', component: WebsiteNewPage },
      { path: '/websites/:website_id/pages/:page_id', component: WebsitePage },
      { path: '/websites/:website_id/snippets', component: WebsiteSnippets },
      { path: '/websites/:website_id/tags', component: WebsiteTags },
      { path: '/websites/:website_id/assets', component: WebsiteAssets },
      { path: '/websites/:website_id/redirects', component: WebsiteRedirects },
      { path: '/websites/:website_id/navigation', component: WebsiteNavigation },

      // Contacts
      { path: '/websites/:website_id/contacts', component: WebsiteContacts },
      { path: '/websites/:website_id/contacts/new', component: WebsiteNewContact },
      { path: '/websites/:website_id/contacts/:contact_id', component: WebsiteContact },

      // Email
      { path: '/websites/:website_id/newsletters', component: WebsiteNewsletters },
      { path: '/websites/:website_id/newsletters/new', component: WebsiteNewNewsletter },
      { path: '/websites/:website_id/newsletters/:newsletter_id', component: WebsiteNewsletter },

      // Store
      { path: '/websites/:website_id/coupons', component: WebsiteCoupons },
      { path: '/websites/:website_id/coupons/new', component: WebsiteNewCoupon },
      { path: '/websites/:website_id/coupons/:coupon_id', component: WebsiteCoupon },
      { path: '/websites/:website_id/products', component: WebsiteProducts },
      { path: '/websites/:website_id/products/:product_id', component: WebsiteProduct },
      { path: '/websites/:website_id/products/:product_id/pages/:page_id', component: WebsiteProductpage },
      { path: '/websites/:website_id/orders', component: WebsiteOrders },
      { path: '/websites/:website_id/orders/:order_id', component: WebsiteOrder },
      { path: '/websites/:website_id/refunds', component: WebsiteRefunds },

      // Website Settings
      { path: '/websites/:website_id/settings', component: WebsiteSettings },
      { path: '/websites/:website_id/settings/code', component: WebsiteSettingsCode },
      { path: '/websites/:website_id/settings/domains', component: WebsiteSettingsDomains },
      { path: '/websites/:website_id/settings/emails', component: WebsiteSettingsEmails },
      { path: '/websites/:website_id/settings/design', component: WebsiteSettingsDesign },

      // Admin
      { path: '/admin', component: Admin },
      { path: '/admin/organizations', component: AdminOrganizations },
      { path: '/admin/organizations/:organization_id', component: AdminOrganization },
      { path: '/admin/websites', component: AdminWebsites },
      { path: '/admin/websites/:website_id', component: AdminWebsite },
      { path: '/admin/websites/:website_id/contacts', component: WebsiteContacts },
      { path: '/admin/websites/:website_id/contacts/:contact_id', component: WebsiteContact },
      { path: '/admin/queue', component: AdminQueue },
      { path: '/admin/tls', component: AdminTls },


      // 404
      { path: '/:path(.*)*', component: Page404, meta: { auth: false } },
    ],
    scrollBehavior() {
      // always scroll to top
      return { top: 0 }
    },
  });

  router.beforeEach(async (to) => {
    const $store = useStore();
    const pingoo = usePingoo();

    // while (true) {
    //   if ($store.appLoaded) {
    //     break;
    //   }
    //   await sleep(20);
    // }

    if ((to.path === '/signup' || to.path.startsWith('/login')) && pingoo.isAuthenticated()) {
      return '/';
    }

    if (to.meta.auth !== false) {
      if (!pingoo.isAuthenticated()) {
        return '/login';
      }
    }

    if (to.path.startsWith('/admin') && $store.isAdmin !== true) {
      return '/';
    }
  });

  return router;
}
