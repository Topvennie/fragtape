// Package model contains all databank models
package model

import "github.com/topvennie/fragtape/pkg/sqlc"

type User struct {
	ID          int
	UID         string
	Name        string
	DisplayName string
	AvatarURL   string
}

func UserModel(user sqlc.User) *User {
	return &User{
		ID:          int(user.ID),
		UID:         user.Uid,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarUrl,
	}
}

// EqualEntry returns true if all non unique values are equal
func (u *User) EqualEntry(u2 User) bool {
	return u.Name == u2.Name && u.DisplayName == u2.DisplayName && u.AvatarURL == u2.AvatarURL
}
