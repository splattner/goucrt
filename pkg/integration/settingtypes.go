package integration

type SettingTypeNumber struct {
	Number SettingTypeNumberDefinition `json:"number"`
}

type SettingTypeNumberDefinition struct {
	Value    float64      `json:"value"`
	Min      float64      `json:"min,omitempty"`
	Max      float64      `json:"max,omitempty"`
	Step     float64      `json:"step,omitempty"`
	Decimals int          `json:"decimal,omitempty"`
	Unit     LanguageText `json:"unit"`
}

type SettingTypeText struct {
	Text SettingTypeTextDefinition `json:"text"`
}

type SettingTypeTextDefinition struct {
	Value string `json:"value"`
	Regex string `json:"regex,omitempty"`
}

type SettingTypeTextArea struct {
	TextArea SettingTypeTextAreaDefinition `json:"textarea"`
}

type SettingTypeTextAreaDefinition struct {
	Value string `json:"value"`
}

type SettingTypePassword struct {
	Password SettingTypePasswordDefinition `json:"password"`
}

type SettingTypePasswordDefinition struct {
	Value string `json:"value"`
	Regex string `json:"regex,omitempty"`
}

type SettingTypeCheckbox struct {
	Checkbox SettingTypeCheckboxDefinition `json:"checkbox"`
}

type SettingTypeCheckboxDefinition struct {
	Value string `json:"value"`
}

type SettingTypeDropdown struct {
	Dropdown SettingTypeDropdowDefinition `json:"dropdown"`
}

type SettingTypeDropdowDefinition struct {
	Value string                              `json:"value"`
	Items []SettingTypeDropdowItemsDefinition `json:"items"`
}
type SettingTypeDropdowItemsDefinition struct {
	Id    string       `json:"id"`
	Label LanguageText `json:"label"`
}

type SettingTypeLabel struct {
	Label SettingTypeLabelDefinition `json:"label"`
}

type SettingTypeLabelDefinition struct {
	Value LanguageText `json:"value"`
}
