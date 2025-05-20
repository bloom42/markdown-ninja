package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bloom42/stdx-go/guid"
	"github.com/bloom42/stdx-go/log/slogx"
	"github.com/stripe/stripe-go/v81"
	stripecustomer "github.com/stripe/stripe-go/v81/customer"
	"markdown.ninja/pkg/server/httpctx"
	"markdown.ninja/pkg/services/kernel"
	"markdown.ninja/pkg/services/organizations"
)

func (service *OrganizationsService) SyncStripe(ctx context.Context, input organizations.SyncStripeInput) (err error) {
	httpCtx := httpctx.FromCtx(ctx)

	actorID, err := service.kernel.CurrentUserID(ctx)
	if err != nil {
		return
	}

	// if current user is not a Markdown Ninja admin then it needs to be an organization admin
	if !httpCtx.AccessToken.IsAdmin {
		var staff organizations.Staff
		staff, err = service.repo.FindStaff(ctx, service.db, actorID, input.OrganizationID)
		if err != nil {
			return
		}

		if staff.Role != organizations.StaffRoleAdministrator {
			err = kernel.ErrPermissionDenied
			return
		}
	}

	_, err = service.syncOrganizationWithStripeCustomer(ctx, input.OrganizationID)
	return
}

func (service *OrganizationsService) syncOrganizationWithStripeCustomer(ctx context.Context, organizationID guid.GUID) (organization organizations.Organization, err error) {
	logger := slogx.FromCtx(ctx).With(slog.String("organization.id", organizationID.String()))
	now := time.Now().UTC()

	tx, err := service.db.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("organizations.syncOrganizationWithStripeCustomer: Starting DB transaction: %w", err)
		return
	}
	defer tx.Rollback()

	organization, err = service.repo.FindOrganizationByID(ctx, tx, organizationID, true)
	if err != nil {
		return
	}

	// Do nothing. Should we send / log an error?
	if organization.StripeCustomerID == nil {
		return
	}

	logger = logger.With(slog.String("stripe.customer.id", *organization.StripeCustomerID))

	getCustomerParams := &stripe.CustomerParams{}
	getCustomerParams.AddExpand("tax_ids")
	getCustomerParams.AddExpand("subscriptions")
	stripeCustomer, err := stripecustomer.Get(*organization.StripeCustomerID, getCustomerParams)
	if err != nil {
		err = fmt.Errorf("organizations.syncOrganizationWithStripeCustomer: fetching stripe customer: %w", err)
		return
	}

	// Update organization from customer

	organization.BillingInformation.Name = stripeCustomer.Name
	organization.BillingInformation.Email = stripeCustomer.Email
	organization.BillingInformation.AddressLine1 = stripeCustomer.Address.Line1
	organization.BillingInformation.AddressLine2 = stripeCustomer.Address.Line2
	organization.BillingInformation.PostalCode = stripeCustomer.Address.PostalCode
	organization.BillingInformation.City = stripeCustomer.Address.City
	organization.BillingInformation.State = stripeCustomer.Address.State
	organization.BillingInformation.CountryCode = stripeCustomer.Address.Country

	if stripeCustomer.TaxIDs != nil && len(stripeCustomer.TaxIDs.Data) != 0 {
		for _, taxID := range stripeCustomer.TaxIDs.Data {
			if taxID.Type == stripe.TaxIDTypeEUVAT {
				organization.BillingInformation.TaxID = &taxID.Value
				break
			}
		}
	}

	// Update subscription if needed
	if stripeCustomer.Subscriptions != nil && len(stripeCustomer.Subscriptions.Data) != 0 {

		nonCanceledSubscriptionsCount := 0
		for _, stripeSubscription := range stripeCustomer.Subscriptions.Data {
			// ignore canceled subscriptions
			if stripeSubscription.Status == stripe.SubscriptionStatusCanceled {
				continue
			}
			nonCanceledSubscriptionsCount += 1

			organization.StripeSubscriptionID = &stripeSubscription.ID

			planStr, planOk := stripeSubscription.Metadata["markdown_ninja_plan"]
			if planOk {
				organization.Plan = kernel.PlanID(planStr)
			} else {
				logger.Error("organizations.syncOrganizationWithStripeCustomer: markdown_ninja_plan metadata is missing for subscription")
			}

			stripeSubscriptionStartDate := time.Unix(stripeSubscription.StartDate, 0).UTC()
			if organization.SubscriptionStartedAt == nil ||
				!organization.SubscriptionStartedAt.Equal(stripeSubscriptionStartDate) {
				organization.SubscriptionStartedAt = &stripeSubscriptionStartDate
			}

			if stripeSubscription.Status == stripe.SubscriptionStatusActive {
				organization.PaymentDueSince = nil
			} else {
				organization.PaymentDueSince = &now
			}
		}

		// if no subscription
		if nonCanceledSubscriptionsCount == 0 {
			organization.StripeSubscriptionID = nil
			organization.SubscriptionStartedAt = nil
			organization.Plan = kernel.PlanFree.ID
			organization.ExtraSlots = 0
		} else if nonCanceledSubscriptionsCount > 1 {
			logger.Error("organizations.syncOrganizationWithStripeCustomer: Organization has more than 1 non-canceled subscription")
		}
	}

	organization.UpdatedAt = now
	err = service.repo.UpdateOrganization(ctx, tx, organization)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("organizations.syncOrganizationWithStripeCustomer: Comitting DB transaction: %w", err)
		return
	}

	return
}
