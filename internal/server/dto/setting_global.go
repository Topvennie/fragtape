package dto

import "github.com/topvennie/fragtape/internal/database/model"

type SettingGlobal struct {
	ID             int    `json:"-"`
	DemoUpload     bool   `json:"demo_upload"`
	CustomCriteria bool   `json:"custom_criteria"`
	ChatCommand    bool   `json:"chat_command"`
	ChatTrigger    string `json:"chat_trigger" validate:"required,min=3"`
}

func SettingGlobalDTO(s *model.SettingGlobal) SettingGlobal {
	return SettingGlobal(*s)
}

func (s *SettingGlobal) ToModel() *model.SettingGlobal {
	setting := model.SettingGlobal(*s)
	return &setting
}
