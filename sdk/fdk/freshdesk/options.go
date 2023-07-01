package freshdesk

type ExportFields struct {
	DefaultFields []string `json:"default_fields,omitempty"`

	CustomFields []string `json:"custom_fields,omitempty"`
}

func (ef *ExportFields) String() string {
	return toString(ef)
}

type ExportOption struct {
	Fields *ExportFields `json:"fields,omitempty"`
}

func (eo *ExportOption) String() string {
	return toString(eo)
}
