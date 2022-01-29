package server

import (
	"net/http"

	"github.com/daniilty/sharenote-friends/internal/core"
)

func convertCoreUsersToResponse(uu []*core.User) *friendsResponse {
	data := make([]*friend, 0, len(uu))

	for i := range uu {
		data = append(data, convertCoreUserToResponse(uu[i]))
	}

	return &friendsResponse{
		Status: http.StatusText(http.StatusOK),
		Data:   data,
	}
}

func convertCoreUserToResponse(u *core.User) *friend {
	return &friend{
		ID:   u.ID,
		Name: u.Name,
	}
}
