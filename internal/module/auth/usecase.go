package auth

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/tsel-ticketmaster/tm-user/internal/pkg/session"
// 	"github.com/sirupsen/logrus"
// )

// type Usecase interface {
// 	Register(context.Context, *RegisterRequest) error
// }

// type usecase struct {
// 	logger  *logrus.Logger
// 	session session.Session
// }

// type UsecaseProperty struct {
// 	Logger  *logrus.Logger
// 	Session session.Session
// }

// func NewUsecase(props *UsecaseProperty) Usecase {
// 	return &usecase{
// 		logger:  props.Logger,
// 		session: props.Session,
// 	}
// }

// Register implements Usecase.
// func (u *usecase) Register(ctx context.Context, req *RegisterRequest) error {
// 	u.logger.WithContext(ctx).Info("ok")
// 	u.logger.WithContext(ctx).WithError(fmt.Errorf("error abcd")).Error()

// 	sessAcc := session.Account{
// 		ID:   "26af1b01-d9a0-4508-a4ae-500927d6b308",
// 		Name: "patrick",
// 		Type: "user",
// 	}
// 	if err := u.session.Set(ctx, sessAcc.ID, sessAcc, time.Hour*1); err != nil {
// 		return err
// 	}

// 	return nil
// }
