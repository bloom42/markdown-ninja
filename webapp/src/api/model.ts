export const MAX_ASSET_SIZE = 100_000_000;


////////////////////////////////////////////////////////////////////////////////////////////////////
// Analytics
////////////////////////////////////////////////////////////////////////////////////////////////////

export type Counter = {
  label: string;
  count: number;
};

export type AnalyticsData = {
  total_page_views: number;
  page_views: Counter[];
  total_visitors: number;
  visitors: Counter[];
  pages: Counter[];
  referrers: Counter[];
  countries: Counter[];
  browsers: Counter[];
  oses: Counter[];
  new_subscribers: number;
}

export type GetAnalyticsDataInput = {
  website_id: string;
}


////////////////////////////////////////////////////////////////////////////////////////////////////
// Contacts
////////////////////////////////////////////////////////////////////////////////////////////////////

export type Address = {
  country_code: string;
  line1: string;
  line2: string;
  postal_code: string;
  city: string;
  state: string;
}

export type Contact = {
  id: string;
  created_at: string;
  updated_at: string;
  name: string;
  email: string;
  country: string;
  subscribed_to_newsletter_at: string | null;
  blocked_at: string | null;

  stripe_customer_id: string | null;

  products: Product[] | null;
  orders: Order[] | null;
}

export type CreateContactInput = {
  website_id: string;
  email: string;
  name: string;
}

export type UpdateContactInput = {
  id: string;
  email?: string;
  name?: string;
  subscribed_to_newsletter?: boolean;
}

export type DeleteContactInput = {
  id: string;
}

export type ListContactsInput = {
  website_id: string;
  query?: string;
}

export type GetContactInput = {
  id: string;
}

export type ImportContactsInput = {
  website_id: string;
  contacts: string;
}

export type ExportContactsInput = {
  website_id: string;
}

export type ExportContactsOutput = {
  contacts: string;
}

export type ExportContactsForProductInput = {
  product_id: string;
}

export type ExportContactsForProductOutput = {
  contacts: string;
}

export type BlockContactInput = {
  id: string;
}

export type UnblockContactInput = {
  id: string;
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Content
////////////////////////////////////////////////////////////////////////////////////////////////////

export enum PageType {
  Page = "page",
  Post = "post",
};

export enum PageStatus {
  Published = "published",
  Draft = "draft",
  Scheduled = "scheduled",
};


export enum AssetType {
  Image = "image",
  Video = "video",
  Audio = "audio",
  File = "file",
  Folder = "folder",
};

export enum VideoStatus {
  Created = "created",
  Uploading = "uploading",
  Uploaded = "uploaded",
  Ready = "ready",
  Transcoding = "transcoding",
  Error = "error",
};

export interface Page extends PageMetadata {
  description: string;
  body_markdown: string;

  tags: Tag[];
}

export interface PageMetadata {
  id: string;
  created_at: string;
  updated_at: string;
  date: string;
  type: PageType;
  path: string;
  title: string;
  size: number;
  body_hash: string;
  metadata_hash: string;
  status: PageStatus;
  language: string;
  send_as_newsletter: boolean;
  newsletter_sent_at: string | null;
};

export type Asset = {
  id: string;
  created_at: string;
  updated_at: string;

  type: AssetType,
  name: string;
  folder: string;
  media_type: string;
  size: number;
  hash: string;
}


export type Tag = {
  id: string;
  created_at: string;
  updated_at: string;
  name: string;
  description: string;
};

export type Snippet = {
  id: string;
  created_at: string;
  updated_at: string;

  name: string;
  content: string;
  hash: string;
  render_in_emails: boolean;
}

export type CreateTagInput = {
  website_id: string;
  name: string;
  description: string;
}

export type UpdateTagInput = {
  id: string;
  name: string;
  description: string;
}

export type DeleteTagInput = {
  id: string;
}

export type CreateSnippetInput = {
  website_id: string;
  name: string;
  content: string;
  render_in_emails: boolean;
}

export type UpdateSnippetInput = {
  id: string;
  name: string;
  content: string;
  render_in_emails: boolean;
}

export type DeleteSnippetInput = {
  id: string;
}

export type CreatePageInput = {
  website_id: string;
  date: string;
  type: PageType;
  title: string;
  path: string;
  tags: string[];
  description: string;
  language: string;
  draft: boolean;
  body_markdown: string;
  send_as_newsletter: boolean;
}

export type UpdatePageInput = {
  id: string;
  date: string;
  title: string;
  path: string;
  tags: string[];
  description?: string;
  language: string;
  draft: boolean;
  body_markdown?: string;
  send_as_newsletter: boolean;
}

export type DeletePageInput = {
  id: string;
}

export type DeleteAssetInput = {
  id: string;
}

export type FetchPageInput = {
  id: string;
}

export type GetTagsInput = {
  website_id: string;
}

export type ListPostsInput = {
  website_id: string;
}

export type ListPagesInput = {
  website_id: string;
}

export type ListAssetsInput = {
  website_id: string;
  folder?: string;
}

export type ListSnippetsInput = {
  website_id: string;
}

export type UploadAssetInput = {
  file: File,
  website_id: string;
  folder?: string;
  product_id?: string;
}

export type CreateFolderInput = {
  website_id: string;
  folder: string;
  name: string;
}


////////////////////////////////////////////////////////////////////////////////////////////////////
// Emails
////////////////////////////////////////////////////////////////////////////////////////////////////


export interface NewsletterMetadata {
  id: string;
  created_at: string;
  updated_at: string;
  scheduled_for: string | null;
  subject: string;
  size: number;
  hash: string;
  sent_at: string | null;
  last_test_sent_at: string | null;
}

export interface Newsletter extends NewsletterMetadata {
  body_markdown: string;
}

export type EmailConfiguration = {
  from_name: string;
  from_address: string;
  domain_verified: string;
  dns_records: EmailDnsRecord[];
}

export type EmailDnsRecord = {
  host: string,
  type: string,
  value: string,
}

export type GetNewslettersInput = {
  website_id: string;
}

export type GetNewsletterInput = {
  id: string;
}

export type CreateNewsletterInput = {
  website_id: string;
  subject: string;
  scheduled_for?: string;
  body_markdown: string;
}

export type UpdateNewsletterInput = {
  id: string;
  subject: string;
  scheduled_for?: string;
  body_markdown?: string;
}

export type DeleteNewsletterInput = {
  id: string;
}

export type UpdateEmailConfigurationInput = {
  website_id: string;
  from_name: string;
  from_address: string;
}

export type GetEmailConfigurationInput = {
  website_id: string;
}

export type VerifyDnsConfigurationInput = {
  website_id: string;
}

export type SendNewsletterInput = {
  id: string;
  test?: boolean;
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Kernel
////////////////////////////////////////////////////////////////////////////////////////////////////

export type PaginatedResult<T> = {
  data: T[];
  // total: number;
}

export enum TwoFaMethod {
  Totp = "totp",
};

export type PricingPlan = {
  id: string,
  name: string,
  description: string,
  price: number,
  features: string[],
}

export type BackgroundJob = {
  id: string;
  created_at: string;
  updated_at: string;
  scheduled_for: string;
  failed_attempts: number;
  status: string;
  type: string;
  data: string;
  retry_max: number
  retry_delay: number
  retry_strategy: string;
  timeout: string;
}


export type InitData = {
  stripe_public_key: string;
  country: string;
  contact_email: string;
  pricing: PricingPlan[],
  // challenge_site_key: string | null;
  pingoo: PingooInitData,
  websites_base_url: string;
}

export type PingooInitData = {
  endpoint: string,
  app_id: string;
}

export type DeleteBackgroundJobInput = {
  id: string;
}


////////////////////////////////////////////////////////////////////////////////////////////////////
// Organizations
////////////////////////////////////////////////////////////////////////////////////////////////////

export enum StaffRole {
  Administrator = "administrator",
};


export type Organization = {
  id: string;
  created_at: string;
  updated_at: string;
  name: string;

  plan: string;
  billing_information: BillingInformation;
  extra_slots: number;
  stripe_customer: boolean;
  payment_due: boolean;

  api_keys: ApiKey[] | null;
  staffs: Staff[] | null;
}

export type BillingInformation = {
  name: string;
  email: string;
  address_line1: string;
  address_line2: string;
  postal_code: string;
  city: string;
  state: string;
  country_code: string;
  tax_id: string | null;
};

export type StaffInvitation = {
  id: string;
  created_at: string;
  updated_at: string;
  role: string;
  invitee_email: string;
  organization_id: string;
  inviter_id: string;
}

export type UserInvitation = {
  id: string;
  created_at: string;
  inviter_name: string;
  inviter_email: string;
  organization_id: string;
  organization_name: string;
  invitee_email: string;
}


export type GetOrganizationsForUserInput = {
  user_id?: string;
}

export type CreateOrganizationInput = {
  name: string;
  plan: string;
  billing_email: string;
}

export type CreateOrganizationOutput = {
  organization: Organization;
  stripe_checkout_session_url: string | null;
}

export type UpdateOrganizationInput = {
  id: string;
  name?: string;
  billing_information?: BillingInformation;

  plan?: string;
  extra_slots?: number;
}

export type GetOrganizationInput = {
  id: string;
  api_keys?: boolean;
  staffs?: boolean;
}

export type GetStaffsInput = {
  organization_id: string;
}

export type Staff = {
  created_at: string;
  role: StaffRole;
  name: string;
  email: string;
  user_id: string;
  organization_id: string;
}

export type DeleteOrganizationInput = {
  id: string;
}

export type ListOrganizationsInput = {
}

export type DeleteApiKeyInput = {
  id: string;
}

export type CreateApiKeyInput = {
  organization_id: string;
  name: string;
}

export type InviteStaffsInput = {
  organization_id: string;
  emails: string[];
}

export type AcceptStaffInvitationInput = {
  id: string;
}

export type DeleteStaffInvitationInput = {
  id: string;
}

export type ListStaffInvitationsForOrganizationInput = {
  organization_id: string;
}

export type RemoveStaffInput = {
  organization_id: string;
  user_id: string;
}

export type OrganizationUpdateSubscriptionInput = {
  organization_id: string;
  plan: string;
  extra_slots: number;
}

export type OrganizationUpdateSubscriptionOutput = {
  stripe_checkout_session_url: string | null;
}

export type OrganizationGetStripeCustomerPortalUrlInput = {
  organization_id: string;
}

export type OrganizationGetStripeCustomerPortalUrlOutput = {
  stripe_customer_portal_url: string;
}

export type OrganizationSyncStripeInput = {
  organization_id: string;
}

export type GetOrganizationBillingUsageInput = {
  organization_id: string;
}

export type OrganizationBillingUsage = {
  used_websites: number;
  allowed_websites: number;
  used_storage: number;
  allowed_storage: number;
  used_staffs: number;
  allowed_staffs: number;
  allowed_emails: number;
  used_emails: number;
}

export type AddStaffs = {
  organization_id: string,
  user_ids: string[],
}

export type OrganizationsAdminStatistics = {
  organizations: number,
  paying_organizations: number,
  monthly_revenue: number,
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Products
////////////////////////////////////////////////////////////////////////////////////////////////////


export enum ProductType {
  Book = "book",
  Course = "course",
  Download = "download",
};

export enum ProductStatus {
  Draft = "draft",
  Active = "active",
};

export enum OrderStatus {
  Pending = "pending",
  Completed = "completed",
  Canceled = "canceled",
};

export enum RefundReason {
  Duplicate = "duplicate",
  Fraudulent = "fraudulent",
  RequestedByCustomer = "requested_by_customer",
  // expired_uncaptured_charge
};


export type Product = {
  id: string;
  created_at: string;
  updated_at: string;

  name: string;
  description: string;
  type: ProductType;
  status: ProductStatus;
  price: number;

  content: ProductPage[] | null;
  assets: Asset[] | null;
}

export type Order = {
  id: string;
  created_at: string;
  updated_at: string;

  total_amount: number;
  currency: string;
  status: OrderStatus;
  country: string;
  email: string;
  additional_invoice_information: string;
  stripe_checkout_session_id: string;
  stripe_payment_intent_id?: string;
  stripe_invoice_id?: string;
  stripe_invoice_url?: string;

  line_items?: OrderLineItem[];
  contact_id: string;
  refunds: Refund[] | null;
}

export type OrderMetadata = {
  id: string;
  created_at: string;
  updated_at: string;

  total_amount: number;
  currency: string;
  status: OrderStatus;
  country: string;
  email: string;
  completed_at?: string;
  canceled_at?: string;

  contact_id: string;
}

export type OrderLineItem = {
  product_name: string;
  original_product_price: number;
  quantity: number;
  product_id: string;
}

export type Coupon = {
  id: string;
  created_at: string;
  updated_at: string;

  code: string;
  expires_at: string | null;
  discount: number;
  uses_limit: number;
  archived: boolean;
  description: string;
  products: string[];
}



export type ProductPage = {
  id: string;
  created_at: string;
  updated_at: string;

  position: number;
  title: string;
  size: number;
  hash: string;
  body_markdown: string;
};

export type Refund = {
  id: string;
  created_at: string;
  updated_at: string;

  amount: number;
  currency: string;
  notes: string;
  status: string;
  reason: RefundReason;
  failure_reason: string | null;
  stripe_refund_id: string | null;

  website_id: string;
  order_id: string;
};

export type CreateProductInput = {
  website_id: string;
  name: string;
  description: string;
  type: ProductType;
  price: number;
}

export type GetProductInput = {
  id: string;
}

export type DeleteProductInput = {
  id: string;
}

export type ListProductsInput = {
  website_id: string;
}

export type ListOrdersInput = {
  website_id: string;
  query?: string;
  limit?: number;
  after?: string;
}

export type UpdateProductInput = {
  id: string;
  name?: string;
  description?: string;
  status?: ProductStatus;
  price?: number;
}

export type CreateCouponInput = {
  website_id: string;
  code: string;
  description: string;
  expires_at?: string;
  discount: number;
  products: string[];
}

export type GetCouponInput = {
  id: string;
}

export type ListCouponsInput = {
  website_id: string;
}

export type UpdateCouponInput = {
  id: string;
  code?: string;
  description?: string;
  expires_at?: string;
  discount?: number;
  archived?: boolean;
  products?: string[];
}


export type CreateProductPageInput = {
  product_id: string;
  title: string;
  body_markdown: string;
};

export type UpdateProductPageInput = {
  id: string;
  title?: string;
  body_markdown?: string;
};


export type DeleteProductPageInput = {
  id: string;
};

export type GetProductPageInput = {
  id: string;
};

export type GiveContactsAccessToProductInput = {
  product_id: string;
  emails: string[];
}

export type RemoveAccessToproduct = {
  product_id: string;
  emails: string[];
}

export type GetOrderInput = {
  id: string;
}

export type ListRefundsInput = {
  website_id: string;
}

export type CreateRefundInput = {
  order_id: string;
  reason: RefundReason;
  notes: string;
  amount: number;
};

////////////////////////////////////////////////////////////////////////////////////////////////////
// Websites
////////////////////////////////////////////////////////////////////////////////////////////////////

export const allCurrencies = ["USD", "EUR"];
export const builtInThemes = ["blog", "docs"];

export type Website = {
  id: string;
  created_at: string;
  updated_at: string;

  name: string;
  slug: string;
  header: string;
  footer: string;
  navigation: WebsiteNavigation;
  language: string;
  used_storage: number;
  primary_domain: string;
  description: string;
  robots_txt: string;
  organization_id: string;
  blocked_at?: string;
  currency: string;
  colors: ThemeColors;
  theme: string;
  ad: string | null;
  announcement: string | null;
  logo: string | null;
  powered_by: boolean,

  domains: Domain[] | null;
  redirects: Redirect[] | null;
  revenue: number | null;
  subscribers: number | null;
};

export type ThemeColors = {
  background: string;
  text: string;
  accent: string;
}

export type Redirect = {
  id: string;
  created_at: string;
  updated_at: string;
  domain: string;
  pattern: string;
  path_pattern: string;
  to: string;
  status: number;
}

export type WebsiteNavigation = {
  primary: SiteNavigationItem[],
  secondary: SiteNavigationItem[],
}

export type SiteNavigationItem = {
  url: string;
  label: string;
};

export type ApiKey = {
  id: string;
  created_at: string;
  updated_at: string;
  name: string;
  token?: string;
};

export type Domain = {
  id: string;
  created_at: string;
  updated_at: string;

  hostname: string;
  tls_active: boolean;
};

export type AddDomainInput = {
  website_id: string;
  hostname: string;
  primary: boolean;
}

export type RemoveDomainInput = {
  id: string;
}

export type SetDomainAsPrimaryInput = {
  domain: string | null;
  website_id: string;
}

export type SaveRedirectsInput = {
  website_id: string;
  redirects: RedirectInput[],
}

export type RedirectInput = {
  pattern: string;
  to: string;
  // status: number;
}


export type CreateWebsiteInput = {
  name: string;
  slug: string;
  organization_id: string;
  currency?: string;
}

export type UpdateWebsiteInput = {
  id: string;
  header?: string;
  footer?: string;
  name?: string;
  description?: string;
  slug?: string;
  navigation?: WebsiteNavigation;
  robots_txt?: string;
  blocked?: boolean;
  currency?: string;
  background_color?: string;
  text_color?: string;
  accent_color?: string;
  theme?: string;
  ad?: string;
  announcement?: string;
  logo?: string;
  powered_by?: boolean,
}

export type DeleteWebsiteInput = {
  id: string;
}

export type GetWebsitesForOrganizationInput = {
  organization_id: string;
}

export type GetWebsiteInput = {
  id: string;
  domains?: boolean;
  redirects?: boolean;
}

export type CheckTlsCertificateForDomainInput = {
  domain_id: string;
}

export type ListWebsitesInput = {
  query?: string;
}

export type UpdateWebsiteIconInput = {
  file: File,
  website_id: string;
}
