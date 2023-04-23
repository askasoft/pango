package freshservice

type Article struct {
	ID int64 `json:"id,omitempty"`

	// Title of the solution article
	Title string `json:"title,omitempty"`

	// Description of the solution article
	Description string `json:"description,omitempty"`

	// The rank of the solution article in the article listing
	Position int `json:"position,omitempty"`

	// The type of the article. ( 1 - permanent, 2 - workaround )
	ArticleType int `json:"article_type,omitempty"`

	// ID of the folder to which the solution article belongs
	FolderID int64 `json:"folder_id,omitempty"`

	// ID of the category to which the solution article belongs
	CagetoryID int64 `json:"category_id,omitempty"`

	// Status of the solution article.  ( 1 - draft, 2 - published )
	Status int `json:"status,omitempty"`

	// Approval status of the article.
	ApprovalStatus int `json:"approval_status,omitempty"`

	// Number of upvotes for the solution article
	ThumbsUp int `json:"thumbs_up,omitempty"`

	// Number of down votes for the solution article
	ThumbsDown int `json:"thumbs_down,omitempty"`

	// ID of the agent who created the solution article
	AgentID int64 `json:"agent_id,omitempty"`

	// Number of views for the solution article
	Views int64 `json:"views,omitempty"`

	// Tags that have been associated with the solution article
	Tags []string `json:"tags,omitempty"`

	// Keywords that have been associated with the solution article
	Keywords []string `json:"keywords,omitempty"`

	// Attachments associated with the article. The total size of all of a article's attachments cannot exceed 25MB.
	Attachments []*Attachment `json:"attachments,omitempty"`

	// Article from external url link.
	URL string `json:"url,omitempty"`

	// Date in future when this article would need to be reviewed again.
	ReviewDate *Time `json:"review_date,omitempty"`

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

func (a *Article) Files() Files {
	fs := make(Files, len(a.Attachments))
	for i, a := range a.Attachments {
		fs[i] = a
	}
	return fs
}

func (a *Article) Values() Values {
	vs := Values{}

	vs.SetString("title", a.Title)
	vs.SetString("description", a.Description)
	vs.SetInt("article_type", a.ArticleType)
	vs.SetInt64("folder_id", a.FolderID)
	vs.SetInt("status", a.Status)
	vs.SetStrings("tags", a.Tags)
	vs.SetStrings("keywords", a.Keywords)
	vs.SetTimePtr("review_date", a.ReviewDate)
	return vs
}

func (a *Article) String() string {
	return toString(a)
}

type articleResult struct {
	Article  *Article   `json:"article,omitempty"`
	Articles []*Article `json:"articles,omitempty"`
}
