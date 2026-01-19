package dto

import "github.com/topvennie/fragtape/internal/database/model"

type User struct {
	ID          int    `json:"id"`
	UID         int    `json:"uid"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Crosshair   string `json:"crosshair"`
}

func UserDTO(user *model.User) User {
	return User{
		ID:          user.ID,
		UID:         user.UID,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
	}
}

func (u *User) ToModel() *model.User {
	user := model.User(*u)
	return &user
}
