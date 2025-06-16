package service

import (
	"fmt"
	"unicode/utf8"

	"markdown.ninja/pkg/errs"
	"markdown.ninja/pkg/services/contacts"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// Contacts
////////////////////////////////////////////////////////////////////////////////////////////////////

func (service *ContactsService) ValidateContactName(name string) error {
	if name == "" {
		return nil
	}

	if len(name) > contacts.ContactNameMaxLength {
		return errs.InvalidArgument(fmt.Sprintf("Contact name is too long. max: %d characters", contacts.ContactNameMaxLength))
	}

	if !utf8.ValidString(name) {
		return contacts.ErrContactNameIsNotValid
	}

	return nil
}
