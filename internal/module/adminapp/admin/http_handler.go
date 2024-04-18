package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/tsel-ticketmaster/tm-user/internal/pkg/middleware"
	"github.com/tsel-ticketmaster/tm-user/pkg/errors"
	publicMiddleware "github.com/tsel-ticketmaster/tm-user/pkg/middleware"
	"github.com/tsel-ticketmaster/tm-user/pkg/response"
	"github.com/tsel-ticketmaster/tm-user/pkg/status"
)

type HTTPHandler struct {
	SessionMiddleware *middleware.AdminSession
	Validate          *validator.Validate
	AdminUseCase      AdminUseCase
}

func InitHTTPHandler(router *http.ServeMux, adminSession *middleware.AdminSession, validate *validator.Validate, adminUseCase AdminUseCase) {
	handler := &HTTPHandler{
		Validate:     validate,
		AdminUseCase: adminUseCase,
	}

	router.HandleFunc("POST /api/v1/adminapp/administrators/signin", publicMiddleware.SetRouteChain(handler.SignIn))
	router.HandleFunc("POST /api/v1/adminapp/administrators", publicMiddleware.SetRouteChain(handler.Create, adminSession.Verify))
	router.HandleFunc("POST /api/v1/adminapp/administrators/signout", publicMiddleware.SetRouteChain(handler.SignOut, adminSession.Verify))
}

func (handler HTTPHandler) validate(ctx context.Context, payload interface{}) error {
	err := handler.Validate.StructCtx(ctx, payload)
	if err == nil {
		return nil
	}

	errorFields := err.(validator.ValidationErrors)

	errMessages := make([]string, len(errorFields))

	for k, errorField := range errorFields {
		errMessages[k] = fmt.Sprintf("invalid '%s' with value '%v'", errorField.Field(), errorField.Value())
	}

	errorMessage := strings.Join(errMessages, ", ")

	return fmt.Errorf(errorMessage)

}

func (handler HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := CreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusUnprocessableEntity, response.RESTEnvelope{
			Status:  status.UNPROCESSABLE_ENTITY,
			Message: err.Error(),
		})

		return
	}

	if err := handler.validate(ctx, req); err != nil {
		response.JSON(w, http.StatusBadRequest, response.RESTEnvelope{
			Status:  status.BAD_REQUEST,
			Message: err.Error(),
		})

		return
	}

	resp, err := handler.AdminUseCase.Create(ctx, req)
	if err != nil {
		ae := errors.Destruct(err)
		response.JSON(w, ae.HTTPStatusCode, response.RESTEnvelope{
			Status:  ae.Status,
			Message: ae.Message,
		})

		return
	}

	response.JSON(w, http.StatusCreated, response.RESTEnvelope{
		Status:  status.CREATED,
		Message: "admin has been successfully created",
		Data:    resp,
	})
}

func (handler HTTPHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := SignInRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusUnprocessableEntity, response.RESTEnvelope{
			Status:  status.UNPROCESSABLE_ENTITY,
			Message: err.Error(),
		})

		return
	}

	if err := handler.validate(ctx, req); err != nil {
		response.JSON(w, http.StatusBadRequest, response.RESTEnvelope{
			Status:  status.BAD_REQUEST,
			Message: err.Error(),
		})

		return
	}

	resp, err := handler.AdminUseCase.SignIn(ctx, req)
	if err != nil {
		ae := errors.Destruct(err)
		response.JSON(w, ae.HTTPStatusCode, response.RESTEnvelope{
			Status:  ae.Status,
			Message: ae.Message,
		})

		return
	}

	response.JSON(w, http.StatusOK, response.RESTEnvelope{
		Status:  status.OK,
		Message: "admin has been successfully signed in",
		Data:    resp,
	})
}

func (handler HTTPHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := handler.AdminUseCase.SignOut(ctx)
	if err != nil {
		ae := errors.Destruct(err)
		response.JSON(w, ae.HTTPStatusCode, response.RESTEnvelope{
			Status:  ae.Status,
			Message: ae.Message,
		})

		return
	}

	response.JSON(w, http.StatusCreated, response.RESTEnvelope{
		Status:  status.OK,
		Message: "admin has been successfully signed out",
	})
}
