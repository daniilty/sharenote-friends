package core

import (
	"context"

	"github.com/daniilty/sharenote-friends/internal/mongo"
	schema "github.com/daniilty/sharenote-grpc-schema"
)

type Service interface {
	// GetFriendRequests - get hser friend request uids.
	GetFriendRequests(context.Context, string) ([]*User, error)
	// RequestFriend - add friend request to user.
	RequestFriend(context.Context, string, string) (bool, error)
	// DeclineFriendRequest decline request from some user.
	DeclineFriendRequest(context.Context, string, string) (bool, error)
	// GetFriends - get user friends.
	GetFriends(context.Context, string) ([]*User, error)
	// AddFriend - add friend from friend request.
	AddFriend(context.Context, string, string) (bool, error)
	// RemoveFriend - remove friend.
	RemoveFriend(context.Context, string, string) (bool, error)
}

type ServiceImpl struct {
	usersClient schema.UsersClient
	db          mongo.DB
}

func NewService(db mongo.DB, usersClient schema.UsersClient) Service {
	return &ServiceImpl{
		usersClient: usersClient,
		db:          db,
	}
}
