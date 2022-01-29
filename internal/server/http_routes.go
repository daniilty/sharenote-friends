package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *HTTP) setRoutes(r *mux.Router) {
	const (
		requestsPath       = "/requests"
		acceptRequestsPath = requestsPath + "/accept"
	)

	api := r.PathPrefix("/api/v1/friends").Subrouter()

	api.HandleFunc("",
		h.getFriendsHandler,
	).Methods(http.MethodGet)

	api.HandleFunc(requestsPath,
		h.getFriendRequestsHandler,
	).Methods(http.MethodGet)

	api.HandleFunc(requestsPath,
		h.requestFriendHandler,
	).Methods(http.MethodPost)

	api.HandleFunc(acceptRequestsPath,
		h.acceptFriendHandler,
	).Methods(http.MethodPost)
}
