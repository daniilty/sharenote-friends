package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ DB = (*DBImpl)(nil)

type DB interface {
	// GetNote - get user by id.
	GetFriendRequests(context.Context, string) (*FriendRequests, error)
	// UpdateFriendRequests - update or insert friend requests for user.
	UpdateFriendRequests(context.Context, *FriendRequests) error
	// RemoveUser - remove user's requests and friends.
	RemoveUser(context.Context, string) error
	// GetFriends - get user friends.
	GetFriends(context.Context, string) (*Friends, error)
	// AddFriend - make transaction and add user to friends list.
	AddFriend(context.Context, string, string) (bool, error)
	// RemoveFriend - make transaction and remove friend from each other list.
	RemoveFriend(context.Context, string, string) (bool, error)
	// UpdateFriends - update user friends.
	UpdateFriends(context.Context, *Friends) error
}

type DBImpl struct {
	mongoDB                  *mongo.Database
	friendRequestsCollection *mongo.Collection
	friendsCollection        *mongo.Collection
}

func NewDBImpl(db *mongo.Database, friendRequestsCollection *mongo.Collection, friendsCollection *mongo.Collection) *DBImpl {
	return &DBImpl{
		mongoDB:                  db,
		friendsCollection:        friendsCollection,
		friendRequestsCollection: friendRequestsCollection,
	}
}

func Connect(ctx context.Context, addr string) (*mongo.Client, error) {
	return mongo.Connect(ctx, options.Client().ApplyURI(addr))
}

func InitIndex(ctx context.Context, collection *mongo.Collection) error {
	const minIndexesLen = 1

	specs, err := collection.Indexes().ListSpecifications(ctx)
	if err != nil {
		return err
	}

	if len(specs) < minIndexesLen {
		return nil
	}

	opts := options.Index().SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"uid": 1}, Options: opts}

	name, err := collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	log.Println("index created", name)

	return nil
}
