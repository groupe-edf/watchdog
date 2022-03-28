package models

import "strconv"

type SettingType string

const (
	SettingBoolType SettingType = "boolean"
	SettingIntType  SettingType = "integer"
)

type Setting struct {
	ID            int64       `json:"id"`
	ContainerID   int64       `json:"container_id"`
	ContainerType string      `json:"container_type"`
	SettingKey    string      `json:"setting_key"`
	SettingType   SettingType `json:"setting_type"`
	SettingValue  string      `json:"setting_value"`
}

func (setting *Setting) CastValue() interface{} {
	switch setting.SettingType {
	case SettingBoolType:
		value, _ := strconv.ParseBool(setting.SettingValue)
		return value
	case SettingIntType:
		value, _ := strconv.ParseInt(setting.SettingValue, 10, 64)
		return value
	default:
		return setting.SettingValue
	}
}

type Settings []Setting
