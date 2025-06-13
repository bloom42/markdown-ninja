package service

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"markdown.ninja/pkg/errs"
	"markdown.ninja/pkg/services/store"
)

func validateAdditionalInvoiceInformation(additionalInformation string) error {
	const ADDITIONAL_INVOICE_INFORMATION_MAX_LENGTH = 800
	if len(additionalInformation) > ADDITIONAL_INVOICE_INFORMATION_MAX_LENGTH {
		return errs.InvalidArgument(fmt.Sprintf("Additional invoice information is too long. Max: %d characters", ADDITIONAL_INVOICE_INFORMATION_MAX_LENGTH))
	}

	return nil
}

func (service *StoreService) validateProductName(name string) error {
	if len(name) < store.ProductNameMinLength {
		return store.ErrProductNameIsTooShort
	}

	if len(name) > store.ProductNameMaxLength {
		return store.ErrProductNameIsTooLong
	}

	if !utf8.ValidString(name) {
		return store.ErrProductNameIsNotValid
	}

	return nil
}

func (service *StoreService) validateProductDescription(description string) error {
	if len(description) > store.ProductDescriptionMaxLength {
		return store.ErrProductDescriptionIsTooLong
	}

	if !utf8.ValidString(description) {
		return store.ErrProductDescriptionIsNotValid
	}

	return nil
}

func (service *StoreService) validateProductType(productType store.ProductType) error {
	if productType != store.ProductTypeBook && productType != store.ProductTypeCourse &&
		productType != store.ProductTypeDigitalDownload {
		return store.ErrProductTypeIsNotValid
	}

	return nil
}

func (service *StoreService) validateProductPrice(price int64) error {
	if price < 0 {
		return store.ErrProductPriceCantBeNegative
	}

	if price > store.ProductPriceMax {
		return store.ErrProductPriceIsTooHigh
	}

	return nil
}

func (service *StoreService) validateCouponCode(code string) error {
	if len(code) < store.CouponCodeMinLength {
		return store.ErrCouponCodeIsTooShort
	}

	if len(code) > store.CouponCodeMaxLength {
		return store.ErrCouponCodeIsTooLong
	}

	if strings.ToUpper(code) != code {
		return store.ErrCouponCodeMustBeUpperCase
	}

	if !utf8.ValidString(code) {
		return store.ErrCouponCodeIsNotValid
	}

	if !store.CouponCodeRegexp.MatchString(code) || strings.Contains(code, "--") {
		return store.ErrCouponCodeIsNotValid
	}

	return nil
}

func (service *StoreService) validateCouponDescription(description string) error {
	if len(description) > store.CouponDescriptionMaxLength {
		return store.ErrCouponDescriptionIsTooLong
	}

	if !utf8.ValidString(description) {
		return store.ErrCouponDescriptionIsNotValid
	}

	return nil
}

func (service *StoreService) validateCouponDiscount(discount int64) error {
	if discount < 1 || discount > 100 {
		return store.ErrCouponDiscountIsNotValid
	}

	return nil
}

func (service *StoreService) validateCouponExpiryDate(expiryDate time.Time) error {
	now := time.Now().UTC()

	if expiryDate.Before(now) {
		return errs.InvalidArgument("coupon's expiry date can't be in the past")
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Pages
////////////////////////////////////////////////////////////////////////////////////////////////////

func (service *StoreService) validateProductPageTitle(title string) error {
	err := service.contentService.ValidatePageTitle(title)
	if err != nil {
		return err
	}

	return nil
}

func (service *StoreService) validateProductPageContent(content string) error {
	err := service.contentService.ValidatePageBodyMarkdown(content)
	if err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Refunds
////////////////////////////////////////////////////////////////////////////////////////////////////

func (service *StoreService) validateRefundReason(reason store.RefundReason) (err error) {
	if !store.ValidRefundReasons.Contains(reason) {
		return store.ErrRefundReasonNotValid
	}

	return nil
}
