package freshdesk

import (
	"net/url"
)

// ---------------------------------------------------
// Solutions

// PerPage: 1 ~ 100, default: 30
type ListCategoriesOption = PageOption

// PerPage: 1 ~ 100, default: 30
type ListFoldersOption = PageOption

// PerPage: 1 ~ 100, default: 30
type ListArticlesOption = PageOption

func (fd *Freshdesk) CreateCategory(category *Category) (*Category, error) {
	url := fd.endpoint("/solutions/categories")
	result := &Category{}
	if err := fd.doPost(url, category, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) CreateCategoryTranslated(cid int64, lang string, category *Category) (*Category, error) {
	url := fd.endpoint("/solutions/categories/%d/%s", cid, lang)
	result := &Category{}
	if err := fd.doPost(url, category, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateCategory(cid int64, category *Category) (*Category, error) {
	url := fd.endpoint("/solutions/categories/%d", cid)
	result := &Category{}
	if err := fd.doPut(url, category, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateCategoryTranslated(cid int64, lang string, category *Category) (*Category, error) {
	url := fd.endpoint("/solutions/categories/%d/%s", cid, lang)
	result := &Category{}
	if err := fd.doPut(url, category, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) GetCategory(cid int64) (*Category, error) {
	url := fd.endpoint("/solutions/categories/%d", cid)
	cat := &Category{}
	err := fd.doGet(url, cat)
	return cat, err
}

func (fd *Freshdesk) GetCategoryTranslated(cid int64, lang string) (*Category, error) {
	url := fd.endpoint("/solutions/categories/%d/%s", cid, lang)
	cat := &Category{}
	err := fd.doGet(url, cat)
	return cat, err
}

func (fd *Freshdesk) ListCategories(lco *ListCategoriesOption) ([]*Category, bool, error) {
	url := fd.endpoint("/solutions/categories")
	categories := []*Category{}
	next, err := fd.doList(url, lco, &categories)
	return categories, next, err
}

func (fd *Freshdesk) IterCategories(lco *ListCategoriesOption, icf func(*Category) error) error {
	if lco == nil {
		lco = &ListCategoriesOption{}
	}
	if lco.Page < 1 {
		lco.Page = 1
	}
	if lco.PerPage < 1 {
		lco.PerPage = 100
	}

	for {
		categories, next, err := fd.ListCategories(lco)
		if err != nil {
			return err
		}
		for _, c := range categories {
			if err = icf(c); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lco.Page++
	}
	return nil
}

func (fd *Freshdesk) ListCategoriesTranslated(lang string, lco *ListCategoriesOption) ([]*Category, bool, error) {
	url := fd.Domain + "/api/v2/solutions/categories/" + lang
	categories := []*Category{}
	next, err := fd.doList(url, lco, &categories)
	return categories, next, err
}

func (fd *Freshdesk) IterCategoriesTranslated(lang string, lco *ListCategoriesOption, icf func(*Category) error) error {
	if lco == nil {
		lco = &ListCategoriesOption{}
	}
	if lco.Page < 1 {
		lco.Page = 1
	}
	if lco.PerPage < 1 {
		lco.PerPage = 100
	}

	for {
		categories, next, err := fd.ListCategoriesTranslated(lang, lco)
		if err != nil {
			return err
		}
		for _, c := range categories {
			if err = icf(c); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lco.Page++
	}
	return nil
}

func (fd *Freshdesk) DeleteCategory(cid int64) error {
	url := fd.endpoint("/solutions/categories/%d", cid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) CreateFolder(cid int64, folder *Folder) (*Folder, error) {
	url := fd.endpoint("/solutions/categories/%d/folders", cid)
	result := &Folder{}
	if err := fd.doPost(url, folder, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) CreateFolderTranslated(fid int64, lang string, folder *Folder) (*Folder, error) {
	url := fd.endpoint("/solutions/folders/%d/%s", fid, lang)
	result := &Folder{}
	if err := fd.doPost(url, folder, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateFolder(fid int64, folder *Folder) (*Folder, error) {
	url := fd.endpoint("/solutions/folders/%d", fid)
	result := &Folder{}
	if err := fd.doPut(url, folder, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateFolderTranslated(fid int64, lang string, folder *Folder) (*Folder, error) {
	url := fd.endpoint("/solutions/folders/%d/%s", fid, lang)
	result := &Folder{}
	if err := fd.doPut(url, folder, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) GetFolder(fid int64) (*Folder, error) {
	url := fd.endpoint("/solutions/folders/%d", fid)
	folder := &Folder{}
	err := fd.doGet(url, folder)
	return folder, err
}

func (fd *Freshdesk) GetFolderTranslated(fid int64, lang string) (*Folder, error) {
	url := fd.endpoint("/solutions/folders/%d/%s", fid, lang)
	folder := &Folder{}
	err := fd.doGet(url, folder)
	return folder, err
}

func (fd *Freshdesk) ListCategoryFolders(cid int64, lfo *ListFoldersOption) ([]*Folder, bool, error) {
	url := fd.endpoint("/solutions/categories/%d/folders", cid)
	folders := []*Folder{}
	next, err := fd.doList(url, lfo, &folders)
	return folders, next, err
}

func (fd *Freshdesk) IterCategoryFolders(cid int64, lfo *ListFoldersOption, iff func(*Folder) error) error {
	if lfo == nil {
		lfo = &ListFoldersOption{}
	}
	if lfo.Page < 1 {
		lfo.Page = 1
	}
	if lfo.PerPage < 1 {
		lfo.PerPage = 100
	}

	for {
		folders, next, err := fd.ListCategoryFolders(cid, lfo)
		if err != nil {
			return err
		}
		for _, f := range folders {
			if err = iff(f); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lfo.Page++
	}
	return nil
}

func (fd *Freshdesk) ListCategoryFoldersTranslated(cid int64, lang string, lfo *ListFoldersOption) ([]*Folder, bool, error) {
	url := fd.endpoint("/solutions/categories/%d/folders/%s", cid, lang)
	folders := []*Folder{}
	next, err := fd.doList(url, lfo, &folders)
	return folders, next, err
}

func (fd *Freshdesk) IterCategoryFoldersTranslated(cid int64, lang string, lfo *ListFoldersOption, iff func(*Folder) error) error {
	if lfo == nil {
		lfo = &ListFoldersOption{}
	}
	if lfo.Page < 1 {
		lfo.Page = 1
	}
	if lfo.PerPage < 1 {
		lfo.PerPage = 100
	}

	for {
		folders, next, err := fd.ListCategoryFoldersTranslated(cid, lang, lfo)
		if err != nil {
			return err
		}
		for _, f := range folders {
			if err = iff(f); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lfo.Page++
	}
	return nil
}

func (fd *Freshdesk) ListSubFolders(fid int64, lfo *ListFoldersOption) ([]*Folder, bool, error) {
	url := fd.endpoint("/solutions/folders/%d/subfolders", fid)
	folders := []*Folder{}
	next, err := fd.doList(url, lfo, &folders)
	return folders, next, err
}

func (fd *Freshdesk) IterSubFolders(fid int64, lfo *ListFoldersOption, iff func(*Folder) error) error {
	if lfo == nil {
		lfo = &ListFoldersOption{}
	}
	if lfo.Page < 1 {
		lfo.Page = 1
	}
	if lfo.PerPage < 1 {
		lfo.PerPage = 100
	}

	for {
		folders, next, err := fd.ListSubFolders(fid, lfo)
		if err != nil {
			return err
		}
		for _, f := range folders {
			if err = iff(f); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lfo.Page++
	}
	return nil
}

func (fd *Freshdesk) ListSubFoldersTranslated(fid int64, lang string, lfo *ListFoldersOption) ([]*Folder, bool, error) {
	url := fd.endpoint("/solutions/folders/%d/subfolders/%s", fid, lang)
	folders := []*Folder{}
	next, err := fd.doList(url, lfo, &folders)
	return folders, next, err
}

func (fd *Freshdesk) IterSubFoldersTranslated(fid int64, lang string, lfo *ListFoldersOption, iff func(*Folder) error) error {
	if lfo == nil {
		lfo = &ListFoldersOption{}
	}
	if lfo.Page < 1 {
		lfo.Page = 1
	}
	if lfo.PerPage < 1 {
		lfo.PerPage = 100
	}

	for {
		folders, next, err := fd.ListSubFoldersTranslated(fid, lang, lfo)
		if err != nil {
			return err
		}
		for _, f := range folders {
			if err = iff(f); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lfo.Page++
	}
	return nil
}

func (fd *Freshdesk) DeleteFolder(fid int64) error {
	url := fd.endpoint("/solutions/folders/%d", fid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) CreateArticle(fid int64, article *Article) (*Article, error) {
	url := fd.endpoint("/solutions/folders/%d/articles", fid)
	result := &Article{}
	if err := fd.doPost(url, article, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) CreateArticleTranslated(aid int64, lang string, article *Article) (*Article, error) {
	url := fd.endpoint("/solutions/articles/%d/%s", aid, lang)
	result := &Article{}
	if err := fd.doPost(url, article, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateArticle(aid int64, article *Article) (*Article, error) {
	url := fd.endpoint("/solutions/articles/%d", aid)
	result := &Article{}
	if err := fd.doPut(url, article, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateArticleTranslated(aid int64, lang string, article *Article) (*Article, error) {
	url := fd.endpoint("/solutions/articles/%d/%s", aid, lang)
	result := &Article{}
	if err := fd.doPut(url, article, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) GetArticle(aid int64) (*Article, error) {
	url := fd.endpoint("/solutions/articles/%d", aid)
	article := &Article{}
	err := fd.doGet(url, article)
	return article, err
}

func (fd *Freshdesk) GetArticleTranslated(aid int64, lang string) (*Article, error) {
	url := fd.endpoint("/solutions/articles/%d/%s", aid, lang)
	article := &Article{}
	err := fd.doGet(url, article)
	return article, err
}

func (fd *Freshdesk) ListFolderArticles(fid int64, lao *ListArticlesOption) ([]*Article, bool, error) {
	url := fd.endpoint("/solutions/folders/%d/articles", fid)
	articles := []*Article{}
	next, err := fd.doList(url, lao, &articles)
	return articles, next, err
}

func (fd *Freshdesk) IterFolderArticles(fid int64, lao *ListArticlesOption, iaf func(*Article) error) error {
	if lao == nil {
		lao = &ListArticlesOption{}
	}
	if lao.Page < 1 {
		lao.Page = 1
	}
	if lao.PerPage < 1 {
		lao.PerPage = 100
	}

	for {
		articles, next, err := fd.ListFolderArticles(fid, lao)
		if err != nil {
			return err
		}
		for _, a := range articles {
			if err = iaf(a); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lao.Page++
	}
	return nil
}

func (fd *Freshdesk) ListFolderArticlesTranslated(fid int64, lang string, lao *ListArticlesOption) ([]*Article, bool, error) {
	url := fd.endpoint("/solutions/folders/%d/farticles/%s", fid, lang)
	articles := []*Article{}
	next, err := fd.doList(url, lao, &articles)
	return articles, next, err
}

func (fd *Freshdesk) IterFolderArticlesTranslated(fid int64, lang string, lao *ListArticlesOption, iaf func(*Article) error) error {
	if lao == nil {
		lao = &ListArticlesOption{}
	}
	if lao.Page < 1 {
		lao.Page = 1
	}
	if lao.PerPage < 1 {
		lao.PerPage = 100
	}

	for {
		articles, next, err := fd.ListFolderArticlesTranslated(fid, lang, lao)
		if err != nil {
			return err
		}
		for _, a := range articles {
			if err = iaf(a); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lao.Page++
	}
	return nil
}

func (fd *Freshdesk) DeleteArticle(aid int64) error {
	url := fd.endpoint("/solutions/articles/%d", aid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) SearchArticles(keyword string) ([]*ArticleEx, error) {
	url := fd.endpoint("/search/solutions?term=%s", url.QueryEscape(keyword))
	articles := []*ArticleEx{}
	err := fd.doGet(url, &articles)
	return articles, err
}
