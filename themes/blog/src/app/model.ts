export const AUTH_CODE_LENGTH  = 8;

export enum PageType {
  Page = "page",
  Post = "post",
};

export enum OrderStatus {
  Pending = "pending",
  Completed = "completed",
  Canceled = "canceled",
};

export enum ProductType {
  Book = "book",
  Course = "course",
  Download = "download",
};

export enum BlockType {
  Html = "html",
};

export type PaginatedResult<T> = {
  data: T[];
}

export type MarkdownNinjaData = {
  website: Website;
  contact: Contact | null;
  page: Page | null;
  country: string;
}

export type Address = {
  country_code: string;
  line1: string;
  line2: string;
  postal_code: string;
  city: string;
  state: string;
}


export type Contact = {
  name: string;
  email: string;
  subscribed_to_newsletter: boolean;
}

export type Website = {
  url: string;
  name: string;
  description: string;
  navigation: WebsiteNavigation;
  language: string;
  ad: string | null;
  announcement: string | null;
  logo: string | null;
  colors: {
    background: string;
    text: string;
    accent: string;
  },
  powered_by: boolean,
}

export type WebsiteNavigation = {
  primary: WebsiteNavigationItem[],
  secondary: WebsiteNavigationItem[],
}

export type WebsiteNavigationItem = {
  url?: string;
  label: string;
  children?: WebsiteNavigationItem[],
};

export type PageMetadata = {
  date: string;
  path: string;
  type: PageType;
  title: string;
  language: string;
  description: string;
  body_hash: string;
  metadata_hash: string;
};

export type Page = PageMetadata & {
  body: string;
  tags: Tag[];
}

// export type Block = {
//   type: BlockType;
//   data: any;
//   design: any;
// }

export type Tag = {
  name: string;
  description: string;
};

export type LoginOutput = {
  session_id: string;
};

export type SubscribeOutput = {
  contact_id: string;
};

export type PlaceOrderOutput = {
  stripe_checkout_url: string;
}

export type Order = {
  id: string;
  created_at: string;

  total_amount: number;
  currency: string;
  status: OrderStatus;
  invoice_url?: string;
}

export type Product = {
  id: string;

  type: ProductType;
  name: string;
  description: string;

  content: ProductPage[] | null;
}

export type ProductPage = {
  id: string;

  position: number;
  title: string;
  body: string;
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Input
////////////////////////////////////////////////////////////////////////////////////////////////////

export type GetPageInput = {
  slug?: string;
}

export type ListPagesInput = {
  tag?: string,
  type?: PageType,
}

export type TrackEventPageViewInput = {
  path: string;
  header_referrer: string;
  query_parameter_ref: string;
}

export type LoginInput = {
  email: string;
}

export type CompleteLoginInput = {
  session_id: string;
  code: string;
}

export type SubscribeInput = {
  email: string;
}

export type UnsubscribeInput = {
  email: string;
  token: string;
}

export type CompleteSubscriptionInput = {
  contact_id: string;
  code: string;
}

export type UpdateMyAccountInput = {
  name?: string;
  subscribed_to_newsletter?: boolean;
  email?: string;
}

export type VerifyEmailInput = {
  token: string;
}

export type VerifyEmailJwt = {
  email: string;
}

export type PlaceOrderInput = {
  products: string[];
  email?: string;
  subscribe_to_newsletter: boolean;
  additional_invoice_information?: string;
}

export type CompleteOrderInput = {
  order_id: string;
}

export type CancelOrderInput = {
  order_id: string;
}

export type GetProductInput = {
  id: string;
}
