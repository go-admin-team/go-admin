package dto

type AutoForm struct {
	Fields        []Field `json:"fields"`
	FormRef       string  `json:"formRef"`
	FormModel     string  `json:"formModel"`
	Size          string  `json:"size"`
	LabelPosition string  `json:"labelPosition"`
	LabelWidth    int     `json:"labelWidth"`
	FormRules     string  `json:"formRules"`
	Gutter        int     `json:"gutter"`
	Disabled      bool    `json:"disabled"`
	Span          int     `json:"span"`
	FormBtns      bool    `json:"formBtns"`
}

type Config struct {
	Label        string        `json:"label"`
	LabelWidth   interface{}   `json:"labelWidth"`
	ShowLabel    bool          `json:"showLabel"`
	ChangeTag    bool          `json:"changeTag"`
	Tag          string        `json:"tag"`
	TagIcon      string        `json:"tagIcon"`
	Required     bool          `json:"required"`
	Layout       string        `json:"layout"`
	Span         int           `json:"span"`
	Document     string        `json:"document"`
	RegList      []interface{} `json:"regList"`
	FormId       int           `json:"formId"`
	RenderKey    int64         `json:"renderKey"`
	DefaultValue interface{}   `json:"defaultValue"`
	ShowTip      bool          `json:"showTip,omitempty"`
	ButtonText   string        `json:"buttonText,omitempty"`
	FileSize     int           `json:"fileSize,omitempty"`
	SizeUnit     string        `json:"sizeUnit,omitempty"`
}

type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Slot struct {
	Prepend  string   `json:"prepend,omitempty"`
	Append   string   `json:"append,omitempty"`
	ListType bool     `json:"list-type,omitempty"`
	Options  []Option `json:"options,omitempty"`
}

type Field struct {
	Config        Config      `json:"__config__"`
	Slot          Slot        `json:"__slot__"`
	Placeholder   string      `json:"placeholder,omitempty"`
	Style         Style       `json:"style,omitempty"`
	Clearable     bool        `json:"clearable,omitempty"`
	PrefixIcon    string      `json:"prefix-icon,omitempty"`
	SuffixIcon    string      `json:"suffix-icon,omitempty"`
	Maxlength     interface{} `json:"maxlength"`
	ShowWordLimit bool        `json:"show-word-limit,omitempty"`
	Readonly      bool        `json:"readonly,omitempty"`
	Disabled      bool        `json:"disabled"`
	VModel        string      `json:"__vModel__"`
	Action        string      `json:"action,omitempty"`
	Accept        string      `json:"accept,omitempty"`
	Name          string      `json:"name,omitempty"`
	AutoUpload    bool        `json:"auto-upload,omitempty"`
	ListType      string      `json:"list-type,omitempty"`
	Multiple      bool        `json:"multiple,omitempty"`
	Filterable    bool        `json:"filterable,omitempty"`
}

type Style struct {
	Width string `json:"width"`
}
