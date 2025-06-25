package api

// we use constants for routes so we can reuse them in clients and other places than router
// without the risk of introducing typos
const (
	// misc
	RoutePurgeCdn = "/purge_cdn"

	// webhooks
	RouteWebhooksStripe = "/webhooks/stripe"
	RouteWebhooksPingoo = "/webhooks/pingoo/{secret}"

	// background jobs
	RouteFailedBackgroundJobs = "/queue/failed_background_jobs"
	RouteDeleteBackgroundJob  = "/queue/delete_background_job"

	// users
	RouteInit = "/init"

	// organizations
	RouteOrganizationsAdminStatistics     = "/organizations/admin-statistics"
	RouteOrganizations                    = "/organizations/list"
	RouteOrganization                     = "/organizations/get"
	RouteCreateOrganization               = "/organizations/create"
	RouteDeleteOrganization               = "/organizations/delete"
	RouteUpdateOrganization               = "/organizations/update"
	RouteAllOrganizations                 = "/organizations/all"
	RouteOrganizationUpdateSubscription   = "/organizations/update_subscription"
	RouteOrganizationStripeCustomerPortal = "/organizations/stripe_customer_portal"
	RouteOrganizationSyncStripe           = "/organizations/sync_stripe"
	RouteOrganizationBillingUsage         = "/organizations/billing_usage"

	// staffs
	RouteStaffs                = "/staffs"
	RouteInviteStaffs          = "/invite_staffs"
	RouteAddStaffs             = "/add_staffs"
	RouteDeleteStaffInvitation = "/delete_staff_invitation"
	RouteRemoveStaff           = "/remove_staff"
	RouteStaffInvitations      = "/staff_invitations"
	RouteUserInvitations       = "/user_invitations"
	RouteAcceptStaffInvitation = "/accept_staff_invitation"

	// website
	RouteCreateWebsite     = "/create_website"
	RouteWebsite           = "/website"
	RouteUpdateWebsite     = "/update_website"
	RouteDeleteWebsite     = "/delete_website"
	RouteWebsites          = "/websites"
	RouteAllWebsites       = "/all_websites"
	RouteWebsiteUpdateIcon = "/websites/update_icon"

	// domains
	RouteAddDomain                    = "/add_domain"
	RouteRemoveDomain                 = "/remove_domain"
	RouteSetDomainAsPrimary           = "/set_domain_as_primary"
	RouteCheckTlsCertificateForDomain = "/check_tls_certificate_for_domain"

	// api Keys
	RouteCreateApiKey = "/create_api_key"
	RouteDeleteApiKey = "/delete_api_key"
	RouteUpdateApiKey = "/update_api_key"

	// pages
	RouteCreatePage = "/create_page"
	RouteUpdatePage = "/update_page"
	RouteDeletePage = "/delete_page"
	RoutePage       = "/page"
	RoutePages      = "/pages"
	RoutePosts      = "/posts"

	// redirects
	RouteSaveRedirect = "/save_redirects"

	// assets
	RouteUploadAsset       = "/upload_asset"
	RouteDeleteAsset       = "/delete_asset"
	RouteAssets            = "/assets"
	RouteCreateAssetFolder = "/create_asset_folder"

	// snippets
	RouteCreateSnippet = "/create_snippet"
	RouteUpdateSnippet = "/update_snippet"
	RouteDeleteSnippet = "/delete_snippet"
	RouteSnippets      = "/snippets"

	// tags
	RouteCreateTag = "/create_tag"
	RouteUpdateTag = "/update_tag"
	RouteDeleteTag = "/delete_tag"
	RouteTags      = "/tags"

	// contacts
	RouteCreateContact            = "/create_contact"
	RouteContacts                 = "/contacts"
	RouteDeleteContact            = "/delete_contact"
	RouteContact                  = "/contact"
	RouteUpdateContact            = "/update_contact"
	RouteImportContacts           = "/import_contacts"
	RouteExportContacts           = "/export_contacts"
	RouteExportContactsForProduct = "/export_contacts_for_product"
	RouteBlockContact             = "/block_contact"
	RouteUnblockContact           = "/unblock_contact"

	// emails configuration
	RouteEmailsConfiguration          = "/emails_configuration"
	RouteUpdateEmailsConfiguration    = "/update_emails_configuration"
	RouteVerifyEmailsDnsConfiguration = "/verify_emails_dns_configuration"

	// newsletters
	RouteNewsletters      = "/newsletters"
	RouteNewsletter       = "/newsletter"
	RouteCreateNewsletter = "/create_newsletter"
	RouteUpdateNewsletter = "/update_newsletter"
	RouteDeleteNewsletter = "/delete_newsletter"
	RouteSendNewsletter   = "/send_newsletter"

	// products
	RouteProduct                     = "/product"
	RouteUpdateProduct               = "/update_product"
	RouteCreateProduct               = "/create_product"
	RouteProducts                    = "/products"
	RouteGiveContactsAccessToProduct = "/give_contacts_access_to_product"
	RouteDeleteProduct               = "/delete_product"
	RouteRemoveAccessToProduct       = "/remove_access_to_product"

	// orders
	RouteOrders = "/orders"
	RouteOrder  = "/order"

	// refunds
	RouteRefunds      = "/refunds"
	RouteCreateRefund = "/create_refund"

	// product pages
	RouteCreateProductPage = "/create_product_page"
	RouteUpdateProductPage = "/update_product_page"
	RouteDeleteProductPage = "/delete_product_page"
	RouteProductPage       = "/product_page"

	// coupons
	RouteCreateCoupon = "/create_coupon"
	RouteUpdateCoupon = "/update_coupon"
	RouteCoupon       = "/coupon"
	RouteCoupons      = "/coupons"

	// analytics
	RouteAnalyticsData = "/analytics_data"
)
