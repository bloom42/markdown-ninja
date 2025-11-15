export const Routes = {
  //////////////////////////////////////////////////////////////////////////////////////////////////
  // kernel
  //////////////////////////////////////////////////////////////////////////////////////////////////
  // misc
  challengeSiteKey: '/challenge_site_key',
  listTlsCertificates: '/list_tls_certificates',
  deleteTlsCertificate: '/delete_tls_certificate',

  // background jobs
  failedBackgroundJobs: '/queue/failed_background_jobs',
  deleteBackgroundJob: '/queue/delete_background_job',

  // users
  init: '/init',
  deleteUser: '/delete_user',
  user: '/user',
  users: '/users',
  updateUser: '/update_user',

  // sessions
  revokeSession: '/revoke_session',
  sessions: '/sessions',

  // organizations
  organizations: '/organizations/list',
  createOrganization: '/organizations/create',
  udpateOrganization: '/organizations/update',
  organization: '/organizations/get',
  deleteOrganization: '/organizations/delete',
  allorganizations: '/organizations/all',
  organizationUpdateSubscription: '/organizations/update_subscription',
  organizationStripeCustomerPortal: '/organizations/stripe_customer_portal',
  organizationSyncStripe: '/organizations/sync_stripe',
  organizationBillingUsage: '/organizations/billing_usage',
  organizationsAdminStatistics: '/organizations/admin-statistics',

  // staffs
  staffs: '/staffs',
  inviteStaffs: '/invite_staffs',
  addStaffs: '/add_staffs',
  removeStaff: '/remove_staff',
  deleteStaffInvitation: '/delete_staff_invitation',
  acceptStaffInvitation: '/accept_staff_invitation',
  staffInvitations: '/staff_invitations',
  userInvitations: '/user_invitations',

  // apiKeys
  deleteApiKey: '/delete_api_key',
  createApiKey: '/create_api_key',


  //////////////////////////////////////////////////////////////////////////////////////////////////
  // websites
  //////////////////////////////////////////////////////////////////////////////////////////////////

  // websites
  createWebsite: '/create_website',
  deleteWebsite: '/delete_website',
  websites: '/websites',
  website: '/website',
  updateWebsite: '/update_website',
  saveRedirects: '/save_redirects',
  allWebsites: '/all_websites',
  websiteUpdateIcon: '/websites/update_icon',
  websitesAdminStatistics: '/websites/admin-statistics',

  // tags
  createTag: '/create_tag',
  updateTag: '/update_tag',
  deleteTag: '/delete_tag',
  tags: '/tags',

  // snippets
  createSnippet: '/create_snippet',
  updateSnippet: '/update_snippet',
  deleteSnippet: '/delete_snippet',
  snippets: '/snippets',

  // pages
  page: '/page',
  createPage: '/create_page',
  updatePage: '/update_page',
  deletePage: '/delete_page',
  pages: '/pages',
  posts: '/posts',

  // assets
  uploadAsset: '/upload_asset',
  deleteAsset: '/delete_asset',
  assets: '/assets',
  createAssetFolder: '/create_asset_folder',

  // domains
  addDomain: '/add_domain',
  removeDomain: '/remove_domain',
  setDomainAsPrimary: '/set_domain_as_primary',
  checkTlsCertificateForDomain: '/check_tls_certificate_for_domain',

  //////////////////////////////////////////////////////////////////////////////////////////////////
  // Contacts
  //////////////////////////////////////////////////////////////////////////////////////////////////

  // contacts
  createContact: '/create_contact',
  contact: '/contact',
  contacts: '/contacts',
  deleteContact: '/delete_contact',
  updateContact: '/update_contact',
  importContacts: '/import_contacts',
  exportContacts: '/export_contacts',
  exportContactsForProduct: '/export_contacts_for_product',
  blockContact: '/block_contact',
  unblockContact: '/unblock_contact',

  // labels
  createLabel: '/create_label',
  deleteLabel: '/delete_label',
  updateLabel: '/update_label',
  labels: '/labels',

  //////////////////////////////////////////////////////////////////////////////////////////////////
  // emails
  //////////////////////////////////////////////////////////////////////////////////////////////////

  // configuration
  emailsConfiguration: '/emails_configuration',
  updateEmailsConfiguration: '/update_emails_configuration',
  verifyEmailsConfiguration: '/verify_emails_dns_configuration',

  // newsletters
  newsletters: '/newsletters',
  newsletter: '/newsletter',
  createNewsletter: '/create_newsletter',
  updateNewsletter: '/update_newsletter',
  deleteNewsletter: '/delete_newsletter',
  sendNewsletter: '/send_newsletter',

  //////////////////////////////////////////////////////////////////////////////////////////////////
  // Products
  //////////////////////////////////////////////////////////////////////////////////////////////////

  // products
  products: '/products',
  product: '/product',
  createProduct: '/create_product',
  updateProduct: '/update_product',
  giveContactsAccessToProduct: '/give_contacts_access_to_product',
  deleteProduct: '/delete_product',
  removeAccessToproduct: '/remove_access_to_product',

  // orders
  orders: '/orders',
  order: '/order',

  // refunds
  refunds: '/refunds',
  createRefund: '/create_refund',

  // pages
  createProductPage: '/create_product_page',
  updateProductPage: '/update_product_page',
  deleteProductPage: '/delete_product_page',
  productPage: '/product_page',

  // coupons
  coupons: '/coupons',
  coupon: '/coupon',
  createCoupon: '/create_coupon',
  updateCoupon: '/update_coupon',

  //////////////////////////////////////////////////////////////////////////////////////////////////
  // Analytics
  //////////////////////////////////////////////////////////////////////////////////////////////////
  analyticsData: '/analytics_data',
}
