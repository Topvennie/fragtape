// Package model contains all databank models
package model

import "github.com/topvennie/fragtape/pkg/sqlc"

type User struct {
	ID          int
	UID         int
	Name        string
	DisplayName string
	AvatarURL   string
	Crosshair   string
	Admin       bool
}

func UserModel(user sqlc.User) *User {
	return &User{
		ID:          int(user.ID),
		UID:         int(user.Uid),
		Name:        fromString(user.Name),
		DisplayName: user.DisplayName,
		AvatarURL:   fromString(user.AvatarUrl),
		Crosshair:   fromString(user.Crosshair),
		Admin:       user.Admin,
	}
}

// EqualEntry returns true if all non unique values are equal
func (u *User) EqualEntry(u2 User) bool {
	return u.Name == u2.Name && u.DisplayName == u2.DisplayName && u.AvatarURL == u2.AvatarURL
}

type UserFilterResult struct {
	Users []User
	Total int
}

type UserFilter struct {
	Name   string
	Admin  *bool
	Real   *bool
	Limit  int
	Offset int
}
