package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type MedicalRecord struct {
	IdentityNumber  string    `db:"identity_number" json:"identity_number"`
	Symptoms        string    `db:"symptoms" json:"symptoms"`
	Medications     string    `db:"medications" json:"medications"`
	CreatedByNip    string    `db:"created_by_nip" json:"createdByNip"`
	CreatedByName   string    `db:"created_by_name" json:"createdByName"`
	CreatedByUserId string    `db:"created_by_user_id" json:"createdByUserId"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
}

type RecordRegistrationPayload struct {
	IdentityNumber int64  `db:"identity_number" json:"identityNumber" form:"identityNumber"`
	Symptoms       string `db:"symptoms" json:"symptoms" form:"symptoms"`
	Medications    string `db:"medications" json:"medications" form:"medications"`
}

type GetRecordQueries struct {
	IdentityNumber  *int64 `db:"identity_number" json:"identityNumber" query:"identityNumber"`
	Limit           int    `json:"limit" query:"limit"`
	Offset          int    `json:"offset" query:"offset"`
	CreatedByNip    string `db:"created_by_nip" json:"createdByNip" query:"createdBy.nip"`
	CreatedByUserId string `db:"created_by_user_id" json:"createdByUserId" query:"createdBy.userId"`
	CreatedAt       string `db:"created_at" json:"createdAt" query:"createdAt"`
}

type GetRecordResponse struct {
	IdentityDetail PatientDetail   `json:"identityDetail"`
	Symptoms       string          `json:"symptoms"`
	Medications    string          `json:"medications"`
	CreatedBy      CreatedByDetail `json:"createdBy"`
	CreatedAt      string          `json:"createdAt"`
}

type CreatedByDetail struct {
	Nip    string `json:"nip"`
	Name   string `json:"name"`
	UserId string `json:"userId"`
}

type PatientDetail struct {
	IdentityNumber      int64  `json:"identityNumber"`
	PhoneNumber         string `json:"phoneNumber"`
	Name                string `json:"name"`
	BirthDate           string `json:"birthDate"`
	Gender              string `json:"gender"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}

func (u *RecordRegistrationPayload) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.IdentityNumber,
			validation.Required.Error("identityNumber is required"),
			validation.By(ValidateIdentityNumber),
		),
		validation.Field(&u.Symptoms,
			validation.Required.Error("symptoms is required"),
			validation.Length(1, 2000).Error("symptoms must be between 1 and 2000 characters"),
		),
		validation.Field(&u.Medications,
			validation.Required.Error("medications is required"),
			validation.Length(1, 2000).Error("medications must be between 1 and 2000 characters"),
		),
	)

	return err
}
