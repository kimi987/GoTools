package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/settings"
)

var settingTypes []shared_proto.SettingType

func init() {
	for idx := range shared_proto.SettingType_name {
		if idx != 0 {
			settingTypes = append(settingTypes, shared_proto.SettingType(idx))
		}
	}
}

func NewHeroSettings(defaultSettings []shared_proto.SettingType) *HeroSettings {
	h := &HeroSettings {
		settings: make([]bool, len(settingTypes)),
		privacySettings: make(map[shared_proto.PrivacySettingType]struct{}),
	}

	h.ResetDefault(defaultSettings)
	h.SetPrivacySettings(settings.DefaultPrivacySettings)

	return h
}

// 英雄设置
type HeroSettings struct {
	settings []bool // 设置的值
	// 隐私设置
	privacySettings map[shared_proto.PrivacySettingType]struct{}
}

func (h *HeroSettings) IsPrivacySettingOpen(t shared_proto.PrivacySettingType) bool {
	_, ok := h.privacySettings[t]
	return ok
}

func (h *HeroSettings) TrySetPrivacySetting(t shared_proto.PrivacySettingType, isOpen bool) bool {
	_, ok := h.privacySettings[t]
	if isOpen {
		if ok {
			return false
		}
		h.privacySettings[t] = struct{}{}
	} else {
		if !ok {
			return false
		}
		delete(h.privacySettings, t)
	}
	return true
}

func (h *HeroSettings) IsValidType(settingType shared_proto.SettingType) bool {
	return settingType > 0 && int(settingType) <= len(h.settings)
}

func (h *HeroSettings) Set(settingType shared_proto.SettingType, open bool) {
	if !h.IsValidType(settingType) {
		logrus.Debugf("设置非法: %v", settingType)
		return
	}

	h.settings[settingType-1] = open
}

func (h *HeroSettings) IsOpen(settingType shared_proto.SettingType) bool {
	if !h.IsValidType(settingType) {
		logrus.Debugf("设置非法: %v", settingType)
		return false
	}

	return h.settings[settingType-1]
}

func (h *HeroSettings) SetPrivacySettings(settings []shared_proto.PrivacySettingType) (changed []shared_proto.PrivacySettingType, isOpen []bool) {
	if len(settings) <= 0 {
		return
	}
	if len(h.privacySettings) <= 0 {
		for _, t := range settings {
			h.privacySettings[t] = struct{}{}
			changed = append(changed, t)
			isOpen = append(isOpen, true)
		}
	} else {
		newSettings := make(map[shared_proto.PrivacySettingType]struct{})
		for _, t := range settings {
			newSettings[t] = struct{}{}
		}
		for t, _ := range h.privacySettings {
			if _, ok := newSettings[t]; ok {
				continue
			}
			delete(h.privacySettings, t)
			changed = append(changed, t)
			isOpen = append(isOpen, false)
		}
		for _, t := range settings {
			if !h.TrySetPrivacySetting(t, true) {
				continue
			}
			changed = append(changed, t)
			isOpen = append(isOpen, true)
		}

	}
	return
}

func (h *HeroSettings) ResetDefault(defaultSettings []shared_proto.SettingType) (result uint64) {
	for idx := range h.settings {
		h.settings[idx] = false
	}
	for _, settingType := range defaultSettings {
		h.Set(settingType, true)
	}
	return
}

func (h *HeroSettings) EncodeToUint64() (result uint64) {
	for i := uint64(0); i < uint64(len(h.settings)); i++ {
		result |= SettingTypeShiftValue(shared_proto.SettingType(i + 1))
	}

	return
}

func (h *HeroSettings) Encode() *shared_proto.HeroSettingsProto {
	proto := &shared_proto.HeroSettingsProto{}
	proto.Settings = h.settings
	if len(h.privacySettings) > 0 {
		for t, _ := range h.privacySettings {
			proto.PrivacySettings = append(proto.PrivacySettings, t)
		}
	}
	return proto
}

func (h *HeroSettings) unmarshal(proto *shared_proto.HeroSettingsProto) {
	if proto == nil {
		return
	}

	copy(h.settings, proto.Settings)
	h.SetPrivacySettings(proto.PrivacySettings)
}

func IsOpen(value uint64, settingType shared_proto.SettingType) bool {
	return value&SettingTypeShiftValue(settingType) != 0
}

func SettingTypeShiftValue(settingType shared_proto.SettingType) (result uint64) {
	return 1 << uint64(settingType)
}
