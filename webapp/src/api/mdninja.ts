import { Routes } from "@/api/routes";
import { useStore } from "@/app/store";
import { type Router } from "vue-router";
import * as model from "@/api/model";
import type { Config } from "@/app/config";
import { usePingoo } from "@/pingoo/pingoo";


type ApiError = {
  message: string;
  code: string;
}

const API_BASE_URL = '/api';
const networkErrorMessage = 'Network error';
const accessTokenIsNotValidErrorMessage = 'Access Token is not valid'; // TODO
const INTERNAL_ERROR_MESSAGE = 'Internal Error. Please try again and contact support if the problem persists.';

let mdninjaService: MdninjaService | null = null;

export function createMdninja(config: Config, router: Router): MdninjaService {
  mdninjaService = new MdninjaService(config, router);
  return mdninjaService;
}

export function useMdninja(): MdninjaService {
  if (!mdninjaService) {
  throw new Error('mdninja service should be created before using it');
  }
  return mdninjaService!;
}

type postOptions = {
  authenticated?: boolean,
  headersInput?: Headers,
}

async function post<I>(route: string, data: I, options: postOptions = { authenticated: true }): Promise<any> {
  const url = `${API_BASE_URL}${route}`;
  let response: any = null;

  const headers = new Headers();
  headers.set('Content-Type', 'application/json');
  headers.set('Accept', 'application/json');

  if (options?.authenticated) {
    const pingoo = usePingoo();
    if (pingoo.isAuthenticated()) {
      headers.set('Authorization', `Bearer ${await pingoo.getAccessToken()}`);
    }
  }
  if (options?.headersInput) {
    for (let [headerName, headerValue] of options?.headersInput) {
      headers.set(headerName, headerValue);
    }
  }

  try {
    response = await fetch(url, {
      method: 'POST',
      headers,
      // For cross-origin requests, such as to an api.xx.com subdomain
      // credentials: 'include',
      body: JSON.stringify(data),
    });
  } catch (err: any) {
    throw new Error(networkErrorMessage);
  }

  return await unwrapApiResponse(response);
}

async function upload(route: string, data: FormData): Promise<any> {
  const url = `${API_BASE_URL}${route}`;
  let response: any = null;
  let pingoo = usePingoo();

  let headers = new Headers();
  headers.set('Accept', 'application/json');
  if (pingoo.isAuthenticated()) {
    headers.set('Authorization', `Bearer ${await pingoo.getAccessToken()}`);
  }

  try {
    response = await fetch(url, {
      method: 'POST',
      headers,
      // For cross-origin requests, such as to an api.xx.com subdomain
      // credentials: 'include',
      body: data,
    });
  } catch (err) {
    throw new Error(networkErrorMessage);
  }

  return await unwrapApiResponse(response);
}

async function unwrapApiResponse(response: Response): Promise<any>  {
  // if the status code is >= 500 or the response is not JSON then something has gone really wrong
  if (response.status >= 500 || !response.headers.get('Content-Type')?.includes('application/json')) {
    throw new Error(INTERNAL_ERROR_MESSAGE);
  }

  const apiRes: any = await response.json();

  if (response.status >= 400) {
    throw new Error((apiRes as ApiError).message);
  }

  return apiRes;
}

export async function importContacts(input: model.ImportContactsInput): Promise<model.Contact[]> {
  return await post(Routes.importContacts, input);
}

export class MdninjaService {
  private config: Config;

  constructor(config: Config, router: Router) {
    this.config = config;
  }

  async init() {
    const $store = useStore();
    const pingoo = usePingoo();

    try {
      if (pingoo.isAuthenticated()) {
        const [accessTokenClaims, organizations] = await Promise.all([
          pingoo.getAccessTokenClaims(),
          this.fetchOrganizationsForUser(),
        ]);
        $store.setOrganizations(organizations);
        $store.setUserId(accessTokenClaims.sub);
        $store.setIsAdmin(accessTokenClaims.is_admin);
        $store.setUserEmail(accessTokenClaims.email);
      }
    } catch (err: any) {
      console.error(err.message);
      $store.setAppLoadingError(err.message);

      if (err.message === accessTokenIsNotValidErrorMessage) {
        this.logout();
        return;
      }
    } finally {
      $store.setAppLoaded();
    }
  }

  async logout() {
    const pingoo = usePingoo();
    pingoo.logout();
    window.location.href = '/';
  }

  //////////////////////////////////////////////////////////////////////////////////////////////////
  // Analytics
  //////////////////////////////////////////////////////////////////////////////////////////////////

  async getAnalyticsData(input: model.GetAnalyticsDataInput): Promise<model.AnalyticsData> {
    const res: model.AnalyticsData = await post(Routes.analyticsData, input);

    return res;
  }

  //////////////////////////////////////////////////////////////////////////////////////////////////
  // Contacts
  //////////////////////////////////////////////////////////////////////////////////////////////////

  async createContact(input: model.CreateContactInput): Promise<model.Contact> {
    const res: model.Contact = await post(Routes.createContact, input);

    return res;
  }

  async listContacts(input: model.ListContactsInput): Promise<model.PaginatedResult<model.Contact>> {
    const res: model.PaginatedResult<model.Contact> = await post(Routes.contacts, input);

    return res;
  }

  async updateContact(input: model.UpdateContactInput): Promise<model.Contact> {
    const res: model.Contact = await post(Routes.updateContact, input);

    return res;
  }

  async deleteContact(id: string): Promise<void> {
    const input: model.DeleteContactInput = {
      id,
    };
    await post(Routes.deleteContact, input);
  }

  async fetchContact(contactId: string): Promise<model.Contact> {
    const input: model.GetContactInput = {
      id: contactId,
    };
    const res: model.Contact = await post(Routes.contact, input);

    return res;
  }

  async exportContacts(input: model.ExportContactsInput): Promise<model.ExportContactsOutput> {
    const res: model.ExportContactsOutput = await post(Routes.exportContacts, input);

    return res;
  }

  async exportContactsForProduct(input: model.ExportContactsForProductInput): Promise<model.ExportContactsForProductOutput> {
    const res: model.ExportContactsForProductOutput = await post(Routes.exportContactsForProduct, input);

    return res;
  }

  async blockContact(input: model.BlockContactInput): Promise<model.Contact> {
    const res: model.Contact = await post(Routes.blockContact, input);

    return res;
  }

  async unblockContact(input: model.UnblockContactInput): Promise<model.Contact> {
    const res: model.Contact = await post(Routes.unblockContact, input);

    return res;
  }

  //////////////////////////////////////////////////////////////////////////////////////////////////
  // Content
  //////////////////////////////////////////////////////////////////////////////////////////////////

  async fetchPage(pageID: string): Promise<model.Page> {
    const input: model.FetchPageInput = {
      id: pageID,
    };
    const content: model.Page = await post(Routes.page, input);

    return content;
  }

  async createPage(input: model.CreatePageInput): Promise<model.Page> {
    const page: model.Page = await post(Routes.createPage, input);
    return page;
  }

  async deletePage(pageID: string) {
    const input: model.DeletePageInput = {
      id: pageID,
    };
    await post(Routes.deletePage, input);
  }

  async updatePage(input: model.UpdatePageInput): Promise<model.Page> {
    const page: model.Page = await post(Routes.updatePage, input);
    return page;
  }

  async uploadAsset(input: model.UploadAssetInput): Promise<model.Asset> {
    const formData = new FormData();
    formData.append('website_id', input.website_id);
    formData.append('file', input.file);
    if (input.folder) {
      formData.append('folder', input.folder);
    }
    if (input.product_id) {
      formData.append('product_id', input.product_id);
    }

    const asset: model.Asset = await upload(Routes.uploadAsset, formData);

    return asset;
  }

  async deleteAsset(assetId: string) {
    const input: model.DeleteAssetInput = {
      id: assetId,
    };

    await post(Routes.deleteAsset, input);
  }

  async createAssetFolder(input: model.CreateFolderInput): Promise<model.Asset> {
    const newFolder: model.Asset = await post(Routes.createAssetFolder, input);
    return newFolder;
  }

  // generateVideoUrl(website: model.Website, asset: model.Asset): string {
  //   return `${location.protocol}//${website.primary_domain}${this.config.sitesPort}/__markdown_ninja/videos/${asset.id}/iframe`;
  // }

  generateAssetPathUrl(website: model.Website, asset: model.Asset): string {
    return `${location.protocol}//${website.primary_domain}${this.config.sitesPort}${asset.folder}/${asset.name}`;
  }

  generateAssetIdUrl(website: model.Website, assetId: string, download = false): string {
    const queryParameters = new URLSearchParams()
    queryParameters.append('id', assetId);
    if (download) {
      queryParameters.append('download', download.toString());
    }
    return `${location.protocol}//${website.primary_domain}${this.config.sitesPort}/assets?${queryParameters.toString()}`;
  }

  generatePagePreviewUrl(website: model.Website, page: model.Page): string {
    return `${location.protocol}//${website.primary_domain}${this.config.sitesPort}/__markdown_ninja/preview/${page.id}`;
  }

  generatePageUrl(website: model.Website, page: model.Page): string {
    return `${location.protocol}//${website.primary_domain}${this.config.sitesPort}${page.path}`;
  }

  generateWebsiteUrl(website: model.Website): string {
    return `${location.protocol}//${website.primary_domain}${this.config.sitesPort}`;
  }

  async deleteSnippet(snippetId: string): Promise<void> {
    const input: model.DeleteSnippetInput = {
      id: snippetId,
    };
    await post(Routes.deleteSnippet, input);
  }

  async createSnippet(input: model.CreateSnippetInput): Promise<model.Snippet> {
    const snippet: model.Snippet = await post(Routes.createSnippet, input);

    return snippet;
  }

  async updateSnippet(input: model.UpdateSnippetInput): Promise<model.Snippet> {
    const snippet: model.Snippet = await post(Routes.updateSnippet, input);

    return snippet;
  }

  async createTag(input: model.CreateTagInput): Promise<model.Tag> {
    const tag: model.Tag = await post(Routes.createTag, input);

    return tag;
  }

  async updateTag(input: model.UpdateTagInput): Promise<model.Tag> {
    const tag: model.Tag = await post(Routes.updateTag, input);

    return tag;
  }

  async deleteTag(tagId: string): Promise<void> {
    const input: model.DeleteTagInput = {
      id: tagId,
    };
    await post(Routes.deleteTag, input);
  }

  async fetchTags(input: model.GetTagsInput): Promise<model.Tag[]> {
    const res: model.Tag[] = await post(Routes.tags, input);

    return res;
  }

  async listPosts(input: model.ListPostsInput): Promise<model.PaginatedResult<model.PageMetadata>> {
    const res: model.PaginatedResult<model.PageMetadata> = await post(Routes.posts, input);
    return res;
  }

  async listPages(input: model.ListPagesInput): Promise<model.PaginatedResult<model.PageMetadata>> {
    const res: model.PaginatedResult<model.PageMetadata>= await post(Routes.pages, input);
    return res;
  }

  async listAssets(input: model.ListAssetsInput): Promise<model.Asset[]> {
    const res: model.Asset[] = await post(Routes.assets, input);

    return res;
  }

  async listSnippets(input: model.ListSnippetsInput): Promise<model.PaginatedResult<model.Snippet>> {
    const res: model.PaginatedResult<model.Snippet> = await post(Routes.snippets, input);

    return res;
  }


  //////////////////////////////////////////////////////////////////////////////////////////////////
  // Emails
  //////////////////////////////////////////////////////////////////////////////////////////////////

  async fetchConfiguration(websiteId: string): Promise<model.EmailConfiguration> {
    const input: model.GetEmailConfigurationInput = {
      website_id: websiteId,
    };
    const res: model.EmailConfiguration = await post(Routes.emailsConfiguration, input);

    return res;
  }

  async updateConfiguration(input: model.UpdateEmailConfigurationInput): Promise<model.EmailConfiguration> {
    const res: model.EmailConfiguration = await post(Routes.updateEmailsConfiguration, input);

    return res;
  }

  async verifyDnsConfiguraiton(websiteId: string): Promise<model.EmailConfiguration> {
    const input: model.VerifyDnsConfigurationInput = {
      website_id: websiteId,
    };
    const res: model.EmailConfiguration = await post(Routes.verifyEmailsConfiguration, input);

    return res;
  }

  async fetchNewsletters(websiteId: string): Promise<model.NewsletterMetadata[]> {
    const input: model.GetNewslettersInput = {
      website_id: websiteId,
    };
    const res: model.NewsletterMetadata[] = await post(Routes.newsletters, input);

    return res;
  }

  async fetchNewsletter(newsletterId: string): Promise<model.Newsletter> {
    const input: model.GetNewsletterInput = {
      id: newsletterId,
    };
    const res: model.Newsletter = await post(Routes.newsletter, input);

    return res;
  }

  async createNewsletter(input: model.CreateNewsletterInput): Promise<model.Newsletter> {
    const res: model.Newsletter = await post(Routes.createNewsletter, input);

    return res;
  }

  async updateNewsletter(input: model.UpdateNewsletterInput): Promise<model.Newsletter> {
    const res: model.Newsletter = await post(Routes.updateNewsletter, input);

    return res;
  }

  async deleteNewsletter(newsletterId: string) {
    const input: model.DeleteNewsletterInput = {
      id: newsletterId,
    };
    await post(Routes.deleteNewsletter, input);
  }

  async sendNewsletter(input: model.SendNewsletterInput): Promise<model.Newsletter> {
    return await post(Routes.sendNewsletter, input);
  }

  //////////////////////////////////////////////////////////////////////////////////////////////////
  // Kernel
  //////////////////////////////////////////////////////////////////////////////////////////////////
  async listFailedBackgroundJobs(): Promise<model.PaginatedResult<model.BackgroundJob>> {
    const res: model.PaginatedResult<model.BackgroundJob> = await post(Routes.failedBackgroundJobs, {});
    return res;
  }

  async deleteBackgroundJob(jobId: string) {
    const input: model.DeleteBackgroundJobInput = {
      id: jobId,
    };
    await post(Routes.deleteBackgroundJob, input);
  }

  //////////////////////////////////////////////////////////////////////////////////////////////////
  // Organizations
  //////////////////////////////////////////////////////////////////////////////////////////////////

  async fetchOrganizationsForUser(userId?: string): Promise<model.Organization[]> {
    const input: model.GetOrganizationsForUserInput = {
      // user_id: userId,
    };
    return await post(Routes.organizations, input);
  }

  async createOrganization(input: model.CreateOrganizationInput): Promise<model.CreateOrganizationOutput> {
    return await post(Routes.createOrganization, input);
  }

  async updateOrganization(input: model.UpdateOrganizationInput): Promise<model.Organization> {
    return await post(Routes.udpateOrganization, input);
  }

  async getOrganization(input: model.GetOrganizationInput): Promise<model.Organization> {
    const res: model.Organization = await post(Routes.organization, input);

    return res;
  }

  async deleteOrganization(organizationId: string) {
    const $store = useStore();

    const input: model.DeleteOrganizationInput = {
      id: organizationId,
    };
    await post(Routes.deleteOrganization, input);

    $store.deleteOrganization(organizationId);
  }

  async listOrganizations(input: model.ListOrganizationsInput): Promise<model.PaginatedResult<model.Organization>> {
    return await post(Routes.allorganizations, input);
  }

  async createApiKey(input: model.CreateApiKeyInput): Promise<model.ApiKey> {
    return await post(Routes.createApiKey, input);
  }

  async deleteApiKey(apiKeyID: string) {
    const input: model.DeleteApiKeyInput = {
      id: apiKeyID,
    };
    await post(Routes.deleteApiKey, input);
  }

  async listStaffInvitationsForOrganization(organizationId: string): Promise<model.PaginatedResult<model.StaffInvitation>> {
    const input: model.ListStaffInvitationsForOrganizationInput = {
      organization_id: organizationId,
    };
    return await post(Routes.staffInvitations, input);
  }

  async listUserInvitations(): Promise<model.PaginatedResult<model.UserInvitation>> {
    return await post(Routes.userInvitations, {});
  }

  async inviteStaffs(input: model.InviteStaffsInput): Promise<model.StaffInvitation> {
    return await post(Routes.inviteStaffs, input);
  }

  async deleteStaffInvitation(id: string) {
    const input: model.DeleteStaffInvitationInput = {
      id,
    }
    await post(Routes.deleteStaffInvitation, input);
  }

  async acceptStaffInvitation(id: string) {
    const input: model.AcceptStaffInvitationInput = {
      id,
    }
    await post(Routes.acceptStaffInvitation, input);
  }

  async removeStaff(input: model.RemoveStaffInput) {
    await post(Routes.removeStaff, input);
  }

  async organizationUpdateSubscription(input: model.OrganizationUpdateSubscriptionInput): Promise<model.OrganizationUpdateSubscriptionOutput> {
    return await post(Routes.organizationUpdateSubscription, input);
  }

  async organizationGetStripeCustomerPortal(input: model.OrganizationGetStripeCustomerPortalUrlInput): Promise<model.OrganizationGetStripeCustomerPortalUrlOutput> {
    return await post(Routes.organizationStripeCustomerPortal, input);
  }

  async organizationSyncStripe(organizationId: string) {
    await post<model.OrganizationSyncStripeInput>(Routes.organizationSyncStripe, { organization_id: organizationId });
  }

  async getorganizationBillingUsage(organizationId: string): Promise<model.OrganizationBillingUsage> {
    return await post<model.GetOrganizationBillingUsageInput>(Routes.organizationBillingUsage, { organization_id: organizationId });
  }

  async addStaffs(input: model.AddStaffs): Promise<model.Staff[]> {
    return await post(Routes.addStaffs, input);
  }

  async getOrganizationsAdminStatistics(): Promise<model.OrganizationsAdminStatistics> {
    return await post(Routes.organizationsAdminStatistics, {});
  }


  //////////////////////////////////////////////////////////////////////////////////////////////////
  // Store
  //////////////////////////////////////////////////////////////////////////////////////////////////
  async createProduct(input: model.CreateProductInput): Promise<model.Product> {
    const res: model.Product = await post(Routes.createProduct, input);

    return res;
  }

  async deleteProduct(input: model.DeleteProductInput) {
    await post(Routes.deleteProduct, input);
  }

  async updateProduct(input: model.UpdateProductInput): Promise<model.Product> {
    const res: model.Product = await post(Routes.updateProduct, input);

    return res;
  }

  async getProduct(productId: string): Promise<model.Product> {
    const input: model.GetProductInput = {
      id: productId,
    };
    const res: model.Product = await post(Routes.product, input);

    return res;
  }

  async listProducts(websiteId: string): Promise<model.PaginatedResult<model.Product>> {
    const input: model.ListProductsInput = {
      website_id: websiteId,
    };
    const res: model.PaginatedResult<model.Product> = await post(Routes.products, input);
    return res;
  }

  async createCoupon(input: model.CreateCouponInput): Promise<model.Coupon> {
    const res: model.Coupon = await post(Routes.createCoupon, input);

    return res;
  }

  async updateCoupon(input: model.UpdateCouponInput): Promise<model.Coupon> {
    const res: model.Coupon = await post(Routes.updateCoupon, input);

    return res;
  }

  async getCoupon(couponId: string): Promise<model.Coupon> {
    const input: model.GetCouponInput = {
      id: couponId,
    };
    const res: model.Coupon = await post(Routes.coupon, input);

    return res;
  }

  async listCoupons(websiteId: string): Promise<model.PaginatedResult<model.Coupon>> {
    const input: model.ListCouponsInput = {
      website_id: websiteId,
    };
    const res: model.PaginatedResult<model.Coupon> = await post(Routes.coupons, input);

    return res;
  }


  async createProductPage(input: model.CreateProductPageInput): Promise<model.ProductPage> {
    const res: model.ProductPage = await post(Routes.createProductPage, input);

    return res;
  }

  async getProductPage(pageId: string): Promise<model.ProductPage> {
    const input: model.GetProductPageInput = {
      id: pageId,
    };

    const res: model.ProductPage = await post(Routes.productPage, input);

    return res;
  }

  async deleteProductPage(pageId: string) {
    const input: model.DeleteProductPageInput = {
      id: pageId,
    };
    await post(Routes.deleteProductPage, input);
  }

  async updateProductPage(input: model.UpdateProductPageInput): Promise<model.ProductPage> {
    const res: model.ProductPage = await post(Routes.updateProductPage, input);

    return res;
  }

  async giveContactAccessToProducts(input: model.GiveContactsAccessToProductInput) {
    await post(Routes.giveContactsAccessToProduct, input);
  }

  async removeAccessToproduct(input: model.RemoveAccessToproduct) {
    await post(Routes.removeAccessToproduct, input);
  }

  async listOrders(input: model.ListOrdersInput): Promise<model.PaginatedResult<model.OrderMetadata>> {
    const res: model.PaginatedResult<model.Order> = await post(Routes.orders, input);

    return res;
  }

  async getOrder(orderId: string): Promise<model.Order> {
    const input: model.GetOrderInput = {
      id: orderId,
    };
    const res: model.Order = await post(Routes.order, input);

    return res;
  }

  async listRefunds(websiteId: string): Promise<model.PaginatedResult<model.Refund>> {
    const input: model.ListRefundsInput = {
      website_id: websiteId,
    };
    const res: model.PaginatedResult<model.Refund> = await post(Routes.refunds, input);

    return res;
  }

  async createRefund(input: model.CreateRefundInput): Promise<model.Refund> {
    const res: model.Refund = await post(Routes.createRefund, input);

    return res;
  }

  //////////////////////////////////////////////////////////////////////////////////////////////////
  // Websites
  //////////////////////////////////////////////////////////////////////////////////////////////////

  async addDomain(input: model.AddDomainInput): Promise<model.Domain> {
    const domain: model.Domain = await post(Routes.addDomain, input);

    return domain;
  }

  async removeDomain(domainId: string) {
    const input: model.RemoveDomainInput = {
      id: domainId,
    };
    await post(Routes.removeDomain, input);
  }

  async setDomainAsPrimary(input: model.SetDomainAsPrimaryInput) {
    await post(Routes.setDomainAsPrimary, input);
  }

  async checkTlsCertificateForDomain(domainId: string) {
    const input: model.CheckTlsCertificateForDomainInput = {
      domain_id: domainId,
    };
    await post(Routes.checkTlsCertificateForDomain, input);
  }

  async createWebsite(input: model.CreateWebsiteInput): Promise<model.Website> {
    const $store = useStore();
    const res: model.Website = await post(Routes.createWebsite, input);
    $store.addWebsites([res]);

    return res;
  }

  async updateWebsite(input: model.UpdateWebsiteInput): Promise<model.Website> {
    const res: model.Website = await post(Routes.updateWebsite, input);

    return res;
  }

  async deleteWebsite(input: model.DeleteWebsiteInput): Promise<void> {
    const $store = useStore();

    await post(Routes.deleteWebsite, input);
    $store.removeWebsites([input.id]);
  }

  async listWebsites(input: model.GetWebsitesForOrganizationInput): Promise<model.Website[]> {
    const $store = useStore();

    const res: model.Website[] = await post(Routes.websites, input);
    $store.addWebsites(res);

    return res;
  }

  async getWebsite(input: model.GetWebsiteInput): Promise<model.Website> {
    const $store = useStore();

    const res: model.Website = await post(Routes.website, input);
    $store.addWebsites([res]);

    return res;
  }

  async saveRedirects(input: model.SaveRedirectsInput): Promise<model.Redirect[]> {
    const res: model.Redirect[] = await post(Routes.saveRedirects, input);

    return res;
  }

  async listAllWebsites(input: model.ListWebsitesInput): Promise<model.PaginatedResult<model.Website>> {
    const res: model.PaginatedResult<model.Website> = await post(Routes.allWebsites, input);
    return res;
  }

  async updateWebsiteIcon(input: model.UpdateWebsiteIconInput) {
    const formData = new FormData();
    formData.append('website_id', input.website_id);
    formData.append('file', input.file);

    await upload(Routes.websiteUpdateIcon, formData);
  }

  async getWebsitessAdminStatistics(): Promise<model.WebsitesAdminStatistics> {
    return await post(Routes.websitesAdminStatistics, {});
  }
}

export async function getInitData() {
  const $store = useStore();

  let initData: model.InitData | null = null;
  if (window.__markdown_ninja_init_data) {
    initData = window.__markdown_ninja_init_data;
  } else {
    initData = await post(Routes.init, {}, { authenticated: false });
  }

  $store.setInitData(initData!);
}
