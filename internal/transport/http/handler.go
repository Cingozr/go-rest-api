package http

import (
	"encoding/json"
	"net/http"

	"github.com/Cingozr/go-rest-api/internal/comment"
	"github.com/gorilla/mux"
)

type Handler struct {
	Router  *mux.Router
	Service *comment.Service
}

type ResponseModel struct {
	Message string
	Error   string
}

func NewHandler(service *comment.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

//sendOkResponse - Retrive error message handler
func sendOkResponse(w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

//sendErrorResponse - Retrive error message handler
func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(ResponseModel{Message: message, Error: err.Error()}); err != nil {
		panic(err)
	}
}
