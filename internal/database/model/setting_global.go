package model

import "github.com/topvennie/fragtape/pkg/sqlc"

type SettingGlobal struct {
	ID             int
	DemoUpload     bool
	CustomCriteria bool
	ChatCommand    bool
	ChatTrigger    string
}

func SettingGlobalModel(s sqlc.SettingGlobal) *SettingGlobal {
	return &SettingGlobal{
		ID:             int(s.ID),
		DemoUpload:     s.DemoUpload,
		CustomCriteria: s.CustomCriteria,
		ChatCommand:    s.ChatCommand,
		ChatTrigger:    s.ChatTrigger,
	}
}
