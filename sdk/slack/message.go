package slack

// Field slack field
type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short string `json:"short,omitempty"`
}

// Attachment slack attachment
type Attachment struct {
	AuthorIcon string   `json:"author_icon,omitempty"`
	AuthorLink string   `json:"author_link,omitempty"`
	AuthorName string   `json:"author_name,omitempty"`
	Fallback   string   `json:"fallback,omitempty"`
	Color      string   `json:"color,omitempty"`
	Pretext    string   `json:"pretext,omitempty"`
	Title      string   `json:"title,omitempty"`
	TitleLink  string   `json:"title_link,omitempty"`
	Text       string   `json:"text,omitempty"`
	ImageURL   string   `json:"image_url,omitempty"`
	ThumbURL   string   `json:"thumb_url,omitempty"`
	Footer     string   `json:"footer,omitempty"`
	FooterIcon string   `json:"footer_icon,omitempty"`
	Fields     []*Field `json:"fields,omitempty"`
}

// AddField add a field
func (sm *Attachment) AddField(sf *Field) {
	sm.Fields = append(sm.Fields, sf)
}

// Message slack message
type Message struct {
	IconEmoji   string        `json:"icon_emoji,omitempty"`
	Text        string        `json:"text,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
	Markdown    bool          `json:"mrkdwn,omitempty"`
}

// AddAttachment add a attachment
func (sm *Message) AddAttachment(sa *Attachment) {
	sm.Attachments = append(sm.Attachments, sa)
}
