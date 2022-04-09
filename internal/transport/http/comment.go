package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Cingozr/go-rest-api/internal/comment"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	ContentType = []string{"Content-Type", "application/json; charset=UTF-8"}
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
			}).Info("Handler Request")
		next.ServeHTTP(w, r)

	})
}

func BasicAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("basic auth endpoint hit")
		user, pass, ok := r.BasicAuth()
		if user == "admin" && pass == "password" && ok {
			log.Infof("User: %s - Password: %s ", user, pass)
			original(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
		}
	}
}

func validateToken(accessToken string) bool {
	var mySigningKey = []byte("cingoz_recai")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there has been an error")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}

func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("jwt authentication hit")
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			w.Header().Set(ContentType[0], ContentType[1])
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			w.Header().Set(ContentType[0], ContentType[1])
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
		}

		if validateToken(authHeaderParts[1]) {
			original(w, r)
		} else {
			w.Header().Set(ContentType[0], ContentType[1])
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
		}
	}
}

func (h *Handler) SetupRoutes() {
	log.Info("Setting Up Routes")
	h.Router = mux.NewRouter()
	h.Router.Use(LoggingMiddleware)

	//Comment Service Routes Start
	h.Router.HandleFunc("/api/comment", JWTAuth(h.GetAllComments)).Methods("GET")
	h.Router.HandleFunc("/api/comment", JWTAuth(h.PostComment)).Methods("POST")
	h.Router.HandleFunc("/api/comment/{id}", JWTAuth(h.GetComment)).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", JWTAuth(h.UpdateComment)).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", JWTAuth(h.DeleteComment)).Methods("DELETE")
	//Comment Service Routes End

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := sendOkResponse(w, ResponseModel{Message: "I'm Alive"}); err != nil {
			panic(err)
		}
	})
}

//GetAllComments - Rerieves all comments
func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.Service.GetAllComments()
	if err != nil {
		sendErrorResponse(w, "Failed retrieve all comments", err)
		return
	}

	if err := sendOkResponse(w, comments); err != nil {
		panic(err)
	}
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var comment comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		sendErrorResponse(w, "Failed to decode Json Body", err)
		return
	}

	comment, err := h.Service.PostComment(comment)
	if err != nil {
		sendErrorResponse(w, "Failed to post new comment", err)
		return
	}

	if err := sendOkResponse(w, comment); err != nil {
		panic(err)
	}
}

func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	commentId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UINT from ID", err)
		return
	}
	comment, err := h.Service.GetComment(uint(commentId))
	if err != nil {
		sendErrorResponse(w, "Error Retrieving Comment By ID", err)
		return
	}

	if err := sendOkResponse(w, comment); err != nil {
		panic(err)
	}
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		sendErrorResponse(w, "Failed to decode JSON body", err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	updateId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UINT from ID", err)
		return
	}

	comment, err = h.Service.UpdateComment(uint(updateId), comment)
	if err != nil {
		sendErrorResponse(w, "Failed to update comment", err)
		return
	}

	if err := sendOkResponse(w, comment); err != nil {
		panic(err)
	}
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	deleteId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UINT from ID", err)
		return
	}

	if err = h.Service.DeleteComment(uint(deleteId)); err != nil {
		sendErrorResponse(w, "Failed to delete comment by id", err)
	}

	if err := sendOkResponse(w, ResponseModel{Message: "Successfully Deleted"}); err != nil {
		panic(err)
	}
}
