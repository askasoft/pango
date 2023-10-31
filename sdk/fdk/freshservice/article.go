package freshservice

import (
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

type ArticleType int
type ArticleStatus int

const (
	ArticleTypePermanent  ArticleType = 1
	ArticleTypeWorkaround ArticleType = 2

	ArticleStatusDraft     ArticleStatus = 1
	ArticleStatusPublished ArticleStatus = 2
)

func (at ArticleType) String() string {
	switch at {
	case ArticleTypePermanent:
		return "Permanent"
	case ArticleTypeWorkaround:
		return "Workaround"
	default:
		return num.Itoa(int(at))
	}
}

func ParseArticleType(s string) ArticleType {
	switch str.ToLower(s) {
	case "permanent":
		return ArticleTypePermanent
	case "workaround":
		return ArticleTypeWorkaround
	default:
		return 0
	}
}

func (as ArticleStatus) String() string {
	switch as {
	case ArticleStatusDraft:
		return "Draft"
	case ArticleStatusPublished:
		return "Published"
	default:
		return num.Itoa(int(as))
	}
}

func ParseArticleStatus(s string) ArticleStatus {
	switch str.ToLower(s) {
	case "draft":
		return ArticleStatusDraft
	case "published":
		return ArticleStatusPublished
	default:
		return 0
	}
}

type Article struct {
	ID int64 `json:"id,omitempty"`

	// Title of the solution article
	Title string `json:"title,omitempty"`

	// Description of the solution article
	Description string `json:"description,omitempty"`

	// The rank of the solution article in the article listing
	Position int `json:"position,omitempty"`

	// The type of the article. ( 1 - permanent, 2 - workaround )
	ArticleType ArticleType `json:"article_type,omitempty"`

	// ID of the folder to which the solution article belongs
	FolderID int64 `json:"folder_id,omitempty"`

	// ID of the category to which the solution article belongs
	CagetoryID int64 `json:"category_id,omitempty"`

	// Status of the solution article.  ( 1 - draft, 2 - published )
	Status ArticleStatus `json:"status,omitempty"`

	// Approval status of the article.
	ApprovalStatus ApprovalStatus `json:"approval_status,omitempty"`

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

	ModifiedBy int64 `json:"modified_by,omitempty"`

	ModifiedAt *Time `json:"modified_at,omitempty"`

	InsertedIntoTickets int `json:"inserted_into_tickets,omitempty"`

	// Date in future when this article would need to be reviewed again.
	ReviewDate *Time `json:"review_date,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (a *Article) AddAttachment(path string, data ...[]byte) {
	aa := NewAttachment(path, data...)
	a.Attachments = append(a.Attachments, aa)
}

func (a *Article) Files() Files {
	return ((Attachments)(a.Attachments)).Files()
}

func (a *Article) Values() Values {
	vs := Values{}

	vs.SetString("title", a.Title)
	vs.SetString("description", a.Description)
	vs.SetInt("article_type", (int)(a.ArticleType))
	vs.SetInt64("folder_id", a.FolderID)
	vs.SetInt("status", (int)(a.Status))
	vs.SetStrings("tags", a.Tags)
	vs.SetStrings("keywords", a.Keywords)
	vs.SetTimePtr("review_date", a.ReviewDate)
	return vs
}

func (a *Article) String() string {
	return toString(a)
}

type articleResult struct {
	Article *Article `json:"article,omitempty"`
}

type ArticleInfo struct {
	Article

	// Attachments associated with the article. The total size of all of a article's attachments cannot exceed 25MB.
	Attachments []string `json:"attachments,omitempty"`

	GroupFolderGroupIDs          []int64 `json:"group_folder_group_ids,omitempty"`
	FolderDepartmentIDs          []int64 `json:"folder_department_ids,omitempty"`
	GroupFolderRequesterGroupIDs []int64 `json:"group_folder_requester_group_ids,omitempty"`
	GroupFolderDepartmentIDs     []int64 `json:"group_folder_department_ids,omitempty"`

	FolderVisibility FolderVisibility `json:"folder_visibility,omitempty"`
}

func (ai *ArticleInfo) normalize() {
	if len(ai.Attachments) > 0 {
		ai.Article.Attachments = make([]*Attachment, len(ai.Attachments))
		for i, s := range ai.Attachments {
			ai.Article.Attachments[i] = &Attachment{Name: s}
		}
	}
}

func (ai *ArticleInfo) String() string {
	return toString(ai)
}

type articlesResult struct {
	Articles []*ArticleInfo `json:"articles,omitempty"`
}
