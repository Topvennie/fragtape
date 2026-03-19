package dto

import "github.com/topvennie/fragtape/internal/database/model"

type SettingUser struct {
	ConnectedSteam  bool `json:"connected_steam"`
	ConnectedFaceit bool `json:"connected_faceit"`
	FirstTimeWizard bool `json:"first_time_wizard"`
}

func SettingUserDTO(s *model.SettingUser) SettingUser {
	return SettingUser{
		ConnectedSteam:  s.SteamMatchToken != "" && s.SteamAuthenticationToken != "",
		ConnectedFaceit: s.FaceitID != "",
		FirstTimeWizard: s.FirstTimeWizard,
	}
}

type SettingUserSteam struct {
	MatchToken          string `json:"match_token" validate:"required,steammatchtoken"`
	AuthenticationToken string `json:"authentication_token" validate:"required,steamauthtoken"`
	ImportOld           *bool  `json:"import_old" validate:"required"`
}

type SettingUserFirsTimeWizard struct {
	FirsTimeWizard *bool `json:"first_time_wizard" validate:"required"`
}
