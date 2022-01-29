package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/daniilty/sharenote-friends/internal/slice"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FriendRequests struct {
	ID        string   `bson:"_id"`
	UID       string   `bson:"uid"`
	FriendIDs []string `bson:"friend_ids"`
}

type Friends struct {
	ID        string   `bson:"_id"`
	UID       string   `bson:"uid"`
	FriendIDs []string `bson:"friend_ids"`
}

func (f *FriendRequests) toBSOND() bson.D {
	return bson.D{
		{Key: "uid", Value: f.UID},
		{Key: "friend_ids", Value: f.FriendIDs},
	}
}

func (f *Friends) toBSOND() bson.D {
	return bson.D{
		{Key: "uid", Value: f.UID},
		{Key: "friend_ids", Value: f.FriendIDs},
	}
}

func (d *DBImpl) GetFriendRequests(ctx context.Context, uid string) (*FriendRequests, error) {
	filter := bson.D{{Key: "uid", Value: uid}}

	res := d.friendRequestsCollection.FindOne(ctx, filter)
	fr := &FriendRequests{}

	err := res.Decode(fr)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}

		fr.UID = uid
	}

	return fr, nil
}

func (d *DBImpl) UpdateFriendRequests(ctx context.Context, fr *FriendRequests) error {

	filter := bson.D{{Key: "uid", Value: fr.UID}}
	update := bson.D{{Key: "$set", Value: fr.toBSOND()}}
	opts := options.Update().SetUpsert(true)

	_, err := d.friendRequestsCollection.UpdateOne(ctx, filter, update, opts)

	return err
}

func (d *DBImpl) GetFriends(ctx context.Context, uid string) (*Friends, error) {
	filter := bson.D{{Key: "uid", Value: uid}}

	res := d.friendsCollection.FindOne(ctx, filter)
	f := &Friends{}

	err := res.Decode(f)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}

		f.UID = uid
	}

	return f, nil
}

func (d *DBImpl) AddFriend(ctx context.Context, from string, to string) (bool, error) {
	session, err := d.mongoDB.Client().StartSession()
	if err != nil {
		return false, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	transaction := d.getAddFriendTransaction(from, to)

	_, err = session.WithTransaction(ctx, transaction)
	if err != nil {
		if errors.Is(err, errNotInFriendRequests) || errors.Is(err, errAlreadyFriends) {
			return true, err
		}

		return false, err
	}

	return true, nil
}

func (d *DBImpl) RemoveFriend(ctx context.Context, from string, to string) (bool, error) {
	session, err := d.mongoDB.Client().StartSession()
	if err != nil {
		return false, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	transaction := d.getRemoveFriendTransaction(from, to)

	_, err = session.WithTransaction(ctx, transaction)
	if err != nil {
		if errors.Is(err, errNotFriends) {
			return true, err
		}

		return false, err
	}

	return true, nil
}

func (d *DBImpl) UpdateFriends(ctx context.Context, f *Friends) error {
	filter := bson.D{{Key: "uid", Value: f.UID}}
	update := bson.D{{Key: "$set", Value: f.toBSOND()}}
	opts := options.Update().SetUpsert(true)

	_, err := d.friendsCollection.UpdateOne(ctx, filter, update, opts)

	return err
}

func (d *DBImpl) getAddFriendTransaction(from string, to string) transactionFunc {
	return func(sessCtx mongo.SessionContext) (interface{}, error) {
		requests, err := d.GetFriendRequests(sessCtx, to)
		if err != nil {
			return false, fmt.Errorf("get friend requests: %w", err)
		}

		if !slice.ContainsString(requests.FriendIDs, from) {
			return true, errNotInFriendRequests
		}

		requests.FriendIDs = slice.RemoveString(requests.FriendIDs, from)

		err = d.UpdateFriendRequests(sessCtx, requests)
		if err != nil {
			return nil, fmt.Errorf("update friend requests: %w", err)
		}

		toFriends, err := d.GetFriends(sessCtx, to)
		if err != nil {
			return nil, fmt.Errorf("get to friends: %w", err)
		}

		if slice.ContainsString(toFriends.FriendIDs, from) {
			return nil, errAlreadyFriends
		}

		fromFriends, err := d.GetFriends(sessCtx, from)
		if err != nil {
			return nil, fmt.Errorf("get from friends: %w", err)
		}

		if slice.ContainsString(fromFriends.FriendIDs, to) {
			return nil, errAlreadyFriends
		}

		// make friends with each other
		toFriends.FriendIDs = append(toFriends.FriendIDs, from)
		fromFriends.FriendIDs = append(fromFriends.FriendIDs, to)

		err = d.UpdateFriends(sessCtx, toFriends)
		if err != nil {
			return nil, fmt.Errorf("update to friends: %w", err)
		}

		err = d.UpdateFriends(sessCtx, fromFriends)
		if err != nil {
			return nil, fmt.Errorf("update from friends: %w", err)
		}

		return nil, nil
	}
}

func (d *DBImpl) getRemoveFriendTransaction(from string, to string) transactionFunc {
	return func(sessCtx mongo.SessionContext) (interface{}, error) {
		toFriends, err := d.GetFriends(sessCtx, to)
		if err != nil {
			return nil, fmt.Errorf("get to friends: %w", err)
		}

		if !slice.ContainsString(toFriends.FriendIDs, from) {
			return nil, errNotFriends
		}

		fromFriends, err := d.GetFriends(sessCtx, from)
		if err != nil {
			return nil, fmt.Errorf("get from friends: %w", err)
		}

		if !slice.ContainsString(fromFriends.FriendIDs, to) {
			return nil, errNotFriends
		}

		// remove friends from each other
		toFriends.FriendIDs = slice.RemoveString(toFriends.FriendIDs, from)
		fromFriends.FriendIDs = slice.RemoveString(toFriends.FriendIDs, to)

		err = d.UpdateFriends(sessCtx, toFriends)
		if err != nil {
			return nil, fmt.Errorf("update to friends: %w", err)
		}

		err = d.UpdateFriends(sessCtx, fromFriends)
		if err != nil {
			return nil, fmt.Errorf("update from friends: %w", err)
		}

		return nil, nil
	}
}
