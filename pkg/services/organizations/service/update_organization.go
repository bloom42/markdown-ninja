package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/stripe/stripe-go/v81"
	stripecustomer "github.com/stripe/stripe-go/v81/customer"
	"markdown.ninja/pkg/errs"
	"markdown.ninja/pkg/server/httpctx"
	"markdown.ninja/pkg/services/kernel"
	"markdown.ninja/pkg/services/organizations"
)

func (service *OrganizationsService) UpdateOrganization(ctx context.Context, input organizations.UpdateOrganizationInput) (org organizations.Organization, err error) {
	httpCtx := httpctx.FromCtx(ctx)

	actorID, err := service.kernel.CurrentUserID(ctx)
	if err != nil {
		return
	}

	// if current user is not a Markdown Ninja admin then it needs to be an organization admin
	if !httpCtx.AccessToken.IsAdmin {
		var staff organizations.Staff
		staff, err = service.repo.FindStaff(ctx, service.db, actorID, input.ID)
		if err != nil {
			return
		}

		if staff.Role != organizations.StaffRoleAdministrator {
			err = kernel.ErrPermissionDenied
			return
		}

		// plan and extra slots can only be updated by admins
		if input.Plan != nil || input.ExtraSlots != nil {
			err = kernel.ErrPermissionDenied
			return
		}
	}

	tx, err := service.db.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("organizations.UpdateOrganization: Starting DB transaction: %w", err)
		return
	}
	defer tx.Rollback()

	org, err = service.repo.FindOrganizationByID(ctx, tx, input.ID, true)
	if err != nil {
		return
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		err = service.validateOrganizationName(name)
		if err != nil {
			return
		}
		org.Name = name
	}

	if input.BillingInformation != nil {
		if service.isSelfHosted {
			return org, errs.InvalidArgument("Billing information can't be updated in self-hosted mode")
		}

		billingInfo := *input.BillingInformation
		err = service.cleanAndValidateBillingInformation(ctx, httpCtx.AccessToken.IsAdmin, &billingInfo)
		if err != nil {
			return
		}
		org.BillingInformation = billingInfo
		if org.StripeCustomerID != nil {
			updateStripeCustomerParams := &stripe.CustomerParams{
				Name:  stripe.String(org.BillingInformation.Name),
				Email: stripe.String(org.BillingInformation.Email),
				Address: &stripe.AddressParams{
					Line1:      stripe.String(org.BillingInformation.AddressLine1),
					Line2:      stripe.String(org.BillingInformation.AddressLine2),
					City:       stripe.String(org.BillingInformation.City),
					PostalCode: stripe.String(org.BillingInformation.PostalCode),
					State:      stripe.String(org.BillingInformation.State),
					Country:    stripe.String(org.BillingInformation.CountryCode),
				},
			}
			_, err = stripecustomer.Update(*org.StripeCustomerID, updateStripeCustomerParams)
			if err != nil {
				err = fmt.Errorf("organizations.UpdateOrganization: error updating stripe customer for organization %s: %w", org.ID.String(), err)
				return
			}

			var stripeCustomerTaxIds []*stripe.TaxID
			stripeCustomerTaxIds, err = service.fetchStripeTaxIDsForCustomer(*org.StripeCustomerID)
			if err != nil {
				err = fmt.Errorf("organizations.UpdateOrganization: fetching taxIDs for organization [%s]: %w", org.ID, err)
				return
			}

			err = service.updateStripeTaxIDIfNeeded(ctx, org, stripeCustomerTaxIds)
			if err != nil {
				err = fmt.Errorf("organizations.UpdateOrganization: updating stripe tax ID for organization [%s]: %w", org.ID, err)
				return
			}
		}
	}

	if input.Plan != nil {
		err = validatePlan(*input.Plan, false)
		if err != nil {
			return
		}
		org.Plan = *input.Plan
	}

	if input.ExtraSlots != nil {
		err = validateExtraSlots(*input.ExtraSlots, org.Plan)
		if err != nil {
			return
		}
		org.ExtraSlots = *input.ExtraSlots
	}

	org.UpdatedAt = time.Now().UTC()
	err = service.repo.UpdateOrganization(ctx, tx, org)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("organizations.UpdateOrganization: Comitting DB transaction: %w", err)
		return
	}

	if httpCtx.AccessToken.IsAdmin {
		org.OrganizationAdminFields = &organizations.OrganizationAdminFields{
			StripeCustomerID:     org.StripeCustomerID,
			StripeSubscriptionID: org.StripeSubscriptionID,
		}
	}

	return
}
