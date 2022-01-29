package core

import (
	"context"
	"fmt"

	"github.com/daniilty/sharenote-friends/internal/slice"
	schema "github.com/daniilty/sharenote-grpc-schema"
)

func (s *ServiceImpl) GetFriendRequests(ctx context.Context, uid string) ([]*User, error) {
	reqs, err := s.db.GetFriendRequests(ctx, uid)
	if err != nil {
		return nil, err
	}

	usersResp, err := s.usersClient.GetUsers(ctx, &schema.GetUsersRequest{
		Ids: reqs.FriendIDs,
	})
	if err != nil {
		return nil, err
	}

	return convertPBUsersToInner(usersResp.GetUsers()), nil
}

func (s *ServiceImpl) RequestFriend(ctx context.Context, from string, to string) (bool, error) {
	reqs, err := s.db.GetFriendRequests(ctx, to)
	if err != nil {
		return false, err
	}

	if slice.ContainsString(reqs.FriendIDs, from) {
		return true, fmt.Errorf("user is already in friend requests")
	}

	reqs.FriendIDs = append(reqs.FriendIDs, from)

	err = s.db.UpdateFriendRequests(ctx, reqs)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *ServiceImpl) DeclineFriendRequest(ctx context.Context, from string, to string) (bool, error) {
	reqs, err := s.db.GetFriendRequests(ctx, to)
	if err != nil {
		return false, err
	}

	if !slice.ContainsString(reqs.FriendIDs, from) {
		return true, fmt.Errorf("user is not in friend requests")
	}

	reqs.FriendIDs = slice.RemoveString(reqs.FriendIDs, from)

	err = s.db.UpdateFriendRequests(ctx, reqs)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *ServiceImpl) GetFriends(ctx context.Context, uid string) ([]*User, error) {
	friends, err := s.db.GetFriends(ctx, uid)
	if err != nil {
		return nil, err
	}

	usersResp, err := s.usersClient.GetUsers(ctx, &schema.GetUsersRequest{
		Ids: friends.FriendIDs,
	})
	if err != nil {
		return nil, err
	}

	return convertPBUsersToInner(usersResp.GetUsers()), nil
}

func (s *ServiceImpl) AddFriend(ctx context.Context, from string, to string) (bool, error) {
	return s.db.AddFriend(ctx, from, to)
}

func (s *ServiceImpl) RemoveFriend(ctx context.Context, from string, to string) (bool, error) {
	return s.db.RemoveFriend(ctx, from, to)
}
