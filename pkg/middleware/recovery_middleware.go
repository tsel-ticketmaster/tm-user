package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"github.com/tsel-ticketmaster/tm-user/pkg/response"
	"github.com/tsel-ticketmaster/tm-user/pkg/status"
)

type Recovery struct {
	logger *logrus.Logger
	debug  bool
}

func NewRecovery(logger *logrus.Logger, debug bool) *Recovery {
	return &Recovery{
		logger: logger,
		debug:  debug,
	}
}

func (rm *Recovery) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {

				rm.logger.WithContext(r.Context()).WithFields(
					logrus.Fields{
						"panicking": "true",
					},
				).Error(recovered)

				if rm.debug {
					stack := debug.Stack()
					rm.logger.WithContext(r.Context()).Error(string(stack))
					w.Header().Set("X-Panic-Response", "true")
				}

				response.JSON(w, http.StatusInternalServerError, response.RESTEnvelope{
					Status:  status.INTERNAL_SERVER_ERROR,
					Message: "an error occured while trying to process the request",
					Data:    nil,
					Meta:    nil,
				})
			}

		}()

		next.ServeHTTP(w, r)
	})
}
