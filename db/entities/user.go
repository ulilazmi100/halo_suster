package entities

import (
	"errors"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	Id                  string    `db:"id" json:"id"`
	Nip                 string    `db:"nip" json:"nip"`
	Name                string    `db:"name" json:"name"`
	Role                string    `db:"role" json:"role"`
	Password            string    `db:"password" json:"password,omitempty"`
	IdentityCardScanImg string    `db:"identity_card_scan_img" json:"identityCardScanImg"`
	Access              bool      `db:"access" json:"access"`
	CreatedAt           time.Time `db:"created_at" json:"createdAt"`
}

type RegistrationPayload struct {
	Nip      string `json:"nip" form:"nip"`
	Name     string `json:"name,omitempty" form:"name"`
	Password string `json:"password" form:"password"`
}

type NurseRegistrationPayload struct {
	Nip                 string `json:"nip" form:"nip"`
	Name                string `json:"name,omitempty" form:"name"`
	IdentityCardScanImg string `json:"identityCardScanImg" form:"identityCardScanImg"`
}

type NurseUpdatePayload struct {
	Nip  string `json:"nip" form:"nip"`
	Name string `json:"name,omitempty" form:"name"`
}

type NurseAccessPayload struct {
	Password string `json:"password" form:"password"`
}

type GetUserQueries struct {
	UserId    string `db:"id" json:"userId" query:"userId"`
	Limit     int    `json:"limit" query:"limit"`
	Offset    int    `json:"offset" query:"offset"`
	Name      string `db:"name" json:"name" query:"name"`
	Nip       string `db:"nip" json:"nip" query:"nip"`
	Role      string `db:"role" json:"role" query:"role"`
	CreatedAt string `db:"created_at" json:"createdAt" query:"createdAt"`
}

type GetUserResponse struct {
	UserId    string `db:"id" json:"userId" query:"userId"`
	Nip       string `db:"nip" json:"nip" query:"nip"`
	Name      string `db:"name" json:"name" query:"name"`
	CreatedAt string `db:"created_at" json:"createdAt" query:"createdAt"`
}

type Credential struct {
	Nip      string `json:"nip"`
	Password string `json:"password"`
}

type JWTPayload struct {
	Id   string
	Nip  string
	Name string
}

type JWTClaims struct {
	Id   string
	Nip  string
	Name string
	jwt.RegisteredClaims
}

func NewUser(nip, name, password string) *User {
	u := &User{
		Nip:      nip,
		Name:     name,
		Password: password,
	}

	return u
}

func (u *Credential) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Nip,
			validation.Required.Error("nip is required"),
		),
		validation.Field(&u.Password,
			validation.Required.Error("password is required"),
			validation.Length(5, 33).Error("password must be between 5 and 33 characters"),
		),
	)

	return err
}

func (u *RegistrationPayload) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Nip,
			validation.Required.Error("nip is required"),
			validation.Length(13, 13).Error("nip number must be 13 characters"),
			validation.By(ValidateITStaffNipFormat),
		),
		validation.Field(&u.Name,
			validation.Required.Error("name is required"),
			validation.Length(5, 50).Error("name must be between 5 and 50 characters"),
		),
		validation.Field(&u.Password,
			validation.Required.Error("password is required"),
			validation.Length(5, 33).Error("password must be between 5 and 33 characters"),
		),
	)

	return err
}

func (u *NurseRegistrationPayload) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Nip,
			validation.Required.Error("nip is required"),
			validation.Length(13, 13).Error("nip number must be 13 characters"),
			validation.By(ValidateNurseNipFormat),
		),
		validation.Field(&u.Name,
			validation.Required.Error("name is required"),
			validation.Length(5, 50).Error("name must be between 5 and 50 characters"),
		),
		validation.Field(&u.IdentityCardScanImg,
			validation.Required.Error("identityCardScanImg is required"),
			validation.By(ValidateImageURL),
		),
	)

	return err
}

func (u *NurseUpdatePayload) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Nip,
			validation.Required.Error("nip is required"),
			validation.Length(13, 13).Error("nip number must be 13 characters"),
		),
		validation.Field(&u.Name,
			validation.Required.Error("name is required"),
			validation.Length(5, 50).Error("name must be between 5 and 50 characters"),
		),
	)

	return err
}

func (u *NurseAccessPayload) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Password,
			validation.Required.Error("password is required"),
			validation.Length(5, 33).Error("password must be between 5 and 33 characters"),
		),
	)

	return err
}

func ValidateNipFormat(value any) error {
	nip, ok := value.(string)
	if !ok {
		return errors.New("parse error")
	}

	pattern := `^(303|615)[12](200[0-9]|201[0-9]|202[0-4])(0[1-9]|1[0-2])([0-9]{3})$`
	rgx := regexp.MustCompile(pattern)
	if !rgx.MatchString(nip) {
		return errors.New("invalid nip format")
	}

	return nil
}

func ValidateNurseNipFormat(value any) error {
	nip, ok := value.(string)
	if !ok {
		return errors.New("parse error")
	}

	pattern := `^303[12](200[0-9]|201[0-9]|202[0-4])(0[1-9]|1[0-2])([0-9]{3})$`
	rgx := regexp.MustCompile(pattern)
	if !rgx.MatchString(nip) {
		return errors.New("invalid nurse nip format")
	}

	return nil
}

func ValidateITStaffNipFormat(value any) error {
	nip, ok := value.(string)
	if !ok {
		return errors.New("parse error")
	}

	pattern := `^615[12](200[0-9]|201[0-9]|202[0-4])(0[1-9]|1[0-2])([0-9]{3})$`
	rgx := regexp.MustCompile(pattern)
	if !rgx.MatchString(nip) {
		return errors.New("invalid IT staff nip format")
	}

	return nil
}

func ValidateImageURL(value any) error {
	url, ok := value.(string)
	if !ok {
		return errors.New("parse error")
	}
	pattern := `^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+(?:.jpg|.jpeg)+$`
	rgx := regexp.MustCompile(pattern)
	if !rgx.MatchString(url) {
		return errors.New("invalid image URL format")
	}

	return nil
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
