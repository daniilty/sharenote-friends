package core

import schema "github.com/daniilty/sharenote-grpc-schema"

type User struct {
	ID   string
	Name string
}

func convertPBUsersToInner(uu []*schema.User) []*User {
	converted := make([]*User, 0, len(uu))

	for i := range uu {
		converted = append(converted, convertPBUserToInner(uu[i]))
	}

	return converted
}

func convertPBUserToInner(u *schema.User) *User {
	return &User{
		ID:   u.GetId(),
		Name: u.GetName(),
	}
}
