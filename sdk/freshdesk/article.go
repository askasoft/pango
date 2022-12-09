package freshdesk

type Article struct {
	ID int64 `json:"id,omitempty"`

	// ID of the agent who created the solution article
	AgentID int64 `json:"agent_id,omitempty"`

	// ID of the category to which the solution article belongs
	CagetoryID int64 `json:"category_id,omitempty"`

	// Title of the solution article
	Title string `json:"title,omitempty"`

	// Description of the solution article
	Description string `json:"description,omitempty"`

	// Description of the solution article in plain text
	DescriptionText string `json:"description_text,omitempty"`

	// ID of the folder to which the solution article belongs
	FolderID int64 `json:"folder_id,omitempty"`

	// Parent category and folders in which the article is placed
	Hierarchy []map[string]any `json:"hierarchy,omitempty"`

	// Number of views for the solution article
	Hits int64 `json:"hits,omitempty"`

	// Status of the solution article
	Status int `json:"status,omitempty"`

	// Meta data for search engine optimization. Allows meta_title, meta_description and meta_keywords
	SeoData map[string]string `json:"seo_data,omitempty"`

	// Tags that have been associated with the solution article
	Tags []string `json:"tags,omitempty"`

	// Number of down votes for the solution article
	ThumbsDown int `json:"thumbs_down,omitempty"`

	// Number of upvotes for the solution article
	ThumbsUp int `json:"thumbs_up,omitempty"`

	// Attachments associated with the article. The total size of all of a article's attachments cannot exceed 25MB.
	Attachments []*Attachment `json:"attachments,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (a *Article) AddAttachment(path string, data ...[]byte) {
	aa := NewAttachment(path, data...)
	a.Attachments = append(a.Attachments, aa)
}

func (a *Article) GetAttachments() []*Attachment {
	return a.Attachments
}

func (a *Article) Values() Values {
	vs := Values{}

	vs.SetString("title", a.Title)
	vs.SetString("description", a.Description)
	vs.SetInt("status", a.Status)
	vs.SetStrings("tags", a.Tags)
	vs.SetMap("seo_data", a.SeoData)
	return vs
}

func (a *Article) String() string {
	return toString(a)
}

type ArticleEx struct {
	Article

	Path         string `json:"path,omitempty"`
	LanguageID   int    `json:"language_id,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
	FolderName   string `json:"folder_name,omitempty"`
}

func (a *ArticleEx) String() string {
	return toString(a)
}
