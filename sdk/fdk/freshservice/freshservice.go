package freshservice

import (
	"fmt"

	"github.com/pandafw/pango/sdk/fdk"
)

type FreshService fdk.FDK

func (fs *FreshService) doGet(url string, result any) error {
	return (*fdk.FDK)(fs).DoGet(url, result)
}

func (fs *FreshService) doList(url string, lo ListOption, ap any) (bool, error) {
	return (*fdk.FDK)(fs).DoList(url, lo, ap)
}

func (fs *FreshService) doPost(url string, source, result any) error {
	return (*fdk.FDK)(fs).DoPost(url, source, result)
}

func (fs *FreshService) doPut(url string, source, result any) error {
	return (*fdk.FDK)(fs).DoPut(url, source, result)
}

func (fs *FreshService) doDelete(url string) error {
	return (*fdk.FDK)(fs).DoDelete(url)
}

func (fs *FreshService) Download(url string) ([]byte, error) {
	return (*fdk.FDK)(fs).DoDownload(url)
}

func (fs *FreshService) SaveFile(url string, filename string) error {
	return (*fdk.FDK)(fs).DoSave(url, filename)
}

// GetHelpdeskAttachmentURL return a permlink for helpdesk attachment/avator URL
func (fs *FreshService) GetHelpdeskAttachmentURL(aid int64) string {
	return fmt.Sprintf("%s/helpdesk/attachments/%d", fs.Domain, aid)
}

func (fs *FreshService) CreateCategory(category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories", fs.Domain)
	result := &categoryResult{}
	err := fs.doPost(url, category, result)
	return result.Category, err
}

func (fs *FreshService) UpdateCategory(cid int64, category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fs.Domain, cid)
	result := &categoryResult{}
	err := fs.doPut(url, category, result)
	return result.Category, err
}

func (fs *FreshService) GetCategory(cid int64) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fs.Domain, cid)
	result := &categoryResult{}
	err := fs.doGet(url, result)
	return result.Category, err
}

func (fs *FreshService) ListCategories() ([]*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories", fs.Domain)
	result := &categoryResult{}
	err := fs.doGet(url, result)
	return result.Categories, err
}

func (fs *FreshService) DeleteCategory(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fs.Domain, cid)
	return fs.doDelete(url)
}

func (fs *FreshService) CreateFolder(folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders", fs.Domain)
	result := &folderResult{}
	err := fs.doPost(url, folder, result)
	return result.Foler, err
}

func (fs *FreshService) UpdateFolder(fid int64, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fs.Domain, fid)
	result := &folderResult{}
	err := fs.doPut(url, folder, result)
	return result.Foler, err
}

func (fs *FreshService) GetFolder(fid int64) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fs.Domain, fid)
	result := &folderResult{}
	err := fs.doGet(url, result)
	return result.Foler, err
}

func (fs *FreshService) ListCategoryFolders(cid int64) ([]*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders?category_id=%d", fs.Domain, cid)
	result := &folderResult{}
	err := fs.doGet(url, result)
	return result.Folders, err
}

func (fs *FreshService) DeleteFolder(fid int64) error {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fs.Domain, fid)
	return fs.doDelete(url)
}

func (fs *FreshService) CreateArticle(fid int64, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fs.Domain, fid)
	result := &articleResult{}
	err := fs.doPost(url, article, result)
	return result.Article, err
}

func (fs *FreshService) SendArticleToApproval(aid int64) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d/send_for_approval", fs.Domain, aid)
	result := &articleResult{}
	err := fs.doPut(url, nil, result)
	return result.Article, err
}

func (fs *FreshService) UpdateArticle(aid int64, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fs.Domain, aid)
	result := &articleResult{}
	err := fs.doPut(url, article, result)
	return result.Article, err
}

func (fs *FreshService) GetArticle(aid int64) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fs.Domain, aid)
	result := &articleResult{}
	err := fs.doGet(url, result)
	return result.Article, err
}

func (fs *FreshService) ListFolderArticles(fid int64) ([]*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles?folder_id=%d", fs.Domain, fid)
	result := &articleResult{}
	err := fs.doGet(url, result)
	return result.Articles, err
}

func (fs *FreshService) DeleteArticle(aid int64) error {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fs.Domain, aid)
	return fs.doDelete(url)
}

func (fs *FreshService) SearchArticles(sao *SearchArticlesOption) ([]*Article, bool, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/search", fs.Domain)
	result := &articleResult{}
	next, err := fs.doList(url, sao, result)
	return result.Articles, next, err
}
