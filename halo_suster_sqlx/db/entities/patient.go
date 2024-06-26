package entities

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Patient struct {
	IdentityNumber      string    `db:"identity_number" json:"id"`
	PhoneNumber         string    `db:"phone_number" json:"phoneNumber"`
	Name                string    `db:"name" json:"name"`
	BirthDate           string    `db:"birth_date" json:"birthDate"`
	Gender              string    `db:"gender" json:"gender"`
	IdentityCardScanImg string    `db:"identity_card_scan_img" json:"identityCardScanImg"`
	CreatedAt           time.Time `db:"created_at" json:"createdAt"`
}

type PatientRegistrationPayload struct {
	IdentityNumber      int64  `db:"identity_number" json:"identityNumber" form:"identityNumber"`
	PhoneNumber         string `db:"phone_number" json:"phoneNumber" form:"phoneNumber"`
	Name                string `db:"name" json:"name" form:"name"`
	BirthDate           string `db:"birth_date" json:"birthDate" form:"birthDate"`
	Gender              string `db:"gender" json:"gender" form:"gender"`
	IdentityCardScanImg string `db:"identity_card_scan_img" json:"identityCardScanImg" form:"identityCardScanImg"`
}

type GetPatientQueries struct {
	IdentityNumber *int64 `db:"identity_number" json:"identityNumber" query:"identityNumber"`
	Limit          int    `json:"limit" query:"limit"`
	Offset         int    `json:"offset" query:"offset"`
	Name           string `db:"name" json:"name" query:"name"`
	PhoneNumber    string `db:"phone_number" json:"phoneNumber" query:"phoneNumber"`
	CreatedAt      string `db:"created_at" json:"createdAt" query:"createdAt"`
}

type GetPatientResponse struct {
	IdentityNumber int64  `db:"identity_number" json:"identityNumber"`
	PhoneNumber    string `db:"phone_number" json:"phoneNumber"`
	Name           string `db:"name" json:"name"`
	BirthDate      string `db:"birth_date" json:"birthDate"`
	Gender         string `db:"gender" json:"gender"`
	CreatedAt      string `db:"created_at" json:"createdAt"`
}

func (u *PatientRegistrationPayload) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.IdentityNumber,
			validation.Required.Error("identityNumber is required"),
			validation.By(ValidateIdentityNumber),
		),
		validation.Field(&u.PhoneNumber,
			validation.Required.Error("phoneNumber is required"),
			validation.Length(10, 15).Error("phoneNumber must be between 10 and 15 characters"),
			validation.By(ValidatePhonePrefix),
		),
		validation.Field(&u.Name,
			validation.Required.Error("name is required"),
			validation.Length(3, 30).Error("name must be between 3 and 30 characters"),
		),
		validation.Field(&u.BirthDate,
			validation.Required.Error("birthDate is required"),
			validation.By(ValidateBirthDate),
		),
		validation.Field(&u.Gender,
			validation.Required.Error("gender is required"),
			validation.In("male", "female"),
		),
		validation.Field(&u.IdentityCardScanImg,
			validation.Required.Error("identityCardScanImg is required"),
			validation.By(ValidateImageURL),
		),
	)

	return err
}

func ValidatePhonePrefix(value any) error {
	phoneNumber, ok := value.(string)
	if !ok {
		return errors.New("parse error")
	}

	// Regex pattern to check if string starts with "+62"
	pattern := `^\+62`
	rgx := regexp.MustCompile(pattern)
	if !rgx.MatchString(phoneNumber) {
		return errors.New("invalid phoneNumber format: must start with +62")
	}

	return nil
}

func ValidateIdentityNumber(value any) error {
	identityNumber, ok := value.(int64)
	if !ok {
		return errors.New("parse error")
	}

	identityNumberString := fmt.Sprintf("%d", identityNumber)

	pattern := `^\d{16}$`
	rgx := regexp.MustCompile(pattern)
	if !rgx.MatchString(identityNumberString) {
		return errors.New("invalid identity number format")
	}

	return nil
}

func ValidateBirthDate(value any) error {
	birthDate, ok := value.(string)
	if !ok {
		return errors.New("parse error")
	}

	// Regex pattern to check if string is in ISO8601 format
	pattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})?$`
	rgx := regexp.MustCompile(pattern)
	if !rgx.MatchString(birthDate) {
		return errors.New("invalid birthDate format: must follow ISO8601")
	}

	return nil
}
