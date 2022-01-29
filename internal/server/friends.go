package server

import (
	"fmt"
	"net/http"

	"github.com/daniilty/sharenote-auth/claims"
)

type friend struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type friendsResponse struct {
	Status string    `json:"status"`
	Data   []*friend `json:"data"`
}

type requestFriendRequest struct {
	FriendID string `json:"friend_id"`
}

func (r *requestFriendRequest) validate() error {
	if r.FriendID == "" {
		return fmt.Errorf(`"friend_id": cannot be empty`)
	}

	return nil
}

type addFriendRequest struct {
	FriendID string `json:"friend_id"`
}

func (r *addFriendRequest) validate() error {
	if r.FriendID == "" {
		return fmt.Errorf(`"friend_id": cannot be empty`)
	}

	return nil
}

func (f *friendsResponse) writeJSON(w http.ResponseWriter) error {
	return writeJSONResponse(w, http.StatusOK, f)
}

func (h *HTTP) getFriendsHandler(w http.ResponseWriter, r *http.Request) {
	resp := h.getFriendsResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getFriendRequestsHandler(w http.ResponseWriter, r *http.Request) {
	resp := h.getFriendRequestsResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) requestFriendHandler(w http.ResponseWriter, r *http.Request) {
	resp := h.getRequestFriendResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) acceptFriendHandler(w http.ResponseWriter, r *http.Request) {
	resp := h.getAcceptFriendResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getFriendsResponse(r *http.Request) response {
	c, err := claims.ParseHTTPHeader(r.Header)
	if err != nil {
		return getUnauthorizedErrorResponse()
	}

	friends, err := h.service.GetFriends(r.Context(), c.UID)
	if err != nil {
		h.logger.Errorw("Get Friends.", "err", err)

		return getInternalServerErrorResponse()
	}

	return convertCoreUsersToResponse(friends)
}

func (h *HTTP) getFriendRequestsResponse(r *http.Request) response {
	c, err := claims.ParseHTTPHeader(r.Header)
	if err != nil {
		return getUnauthorizedErrorResponse()
	}

	friends, err := h.service.GetFriendRequests(r.Context(), c.UID)
	if err != nil {
		h.logger.Errorw("Get Friend requests.", "err", err)

		return getInternalServerErrorResponse()
	}

	return convertCoreUsersToResponse(friends)
}

func (h *HTTP) getRequestFriendResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no body")
	}

	c, err := claims.ParseHTTPHeader(r.Header)
	if err != nil {
		return getUnauthorizedErrorResponse()
	}

	req := &requestFriendRequest{}

	err = unmarshalReader(r.Body, req)
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error())
	}

	err = req.validate()
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error())
	}

	if req.FriendID == c.UID {
		return getBadRequestWithMsgResponse("you cannot be friend with yourself")
	}

	ok, err := h.service.RequestFriend(r.Context(), c.UID, req.FriendID)
	if err != nil {
		if ok {
			return getBadRequestWithMsgResponse(err.Error())
		}

		h.logger.Errorw("Request Friend.", "err", err)

		return getInternalServerErrorResponse()
	}

	return getEmptyOKResponse()
}

func (h *HTTP) getAcceptFriendResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no body")
	}

	c, err := claims.ParseHTTPHeader(r.Header)
	if err != nil {
		return getUnauthorizedErrorResponse()
	}

	req := &addFriendRequest{}

	err = unmarshalReader(r.Body, req)
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error())
	}

	err = req.validate()
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error())
	}

	ok, err := h.service.AddFriend(r.Context(), req.FriendID, c.UID)
	if err != nil {
		if ok {
			return getBadRequestWithMsgResponse(err.Error())
		}

		h.logger.Errorw("Add friend.", "err", err)

		return getInternalServerErrorResponse()
	}

	return getEmptyOKResponse()
}
