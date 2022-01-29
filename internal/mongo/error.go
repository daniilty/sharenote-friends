package mongo

import "errors"

var (
	errNotInFriendRequests = errors.New("not in friend requests")
	errAlreadyFriends      = errors.New("users are already friends")
	errNotFriends          = errors.New("users are not friends")
)
