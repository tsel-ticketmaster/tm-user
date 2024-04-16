package session

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/tsel-ticketmaster/tm-user/pkg/errors"
	"github.com/tsel-ticketmaster/tm-user/pkg/status"
)

type AccountContextKey struct{}

type Account struct {
	ID   int64
	Name string
	Type string
}

type Session interface {
	Set(ctx context.Context, key string, acc Account, ttl time.Duration) error
	Delete(ctx context.Context) error
	Get(ctx context.Context, key string) (Account, error)
}

type redisSessionStore struct {
	l *logrus.Logger
	r redis.UniversalClient
}

// Delete implements Session.
func (s *redisSessionStore) Delete(ctx context.Context) error {
	panic("unimplemented")
}

// Get implements Session.
func (s *redisSessionStore) Get(ctx context.Context, key string) (Account, error) {
	sessionKey := fmt.Sprintf("user:session:%s", key)
	acc := Account{}
	dataBuff, err := s.r.Get(ctx, sessionKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return acc, errors.New(http.StatusNotFound, status.NOT_FOUND, "user session is not found")
		}

		s.l.WithContext(ctx).WithError(err).Error()
		return acc, errors.New(http.StatusInternalServerError, status.INTERNAL_SERVER_ERROR, "")
	}

	json.Unmarshal(dataBuff, &acc)

	return acc, nil
}

// Set implements Session.
func (s *redisSessionStore) Set(ctx context.Context, key string, acc Account, ttl time.Duration) error {
	sessionKey := fmt.Sprintf("user:session:%s", key)

	transaction := func(tx *redis.Tx) error {
		err := tx.Get(ctx, sessionKey).Err()

		if err == nil {
			return errors.New(http.StatusForbidden, status.ALREADY_SIGNED_IN, "")
		}
		if err != redis.Nil {
			s.l.WithContext(ctx).WithError(err).Error()
			return errors.New(http.StatusInternalServerError, status.INTERNAL_SERVER_ERROR, "")
		}

		_, err = tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
			accBuff, _ := json.Marshal(acc)
			if err := p.Set(ctx, sessionKey, accBuff, ttl).Err(); err != nil {
				s.l.WithContext(ctx).WithError(err).Error()
				return errors.New(http.StatusInternalServerError, status.INTERNAL_SERVER_ERROR, "")
			}

			return nil
		})

		return err
	}

	err := s.r.Watch(ctx, transaction, sessionKey)
	if err != nil {
		if _, ok := err.(*errors.AppError); !ok {
			s.l.WithContext(ctx).WithError(err).Error()
		}
	}

	return err
}

func NewRedisSessionStore(l *logrus.Logger, r redis.UniversalClient) Session {
	return &redisSessionStore{
		l: l,
		r: r,
	}
}

func GetAccountFromCtx(ctx context.Context) (Account, error) {
	if ctx == nil {
		return Account{}, errors.New(http.StatusForbidden, status.FORBIDDEN, "request has an empty context")
	}
	value := ctx.Value(AccountContextKey{})

	acc, ok := value.(Account)
	if !ok {
		return Account{}, errors.New(http.StatusForbidden, status.FORBIDDEN, "request has an invalid context")
	}

	return acc, nil
}
