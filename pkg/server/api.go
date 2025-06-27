package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"markdown.ninja/pkg/server/api"
	"markdown.ninja/pkg/server/apiutil"
)

// type Api struct {
// 	kernelService        kernel.Service
// 	websitesService      websites.Service
// 	contactsService      contacts.Service
// 	emailsService        emails.Service
// 	storeService         store.Service
// 	eventsService        events.Service
// 	contentService       content.Service
// 	organizationsService organizations.Service

// 	// env          config.Env
// 	webappDomain string
// }

// func NewApi(webappDomain string, kernelService kernel.Service, websitesService websites.Service,
// 	contactsService contacts.Service, emailsService emails.Service, storeService store.Service,
// 	eventsService events.Service, contentService content.Service, organizationsService organizations.Service,
// ) *Api {
// 	return &Api{
// 		kernelService,
// 		websitesService,
// 		contactsService,
// 		emailsService,
// 		storeService,
// 		eventsService,
// 		contentService,
// 		organizationsService,

// 		webappDomain,
// 	}
// }

func (server *server) Api(apiRouter chi.Router) {
	// apiRouter = chi.NewRouter()

	// we AllowCredentials to enable cookies to be set from an API subdomain
	// allowedDomain := fmt.Sprintf("https://%s", server.webappDomain)
	// if server.env == config.EnvDev {
	// 	allowedDomain = fmt.Sprintf("http://%s:3000", server.webappDomain)
	// }
	// cors := cors.New(cors.Options{
	// 	// we can't set Access-Control-Allow-Origin to "*" if we Allow Credentials
	// 	AllowedOrigins: []string{allowedDomain},
	// 	AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	// 	AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "Origin", "Set-Cookie", "Cookie"},
	// 	ExposedHeaders: []string{},
	// 	// to let Cookies flow
	// 	AllowCredentials: true,
	// 	MaxAge:           3600,
	// })
	// apiRouter.Use(cors.)

	apiRouter.Use(middleware.NoCache)
	// router.Use(middleware.MaxBodySize(server.MaxBodySize))

	apiRouter.Get("/", apiutil.IndexHandler)
	apiRouter.NotFound(apiutil.NotFoundHandler)

	////////////////////////////////////////////////////////////////////////////////////////////////
	// Webhooks
	////////////////////////////////////////////////////////////////////////////////////////////////
	apiRouter.Post(api.RouteWebhooksStripe, server.stripeWebhook)
	apiRouter.Post(api.RouteWebhooksPingoo, server.pingooWebhookHandler)

	////////////////////////////////////////////////////////////////////////////////////////////////
	// Kernel
	////////////////////////////////////////////////////////////////////////////////////////////////

	// misc
	apiRouter.Post(api.RouteInit, apiutil.JsonEndpoint(server.kernelService.Init))
	apiRouter.Post("/list_tls_certificates", apiutil.JsonEndpoint(server.certManager.ListCertificates))

	// jobs queue
	apiRouter.Post(api.RouteFailedBackgroundJobs, apiutil.JsonEndpoint(server.kernelService.ListFailedBackgroundJobs))
	apiRouter.Post(api.RouteDeleteBackgroundJob, apiutil.JsonEndpointOk(server.kernelService.DeleteBackgroundJob))

	// organizations
	apiRouter.Post(api.RouteOrganization, apiutil.JsonEndpoint(server.organizationsService.GetOrganization))
	apiRouter.Post(api.RouteOrganizations, apiutil.JsonEndpoint(server.organizationsService.GetOrganizationsForUser))
	apiRouter.Post(api.RouteCreateOrganization, apiutil.JsonEndpoint(server.organizationsService.CreateOrganization))
	apiRouter.Post(api.RouteUpdateOrganization, apiutil.JsonEndpoint(server.organizationsService.UpdateOrganization))
	apiRouter.Post(api.RouteDeleteOrganization, apiutil.JsonEndpointOk(server.organizationsService.DeleteOrganization))
	apiRouter.Post(api.RouteAllOrganizations, apiutil.JsonEndpoint(server.organizationsService.ListOrganizations))
	apiRouter.Post(api.RouteOrganizationUpdateSubscription, apiutil.JsonEndpoint(server.organizationsService.UpdateSubscription))
	apiRouter.Post(api.RouteOrganizationStripeCustomerPortal, apiutil.JsonEndpoint(server.organizationsService.GetStripeCustomerPortalUrl))
	apiRouter.Post(api.RouteOrganizationSyncStripe, apiutil.JsonEndpointOk(server.organizationsService.SyncStripe))
	apiRouter.Post(api.RouteOrganizationBillingUsage, apiutil.JsonEndpoint(server.organizationsService.GetBillingUsage))
	apiRouter.Post(api.RouteOrganizationsAdminStatistics, apiutil.JsonEndpoint(server.organizationsService.GetAdminStatistics))

	// staffs
	apiRouter.Post(api.RouteInviteStaffs, apiutil.JsonEndpoint(server.organizationsService.InviteStaffs))
	apiRouter.Post(api.RouteAddStaffs, apiutil.JsonEndpoint(server.organizationsService.AddStaffs))
	apiRouter.Post(api.RouteDeleteStaffInvitation, apiutil.JsonEndpointOk(server.organizationsService.DeleteStaffInvitation))
	apiRouter.Post(api.RouteStaffInvitations, apiutil.JsonEndpoint(server.organizationsService.ListStaffInvitationsForOrganization))
	apiRouter.Post(api.RouteUserInvitations, apiutil.JsonEndpoint(server.organizationsService.ListUserInvitations))
	apiRouter.Post(api.RouteAcceptStaffInvitation, apiutil.JsonEndpointOk(server.organizationsService.AcceptStaffInvitation))
	apiRouter.Post(api.RouteRemoveStaff, apiutil.JsonEndpointOk(server.organizationsService.RemoveStaff))

	// apiKeys
	apiRouter.Post(api.RouteCreateApiKey, apiutil.JsonEndpoint(server.organizationsService.CreateApiKey))
	apiRouter.Post(api.RouteDeleteApiKey, apiutil.JsonEndpointOk(server.organizationsService.DeleteApiKey))
	apiRouter.Post(api.RouteUpdateApiKey, apiutil.JsonEndpoint(server.organizationsService.UpdateApiKey))

	////////////////////////////////////////////////////////////////////////////////////////////////
	// Websites
	////////////////////////////////////////////////////////////////////////////////////////////////

	// websites
	apiRouter.Post(api.RouteCreateWebsite, apiutil.JsonEndpoint(server.websitesService.CreateWebsite))
	apiRouter.Post(api.RouteDeleteWebsite, apiutil.JsonEndpointOk(server.websitesService.DeleteWebsite))
	apiRouter.Post(api.RouteWebsites, apiutil.JsonEndpoint(server.websitesService.GetWebsitesForOrganization))
	apiRouter.Post(api.RouteWebsite, apiutil.JsonEndpoint(server.websitesService.GetWebsite))
	apiRouter.Post(api.RouteUpdateWebsite, apiutil.JsonEndpoint(server.websitesService.UpdateWebsite))
	apiRouter.Post(api.RouteSaveRedirect, apiutil.JsonEndpoint(server.websitesService.SaveRedirects))
	apiRouter.Post(api.RouteAllWebsites, apiutil.JsonEndpoint(server.websitesService.ListWebsites))
	apiRouter.Post(api.RouteWebsiteUpdateIcon, server.websiteUpdateIcon)
	apiRouter.Post(api.RouteWebsitesAdminStatistics, apiutil.JsonEndpoint(server.websitesService.GetAdminStatistics))

	// snippets
	apiRouter.Post(api.RouteCreateSnippet, apiutil.JsonEndpoint(server.contentService.CreateSnippet))
	apiRouter.Post(api.RouteDeleteSnippet, apiutil.JsonEndpointOk(server.contentService.DeleteSnippet))
	apiRouter.Post(api.RouteUpdateSnippet, apiutil.JsonEndpoint(server.contentService.UpdateSnippet))
	apiRouter.Post(api.RouteSnippets, apiutil.JsonEndpoint(server.contentService.ListSnippets))

	// tags
	apiRouter.Post(api.RouteCreateTag, apiutil.JsonEndpoint(server.contentService.CreateTag))
	apiRouter.Post(api.RouteUpdateTag, apiutil.JsonEndpoint(server.contentService.UpdateTag))
	apiRouter.Post(api.RouteDeleteTag, apiutil.JsonEndpointOk(server.contentService.DeleteTag))
	apiRouter.Post(api.RouteTags, apiutil.JsonEndpoint(server.contentService.GetTags))

	// pages
	apiRouter.Post(api.RouteCreatePage, apiutil.JsonEndpoint(server.contentService.CreatePage))
	apiRouter.Post(api.RouteUpdatePage, apiutil.JsonEndpoint(server.contentService.UpdatePage))
	apiRouter.Post(api.RouteDeletePage, apiutil.JsonEndpointOk(server.contentService.DeletePage))
	apiRouter.Post(api.RoutePage, apiutil.JsonEndpoint(server.contentService.GetPage))
	apiRouter.Post(api.RoutePages, apiutil.JsonEndpoint(server.contentService.ListPages))
	apiRouter.Post(api.RoutePosts, apiutil.JsonEndpoint(server.contentService.ListPosts))

	// assets
	apiRouter.Post(api.RouteUploadAsset, server.uploadAsset)
	apiRouter.Post(api.RouteDeleteAsset, apiutil.JsonEndpointOk(server.contentService.DeleteAsset))
	apiRouter.Post(api.RouteAssets, apiutil.JsonEndpoint(server.contentService.ListAssets))
	apiRouter.Post(api.RouteCreateAssetFolder, apiutil.JsonEndpoint(server.contentService.CreateAssetFolder))

	// domains
	apiRouter.Post(api.RouteAddDomain, apiutil.JsonEndpoint(server.websitesService.AddDomain))
	apiRouter.Post(api.RouteRemoveDomain, apiutil.JsonEndpointOk(server.websitesService.RemoveDomain))
	apiRouter.Post(api.RouteSetDomainAsPrimary, apiutil.JsonEndpointOk(server.websitesService.SetDomainAsPrimary))
	apiRouter.Post(api.RouteCheckTlsCertificateForDomain, apiutil.JsonEndpointOk(server.websitesService.CheckTlsCertificateForDomain))

	////////////////////////////////////////////////////////////////////////////////////////////////
	// Contacts
	////////////////////////////////////////////////////////////////////////////////////////////////

	// contacts
	apiRouter.Post(api.RouteCreateContact, apiutil.JsonEndpoint(server.contactsService.CreateContact))
	apiRouter.Post(api.RouteContacts, apiutil.JsonEndpoint(server.contactsService.ListContacts))
	apiRouter.Post(api.RouteDeleteContact, apiutil.JsonEndpointOk(server.contactsService.DeleteContact))
	apiRouter.Post(api.RouteContact, apiutil.JsonEndpoint(server.contactsService.GetContact))
	apiRouter.Post(api.RouteUpdateContact, apiutil.JsonEndpoint(server.contactsService.UpdateContact))
	apiRouter.Post(api.RouteImportContacts, apiutil.JsonEndpoint(server.contactsService.ImportContacts))
	apiRouter.Post(api.RouteExportContacts, apiutil.JsonEndpoint(server.contactsService.ExportContacts))
	apiRouter.Post(api.RouteExportContactsForProduct, apiutil.JsonEndpoint(server.contactsService.ExportContactsForProduct))
	apiRouter.Post(api.RouteBlockContact, apiutil.JsonEndpoint(server.contactsService.BlockContact))
	apiRouter.Post(api.RouteUnblockContact, apiutil.JsonEndpoint(server.contactsService.UnblockContact))

	////////////////////////////////////////////////////////////////////////////////////////////////
	// Emails
	////////////////////////////////////////////////////////////////////////////////////////////////

	// configuration
	apiRouter.Post(api.RouteUpdateEmailsConfiguration, apiutil.JsonEndpoint(server.emailsService.UpdateWebsiteConfiguration))
	apiRouter.Post(api.RouteEmailsConfiguration, apiutil.JsonEndpoint(server.emailsService.GetWebsiteConfiguration))
	apiRouter.Post(api.RouteVerifyEmailsDnsConfiguration, apiutil.JsonEndpoint(server.emailsService.VerifyDnsConfiguration))

	// newsletters
	apiRouter.Post(api.RouteNewsletters, apiutil.JsonEndpoint(server.emailsService.GetNewsletters))
	apiRouter.Post(api.RouteNewsletter, apiutil.JsonEndpoint(server.emailsService.GetNewsletter))
	apiRouter.Post(api.RouteCreateNewsletter, apiutil.JsonEndpoint(server.emailsService.CreateNewsletter))
	apiRouter.Post(api.RouteUpdateNewsletter, apiutil.JsonEndpoint(server.emailsService.UpdateNewsletter))
	apiRouter.Post(api.RouteDeleteNewsletter, apiutil.JsonEndpointOk(server.emailsService.DeleteNewsletter))
	apiRouter.Post(api.RouteSendNewsletter, apiutil.JsonEndpoint(server.emailsService.SendNewsletter))

	////////////////////////////////////////////////////////////////////////////////////////////////
	// Store
	////////////////////////////////////////////////////////////////////////////////////////////////

	// products
	apiRouter.Post(api.RouteCreateProduct, apiutil.JsonEndpoint(server.storeService.CreateProduct))
	apiRouter.Post(api.RouteUpdateProduct, apiutil.JsonEndpoint(server.storeService.UpdateProduct))
	apiRouter.Post(api.RouteProduct, apiutil.JsonEndpoint(server.storeService.GetProduct))
	apiRouter.Post(api.RouteProducts, apiutil.JsonEndpoint(server.storeService.ListProducts))
	apiRouter.Post(api.RouteGiveContactsAccessToProduct, apiutil.JsonEndpointOk(server.storeService.GiveContactsAccessToProduct))
	apiRouter.Post(api.RouteDeleteProduct, apiutil.JsonEndpointOk(server.storeService.DeleteProduct))
	apiRouter.Post(api.RouteRemoveAccessToProduct, apiutil.JsonEndpointOk(server.storeService.RemoveAccessToProduct))

	// orders
	apiRouter.Post(api.RouteOrders, apiutil.JsonEndpoint(server.storeService.ListOrders))
	apiRouter.Post(api.RouteOrder, apiutil.JsonEndpoint(server.storeService.GetOrder))

	// refunds
	apiRouter.Post(api.RouteRefunds, apiutil.JsonEndpoint(server.storeService.ListRefunds))
	apiRouter.Post(api.RouteCreateRefund, apiutil.JsonEndpoint(server.storeService.CreateRefund))

	// coupons
	apiRouter.Post(api.RouteCreateCoupon, apiutil.JsonEndpoint(server.storeService.CreateCoupon))
	apiRouter.Post(api.RouteUpdateCoupon, apiutil.JsonEndpoint(server.storeService.UpdateCoupon))
	apiRouter.Post(api.RouteCoupon, apiutil.JsonEndpoint(server.storeService.GetCoupon))
	apiRouter.Post(api.RouteCoupons, apiutil.JsonEndpoint(server.storeService.ListCoupons))

	// courses
	apiRouter.Post(api.RouteCreateProductPage, apiutil.JsonEndpoint(server.storeService.CreateProductPage))
	apiRouter.Post(api.RouteUpdateProductPage, apiutil.JsonEndpoint(server.storeService.UpdateProductPage))
	apiRouter.Post(api.RouteDeleteProductPage, apiutil.JsonEndpointOk(server.storeService.DeleteProductPage))
	apiRouter.Post(api.RouteProductPage, apiutil.JsonEndpoint(server.storeService.GetProductPage))

	////////////////////////////////////////////////////////////////////////////////////////////////
	// Analytics
	////////////////////////////////////////////////////////////////////////////////////////////////
	apiRouter.Post(api.RouteAnalyticsData, apiutil.JsonEndpoint(server.eventsService.GetAnalyticsData))
}
