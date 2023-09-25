package freshdesk

import (
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

type ArticleStatus int
type ArticleHierarchyType string

const (
	ArticleStatusDraft     ArticleStatus = 1
	ArticleStatusPublished ArticleStatus = 2

	ArticleHierarchyTypeCategory ArticleHierarchyType = "category"
	ArticleHierarchyTypeFolder   ArticleHierarchyType = "folder"
)

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

type ArticleSeoData struct {
	MetaTitle       string `json:"meta_title,omitempty"`
	MetaKeywords    string `json:"meta_keywords,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`
}

type ArticleHierarchyData struct {
	ID       int64  `json:"id,omitempty"`
	Language string `json:"language,omitempty"`
	Name     string `json:"name,omitempty"`
}

type ArticleHierarchyItem struct {
	Data  *ArticleHierarchyData `json:"data,omitempty"`
	Level int                   `json:"level,omitempty"`
	Type  ArticleHierarchyType  `json:"type,omitempty"`
}

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
	Hierarchy []*ArticleHierarchyItem `json:"hierarchy,omitempty"`

	// Number of views for the solution article
	Hits int64 `json:"hits,omitempty"`

	// Status of the solution article
	Status ArticleStatus `json:"status,omitempty"`

	// Meta data for search engine optimization. Allows meta_title, meta_description and meta_keywords
	SeoData *ArticleSeoData `json:"seo_data,omitempty"`

	// Tags that have been associated with the solution article
	Tags []string `json:"tags,omitempty"`

	// Omnichannel: available for chat platforms ("web", "ios", "android")
	Platforms []string `json:"platforms,omitempty"`

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

func (a *Article) Files() Files {
	return ((Attachments)(a.Attachments)).Files()
}

func (a *Article) Values() Values {
	vs := Values{}

	vs.SetString("title", a.Title)
	vs.SetString("description", a.Description)
	vs.SetInt("status", (int)(a.Status))
	vs.SetStrings("tags", a.Tags)
	vs.SetStrings("platforms", a.Platforms)
	if a.SeoData != nil {
		vs.SetString("seo_data[meta_title]", a.SeoData.MetaTitle)
		vs.SetString("seo_data[meta_keywords]", a.SeoData.MetaKeywords)
		vs.SetString("seo_data[meta_description]", a.SeoData.MetaDescription)
	}
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
