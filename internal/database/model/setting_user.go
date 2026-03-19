package model

import "github.com/topvennie/fragtape/pkg/sqlc"

type SettingUser struct {
	ID                       int
	UserID                   int
	SteamMatchToken          string
	SteamAuthenticationToken string
	SteamImportOld           bool
	FaceitID                 string
	FirstTimeWizard          bool
}

func SettingUserModel(s sqlc.SettingUser) *SettingUser {
	return &SettingUser{
		ID:                       int(s.ID),
		UserID:                   int(s.UserID),
		SteamMatchToken:          fromString(s.SteamMatchToken),
		SteamAuthenticationToken: fromString(s.SteamAuthenticationToken),
		SteamImportOld:           fromBool(s.SteamImportOld),
		FaceitID:                 fromString(s.FaceitID),
		FirstTimeWizard:          s.FirstTimeWizard,
	}
}
