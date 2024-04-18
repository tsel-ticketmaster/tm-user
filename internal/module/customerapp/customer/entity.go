package customer

import "time"

const (
	changeEmailVerificationKeyPrefix = "user:change_email_verification:customer:token:%s"
	VerificationURLPath              = "/v1/customerapp/customers/change-email-verification"

	VerficationStatusVerified    = "VERIFIED"
	VerificationStatusUnverified = "UNVERIFIED"

	MemberStatusActive   = "ACTIVE"
	MemberStatusInactive = "INACTIVE"
)

type Customer struct {
	ID                 int64
	Name               string
	Email              string
	Password           string
	PasswordSalt       string
	VerificationStatus string
	MemberStatus       string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
