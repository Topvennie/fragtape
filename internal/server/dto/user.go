package dto

import (
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/utils"
)

type User struct {
	ID          int    `json:"id"`
	UID         int    `json:"uid"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Crosshair   string `json:"crosshair"`
	Admin       bool   `json:"admin"`
}

func UserDTO(user *model.User) User {
	return User{
		ID:          user.ID,
		UID:         user.UID,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Admin:       user.Admin,
	}
}

func (u *User) ToModel() *model.User {
	user := model.User(*u)
	return &user
}

type UserFilterResult struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

func UserFilterResultDTO(u *model.UserFilterResult) UserFilterResult {
	return UserFilterResult{
		Users: utils.SliceMap(u.Users, func(u model.User) User { return UserDTO(&u) }),
		Total: u.Total,
	}
}

type UserFilter struct {
	Name   string
	Admin  *bool
	Real   *bool
	Limit  int
	Offset int
}

func (u *UserFilter) ToModel() *model.UserFilter {
	return &model.UserFilter{
		Name:   u.Name,
		Admin:  u.Admin,
		Real:   u.Real,
		Limit:  u.Limit,
		Offset: u.Offset,
	}
}
