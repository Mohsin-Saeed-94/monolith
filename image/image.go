package image

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/turnkeyca/monolith/auth"
)

type Handler struct {
	logger *log.Logger
	client *Client
}

func NewHandler(logger *log.Logger, client *Client) *Handler {
	return &Handler{
		logger: logger,
		client: client,
	}
}

type GenericError struct {
	Message string `json:"message"`
}

func ConfigureImageRoutes(router *mux.Router, logger *log.Logger, client *Client, authenticator *auth.Authenticator) {
	imageHandler := NewHandler(logger, client)
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/api/image/{id:%s}", imageHandler.HandleGetImage)
	getRouter.Use(authenticator.AuthenticateHttp, imageHandler.GetIdFromPath)
	putRouter := router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/api/image", imageHandler.HandlePutImage)
	putRouter.Use(authenticator.AuthenticateHttp, imageHandler.GetBody)
	deleteRouter := router.Methods(http.MethodPut).Subrouter()
	deleteRouter.HandleFunc("/api/image/{id:%s}", imageHandler.HandleDeleteImage)
	deleteRouter.Use(authenticator.AuthenticateHttp, imageHandler.GetIdFromPath)

}
