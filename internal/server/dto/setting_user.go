package dto

import "github.com/topvennie/fragtape/internal/database/model"

type SettingUser struct {
	ConnectedSteam  bool `json:"connected_steam"`
	ConnectedFaceit bool `json:"connected_faceit"`
}

func SettingUserDTO(s *model.SettingUser) SettingUser {
	return SettingUser{
		ConnectedSteam:  s.SteamMatchToken != "" && s.SteamAuthenticationToken != "",
		ConnectedFaceit: s.FaceitID != "",
	}
}

type SettingUserSteam struct {
	MatchToken          string `json:"match_token" validate:"required,steammatchtoken"`
	AuthenticationToken string `json:"authentication_token" validate:"required,steamauthtoken"`
}
