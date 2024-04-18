package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/tsel-ticketmaster/tm-user/internal/pkg/jwt"
	"github.com/tsel-ticketmaster/tm-user/internal/pkg/session"
	"github.com/tsel-ticketmaster/tm-user/internal/pkg/util"
	"github.com/tsel-ticketmaster/tm-user/pkg/errors"
	"github.com/tsel-ticketmaster/tm-user/pkg/status"
)

type CustomerUseCase interface {
	SignUp(ctx context.Context, req SignUpRequest) (SignUpResponse, error)
	SignIn(ctx context.Context, req SignInRequest) (SignInResponse, error)
	SignOut(ctx context.Context) error
	GetProfile(ctx context.Context) (GetProfileResponse, error)
	UpdateProfile(ctx context.Context, req UpdateProfileRequest) error
	ChangeEmail(ctx context.Context, req ChangeEmailRequest) (ChangeEmailResponse, error)
	ChangePassword(ctx context.Context, req ChangePasswordRequest) error
}

type CustomerUseCaseProperty struct {
	Logger             *logrus.Logger
	Timeout            time.Duration
	UserBaseURL        string
	JSONWebToken       *jwt.JSONWebToken
	Session            session.Session
	Cache              redis.UniversalClient
	CustomerRepository CustomerRepository
}

type customerUseCase struct {
	logger             *logrus.Logger
	timeout            time.Duration
	userBaseURL        string
	jsonWebToken       *jwt.JSONWebToken
	session            session.Session
	cache              redis.UniversalClient
	customerRepository CustomerRepository
}

// ChangeEmail implements CustomerUseCase.
func (u *customerUseCase) ChangeEmail(ctx context.Context, req ChangeEmailRequest) (ChangeEmailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	acc, err := session.GetAccountFromCtx(ctx)
	if err != nil {
		return ChangeEmailResponse{}, err
	}

	c, err := u.customerRepository.FindByID(ctx, acc.ID, nil)
	if err != nil {
		return ChangeEmailResponse{}, err
	}

	if req.Email == c.Email {
		return ChangeEmailResponse{}, errors.New(http.StatusBadRequest, status.BAD_REQUEST, "the new email is the same as existing email")
	}

	now := time.Now()
	linkExpiresIn := time.Minute * 5
	linkExpiresAt := now.Add(linkExpiresIn)
	verificationToken := util.GenerateRandomHEX(32)
	verificationKey := fmt.Sprintf(changeEmailVerificationKeyPrefix, verificationToken)
	verificationLink := fmt.Sprintf("%s%s?token=%s", u.userBaseURL, VerificationURLPath, verificationToken)
	changeEmailEvent := ChangeEmailEvent{
		ID:              c.ID,
		Name:            c.Name,
		ExistingEmail:   c.Email,
		NewEmail:        req.Email,
		VerficationLink: verificationLink,
	}

	changeEmailEventBuff, _ := json.Marshal(changeEmailEvent)

	if err := u.cache.Set(ctx, verificationKey, changeEmailEventBuff, linkExpiresIn).Err(); err != nil {
		u.logger.WithContext(ctx).WithError(err).Error()
		return ChangeEmailResponse{}, errors.New(http.StatusInternalServerError, status.INTERNAL_SERVER_ERROR, "an error occured while changing customer's email")
	}

	// publish for email notification

	if err := u.session.Delete(ctx, fmt.Sprintf("customer:%d", c.ID)); err != nil {
		return ChangeEmailResponse{}, err
	}

	resp := ChangeEmailResponse{
		VerificationExpiresAt: linkExpiresAt,
	}

	return resp, nil
}

// ChangePassword implements CustomerUseCase.
func (u *customerUseCase) ChangePassword(ctx context.Context, req ChangePasswordRequest) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	acc, err := session.GetAccountFromCtx(ctx)
	if err != nil {
		return err
	}

	c, err := u.customerRepository.FindByID(ctx, acc.ID, nil)
	if err != nil {
		return err
	}

	hashedExistingPassword := util.GenerateSecret(req.ExistingPassword, c.PasswordSalt)
	if hashedExistingPassword != c.Password {
		return errors.New(http.StatusBadRequest, status.BAD_REQUEST, "invalid customer's existing password")
	}

	newPasswordSalt := util.GenerateRandomHEX(16)
	newHashedPassword := util.GenerateSecret(req.NewPassword, newPasswordSalt)

	c.Password = newHashedPassword
	c.PasswordSalt = newPasswordSalt

	if err := u.customerRepository.Update(ctx, c.ID, c, nil); err != nil {
		return err
	}

	if err := u.session.Delete(ctx, fmt.Sprintf("customer:%d", c.ID)); err != nil {
		return err
	}

	return nil
}

// GetProfile implements CustomerUseCase.
func (u *customerUseCase) GetProfile(ctx context.Context) (GetProfileResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	acc, err := session.GetAccountFromCtx(ctx)
	if err != nil {
		return GetProfileResponse{}, err
	}

	c, err := u.customerRepository.FindByID(ctx, acc.ID, nil)
	if err != nil {
		return GetProfileResponse{}, err
	}

	resp := GetProfileResponse{
		ID:                 c.ID,
		Name:               c.Name,
		Email:              c.Email,
		VerificationStatus: c.VerificationStatus,
		MemberStatus:       c.MemberStatus,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}

	return resp, nil
}

// SignIn implements CustomerUseCase.
func (c *customerUseCase) SignIn(ctx context.Context, req SignInRequest) (SignInResponse, error) {
	panic("unimplemented")
}

// SignOut implements CustomerUseCase.
func (c *customerUseCase) SignOut(ctx context.Context) error {
	panic("unimplemented")
}

// SignUp implements CustomerUseCase.
func (c *customerUseCase) SignUp(ctx context.Context, req SignUpRequest) (SignUpResponse, error) {
	panic("unimplemented")
}

// UpdateProfile implements CustomerUseCase.
func (c *customerUseCase) UpdateProfile(ctx context.Context, req UpdateProfileRequest) error {
	panic("unimplemented")
}

func NewCustomerUseCase(props CustomerUseCaseProperty) CustomerUseCase {
	return &customerUseCase{
		logger:             props.Logger,
		timeout:            props.Timeout,
		jsonWebToken:       props.JSONWebToken,
		session:            props.Session,
		customerRepository: props.CustomerRepository,
	}
}
