package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bloom42/stdx-go/guid"
	"github.com/bloom42/stdx-go/opt"
	"github.com/bloom42/stdx-go/queue"
	"markdown.ninja/pkg/services/kernel"
	"markdown.ninja/pkg/services/organizations"
	"markdown.ninja/pkg/timeutil"
)

func (service *OrganizationsService) JobDispatchInvoiceMonthlyUsage(ctx context.Context, _ organizations.JobDispatchInvoiceMonthlyUsage) error {
	now := time.Now().UTC()

	// Only execute this job on the first monday of every months
	if now.Weekday() != time.Monday || now.Day() != timeutil.GetFirstMondayOfTheMonth(now).Day() {
		return nil
	}

	proOrganizations, err := service.repo.FindOrganizationsByPlan(ctx, service.db, kernel.PlanPro.ID)
	if err != nil {
		return err
	}

	jobs := make([]queue.NewJobInput, 0, len(proOrganizations))
	jobScheduledFor := time.Now().UTC()

	for i, organization := range proOrganizations {
		if organization.StripeCustomerID == nil || organization.StripeSubscriptionID == nil {
			continue
		}

		// Limit the number of requests per second
		if i != 0 && i%10 == 0 {
			jobScheduledFor = jobScheduledFor.Add(time.Second)
		}

		job := queue.NewJobInput{
			Data: organizations.JobInvoiceMonthlyUsage{
				OrganizationID: organization.ID,
				IdempotencyKey: guid.NewTimeBased().String(),
			},
			ScheduledFor: &jobScheduledFor,
			Timeout:      opt.Int64(300),
			RetryDelay:   opt.Int64(300),
		}
		jobs = append(jobs, job)
	}

	err = service.queue.PushMany(ctx, nil, jobs)
	if err != nil {
		return fmt.Errorf("organizations.JobDispatchInvoiceMonthlyUsage: pushing jobs to queue: %w", err)
	}

	return nil
}
