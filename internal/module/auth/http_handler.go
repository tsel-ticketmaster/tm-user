package auth

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/go-playground/validator/v10"
// 	"github.com/tsel-ticketmaster/tm-user/pkg/errors"
// 	"github.com/tsel-ticketmaster/tm-user/pkg/response"
// 	"github.com/tsel-ticketmaster/tm-user/pkg/status"
// )

// type HTTPHandler struct {
// 	Router   *http.ServeMux
// 	Validate *validator.Validate
// 	Usecase  Usecase
// }

// func (handler HTTPHandler) Handle() {
// 	handler.Router.HandleFunc("POST /api/v1/auth/register", handler.Register)

// }

// func (handler HTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
// 	var req RegisterRequest

// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		response.JSON(w, http.StatusUnprocessableEntity, response.RESTEnvelope{
// 			Status:  status.UNPROCESSABLE_ENTITY,
// 			Message: err.Error(),
// 		})
// 		return
// 	}

// 	err = handler.Usecase.Register(r.Context(), nil)
// 	if err != nil {
// 		ae := errors.Destruct(err)
// 		response.JSON(w, ae.HTTPStatusCode, response.RESTEnvelope{
// 			Status:  ae.Status,
// 			Message: ae.Message,
// 		})
// 		return
// 	}

// 	response.JSON(w, http.StatusCreated, response.RESTEnvelope{
// 		Status:  status.OK,
// 		Message: "account is already registered",
// 	})
// }
