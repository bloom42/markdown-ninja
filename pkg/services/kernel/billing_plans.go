package kernel

const (
	// number of bytes of assets storage per slot
	StoragePerSlot = 5_000_000_000 // 5 GB
)

type PlanID string

type Plan struct {
	ID          PlanID   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int64    `json:"price"`
	Features    []string `json:"features"`

	// Allowed storage in bytes
	MaxExtraSlots           int64 `json:"-"`
	AllowedStorage          int64 `json:"-"`
	AllowedEmails           int64 `json:"-"`
	AllowedPages            int64 `json:"-"`
	AllowedAssets           int64 `json:"-"`
	SelfServe               bool  `json:"-"`
	MaxAssetSize            int64 `json:"-"`
	CustomDomainsPerWebsite int64 `json:"-"`
}

var PlanFree = Plan{
	ID:          "free",
	Name:        "Free",
	Description: "The perfect place to start your blogging journey or replace a static site",
	Price:       0,
	SelfServe:   true,

	AllowedStorage:          10_000_000, // 10 MB
	MaxExtraSlots:           0,
	AllowedEmails:           0,
	AllowedPages:            25,
	AllowedAssets:           50,
	MaxAssetSize:            1_000_000, // 1 MB
	CustomDomainsPerWebsite: 0,

	Features: []string{
		"1 Website",
		"1 Staff",
		// "10 MB assets storage",
		"Built-in Analytics",
		"Headless CMS",
	},
}

var PlanPro = Plan{
	ID:   "pro",
	Name: "Pro",
	// Description: "Everything you need to scale your audience",
	Description: "For bloggers, startups and agencies looking to grow their audience",
	Price:       5,
	SelfServe:   true,

	AllowedStorage:          StoragePerSlot, // base: 5GB
	MaxExtraSlots:           1000,
	AllowedEmails:           500_000,
	AllowedPages:            3000,
	AllowedAssets:           3000,
	MaxAssetSize:            MaxAssetSize,
	CustomDomainsPerWebsite: 5,

	Features: []string{
		// 'No additional transaction fees',
		// 'Website with custom domain',
		"1 Website / slot",
		"1 Staff / slot",
		"5 GB of assets storage / slot",
		"1€ / 1000 emails",
		"Custom Domains",
		"Unlimited page views",
		// 'Unlimited contacts',
		"Advanced Analytics",
		// 'Priority support'
	},
}

var PlanEnterprise = Plan{
	ID:          "enterprise",
	Name:        "Enterprise",
	Description: "Control, compliance, and support tailored for large scale organizations",
	Price:       2000,
	SelfServe:   false,

	AllowedStorage:          StoragePerSlot * 10,
	MaxExtraSlots:           10000,
	AllowedEmails:           10_000_000,
	AllowedPages:            200_000,
	AllowedAssets:           200_000,
	MaxAssetSize:            MaxAssetSize,
	CustomDomainsPerWebsite: 50,

	Features: []string{
		"Unlimited Staffs",
		"Premium support",
		"Self-hosting support",
		"Commercial license",
		"And a lot of other unfair advantages",
		// 'No additional transaction fees',
		// 'Website with custom domain',
		// 'Unlimited page views',
		// 'Unlimited contacts',
		// '10 Websites',
		// '10€ / extra Website',
		// // 'Unlimited products',
		// 'Advanced Analytics',
		// // 'Custom themes',
		// '10 Staffs / 10€ extra staff',
		// '50 GB assets storage',
		// '1€ / extra GB',
		// '80,000 Emails',
		// '1€ / extra 1000 emails',
		// 'Advanced Analytics',
		// 'Priority support'
	},
}

var AllPlans = map[PlanID]Plan{
	PlanFree.ID:       PlanFree,
	PlanPro.ID:        PlanPro,
	PlanEnterprise.ID: PlanEnterprise,
}
