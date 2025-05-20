package service

import (
	"context"
	"strings"
	"time"

	"github.com/bloom42/stdx-go/guid"
	"github.com/bloom42/stdx-go/log/slogx"
	"github.com/bloom42/stdx-go/queue"
	"github.com/bloom42/stdx-go/slicesx"
	"markdown.ninja/pkg/services/kernel"
	"markdown.ninja/pkg/services/organizations"
)

func (service *OrganizationsService) InviteStaffs(ctx context.Context, input organizations.InviteStaffsInput) (invitations []organizations.StaffInvitation, err error) {
	actorID, err := service.kernel.CurrentUserID(ctx)
	if err != nil {
		return
	}

	logger := slogx.FromCtx(ctx)

	staff, err := service.CheckUserIsStaff(ctx, service.db, actorID, input.OrganizationID)
	if err != nil {
		return
	}

	if staff.Role != organizations.StaffRoleAdministrator {
		err = kernel.ErrPermissionDenied
		return
	}

	emails := slicesx.Unique(input.Emails)

	existingStaffs, err := service.getStaffsWithDetails(ctx, service.db, input.OrganizationID)
	if err != nil {
		return
	}
	existingStaffsByEmail := make(map[string]organizations.Staff, len(existingStaffs))
	for _, staff := range existingStaffs {
		existingStaffsByEmail[staff.Email] = staff.Staff
	}

	existingInvitations, err := service.repo.FindStaffInvitationsForOrganization(ctx, service.db, input.OrganizationID)
	if err != nil {
		return
	}
	existingInvitationsByEmail := make(map[string]organizations.StaffInvitation, len(existingInvitations))
	for _, invitation := range existingInvitations {
		existingInvitationsByEmail[invitation.InviteeEmail] = invitation
	}

	err = service.CheckBillingGatedAction(ctx, service.db, input.OrganizationID, organizations.BillingGatedActionInviteStaffs{NewStaffs: int64(len(emails))})
	if err != nil {
		return
	}

	now := time.Now().UTC()

	invitations = make([]organizations.StaffInvitation, 0, len(emails))
	invitationIDs := make([]guid.GUID, 0, len(emails))
	for _, email := range emails {
		// if staff / invitation already exist then we don't create / send new invitations
		_, staffExists := existingStaffsByEmail[email]
		_, invitationExists := existingInvitationsByEmail[email]

		if staffExists || invitationExists {
			continue
		}

		email = strings.TrimSpace(email)
		if email == "" {
			continue
		}

		err = service.kernel.ValidateEmail(ctx, email, true)
		if err != nil {
			return
		}

		invitationID := guid.NewTimeBased()
		invitationIDs = append(invitationIDs, invitationID)
		invitation := organizations.StaffInvitation{
			ID:             invitationID,
			CreatedAt:      now,
			UpdatedAt:      now,
			Role:           organizations.StaffRoleAdministrator,
			InviteeEmail:   email,
			OrganizationID: input.OrganizationID,
			InviterID:      actorID,
		}
		invitations = append(invitations, invitation)
	}

	err = service.repo.CreateStaffInvitations(ctx, service.db, invitations)
	if err != nil {
		return
	}

	err = service.queue.Push(ctx, nil, queue.NewJobInput{
		Data: organizations.JobSendStaffInvitations{
			InvitationIDs: invitationIDs,
		},
	})
	if err != nil {
		logger.Error("organizations.InviteStaffs: pushing job to queue", slogx.Err(err))
		err = nil
	}

	return
}
