package freshdesk

import (
	"context"
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

func (fd *Freshdesk) CreateCategory(ctx context.Context, category *CategoryCreate) (*Category, error) {
	url := fd.Endpoint("/solutions/categories")
	result := &Category{}
	if err := fd.DoPost(ctx, url, category, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) CreateCategoryTranslated(ctx context.Context, cid int64, lang string, category *CategoryCreate) (*Category, error) {
	url := fd.Endpoint("/solutions/categories/%d/%s", cid, lang)
	result := &Category{}
	if err := fd.DoPost(ctx, url, category, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateCategory(ctx context.Context, cid int64, category *CategoryUpdate) (*Category, error) {
	url := fd.Endpoint("/solutions/categories/%d", cid)
	result := &Category{}
	if err := fd.DoPut(ctx, url, category, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateCategoryTranslated(ctx context.Context, cid int64, lang string, category *CategoryUpdate) (*Category, error) {
	url := fd.Endpoint("/solutions/categories/%d/%s", cid, lang)
	result := &Category{}
	if err := fd.DoPut(ctx, url, category, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) GetCategory(ctx context.Context, cid int64) (*Category, error) {
	url := fd.Endpoint("/solutions/categories/%d", cid)
	cat := &Category{}
	err := fd.DoGet(ctx, url, cat)
	return cat, err
}

func (fd *Freshdesk) GetCategoryTranslated(ctx context.Context, cid int64, lang string) (*Category, error) {
	url := fd.Endpoint("/solutions/categories/%d/%s", cid, lang)
	cat := &Category{}
	err := fd.DoGet(ctx, url, cat)
	return cat, err
}

func (fd *Freshdesk) ListCategories(ctx context.Context, lco *ListCategoriesOption) ([]*Category, bool, error) {
	url := fd.Endpoint("/solutions/categories")
	categories := []*Category{}
	next, err := fd.DoList(ctx, url, lco, &categories)
	return categories, next, err
}

func (fd *Freshdesk) IterCategories(ctx context.Context, lco *ListCategoriesOption, icf func(*Category) error) error {
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
		categories, next, err := fd.ListCategories(ctx, lco)
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

func (fd *Freshdesk) ListCategoriesTranslated(ctx context.Context, lang string, lco *ListCategoriesOption) ([]*Category, bool, error) {
	url := fd.Domain + "/api/v2/solutions/categories/" + lang
	categories := []*Category{}
	next, err := fd.DoList(ctx, url, lco, &categories)
	return categories, next, err
}

func (fd *Freshdesk) IterCategoriesTranslated(ctx context.Context, lang string, lco *ListCategoriesOption, icf func(*Category) error) error {
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
		categories, next, err := fd.ListCategoriesTranslated(ctx, lang, lco)
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

func (fd *Freshdesk) DeleteCategory(ctx context.Context, cid int64) error {
	url := fd.Endpoint("/solutions/categories/%d", cid)
	return fd.DoDelete(ctx, url)
}

func (fd *Freshdesk) CreateFolder(ctx context.Context, cid int64, folder *FolderCreate) (*Folder, error) {
	url := fd.Endpoint("/solutions/categories/%d/folders", cid)
	result := &Folder{}
	if err := fd.DoPost(ctx, url, folder, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) CreateFolderTranslated(ctx context.Context, fid int64, lang string, folder *FolderCreate) (*Folder, error) {
	url := fd.Endpoint("/solutions/folders/%d/%s", fid, lang)
	result := &Folder{}
	if err := fd.DoPost(ctx, url, folder, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateFolder(ctx context.Context, fid int64, folder *FolderUpdate) (*Folder, error) {
	url := fd.Endpoint("/solutions/folders/%d", fid)
	result := &Folder{}
	if err := fd.DoPut(ctx, url, folder, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateFolderTranslated(ctx context.Context, fid int64, lang string, folder *FolderUpdate) (*Folder, error) {
	url := fd.Endpoint("/solutions/folders/%d/%s", fid, lang)
	result := &Folder{}
	if err := fd.DoPut(ctx, url, folder, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) GetFolder(ctx context.Context, fid int64) (*Folder, error) {
	url := fd.Endpoint("/solutions/folders/%d", fid)
	folder := &Folder{}
	err := fd.DoGet(ctx, url, folder)
	return folder, err
}

func (fd *Freshdesk) GetFolderTranslated(ctx context.Context, fid int64, lang string) (*Folder, error) {
	url := fd.Endpoint("/solutions/folders/%d/%s", fid, lang)
	folder := &Folder{}
	err := fd.DoGet(ctx, url, folder)
	return folder, err
}

func (fd *Freshdesk) ListCategoryFolders(ctx context.Context, cid int64, lfo *ListFoldersOption) ([]*Folder, bool, error) {
	url := fd.Endpoint("/solutions/categories/%d/folders", cid)
	folders := []*Folder{}
	next, err := fd.DoList(ctx, url, lfo, &folders)
	return folders, next, err
}

func (fd *Freshdesk) IterCategoryFolders(ctx context.Context, cid int64, lfo *ListFoldersOption, iff func(*Folder) error) error {
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
		folders, next, err := fd.ListCategoryFolders(ctx, cid, lfo)
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

func (fd *Freshdesk) ListCategoryFoldersTranslated(ctx context.Context, cid int64, lang string, lfo *ListFoldersOption) ([]*Folder, bool, error) {
	url := fd.Endpoint("/solutions/categories/%d/folders/%s", cid, lang)
	folders := []*Folder{}
	next, err := fd.DoList(ctx, url, lfo, &folders)
	return folders, next, err
}

func (fd *Freshdesk) IterCategoryFoldersTranslated(ctx context.Context, cid int64, lang string, lfo *ListFoldersOption, iff func(*Folder) error) error {
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
		folders, next, err := fd.ListCategoryFoldersTranslated(ctx, cid, lang, lfo)
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

func (fd *Freshdesk) ListSubFolders(ctx context.Context, fid int64, lfo *ListFoldersOption) ([]*Folder, bool, error) {
	url := fd.Endpoint("/solutions/folders/%d/subfolders", fid)
	folders := []*Folder{}
	next, err := fd.DoList(ctx, url, lfo, &folders)
	return folders, next, err
}

func (fd *Freshdesk) IterSubFolders(ctx context.Context, fid int64, lfo *ListFoldersOption, iff func(*Folder) error) error {
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
		folders, next, err := fd.ListSubFolders(ctx, fid, lfo)
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

func (fd *Freshdesk) ListSubFoldersTranslated(ctx context.Context, fid int64, lang string, lfo *ListFoldersOption) ([]*Folder, bool, error) {
	url := fd.Endpoint("/solutions/folders/%d/subfolders/%s", fid, lang)
	folders := []*Folder{}
	next, err := fd.DoList(ctx, url, lfo, &folders)
	return folders, next, err
}

func (fd *Freshdesk) IterSubFoldersTranslated(ctx context.Context, fid int64, lang string, lfo *ListFoldersOption, iff func(*Folder) error) error {
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
		folders, next, err := fd.ListSubFoldersTranslated(ctx, fid, lang, lfo)
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

func (fd *Freshdesk) DeleteFolder(ctx context.Context, fid int64) error {
	url := fd.Endpoint("/solutions/folders/%d", fid)
	return fd.DoDelete(ctx, url)
}

func (fd *Freshdesk) CreateArticle(ctx context.Context, fid int64, article *ArticleCreate) (*Article, error) {
	url := fd.Endpoint("/solutions/folders/%d/articles", fid)
	result := &Article{}
	if err := fd.DoPost(ctx, url, article, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) CreateArticleTranslated(ctx context.Context, aid int64, lang string, article *ArticleCreate) (*Article, error) {
	url := fd.Endpoint("/solutions/articles/%d/%s", aid, lang)
	result := &Article{}
	if err := fd.DoPost(ctx, url, article, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateArticle(ctx context.Context, aid int64, article *ArticleUpdate) (*Article, error) {
	url := fd.Endpoint("/solutions/articles/%d", aid)
	result := &Article{}
	if err := fd.DoPut(ctx, url, article, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateArticleTranslated(ctx context.Context, aid int64, lang string, article *ArticleUpdate) (*Article, error) {
	url := fd.Endpoint("/solutions/articles/%d/%s", aid, lang)
	result := &Article{}
	if err := fd.DoPut(ctx, url, article, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) GetArticle(ctx context.Context, aid int64) (*Article, error) {
	url := fd.Endpoint("/solutions/articles/%d", aid)
	article := &Article{}
	err := fd.DoGet(ctx, url, article)
	return article, err
}

func (fd *Freshdesk) GetArticleTranslated(ctx context.Context, aid int64, lang string) (*Article, error) {
	url := fd.Endpoint("/solutions/articles/%d/%s", aid, lang)
	article := &Article{}
	err := fd.DoGet(ctx, url, article)
	return article, err
}

func (fd *Freshdesk) ListFolderArticles(ctx context.Context, fid int64, lao *ListArticlesOption) ([]*Article, bool, error) {
	url := fd.Endpoint("/solutions/folders/%d/articles", fid)
	articles := []*Article{}
	next, err := fd.DoList(ctx, url, lao, &articles)
	return articles, next, err
}

func (fd *Freshdesk) IterFolderArticles(ctx context.Context, fid int64, lao *ListArticlesOption, iaf func(*Article) error) error {
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
		articles, next, err := fd.ListFolderArticles(ctx, fid, lao)
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

func (fd *Freshdesk) ListFolderArticlesTranslated(ctx context.Context, fid int64, lang string, lao *ListArticlesOption) ([]*Article, bool, error) {
	url := fd.Endpoint("/solutions/folders/%d/farticles/%s", fid, lang)
	articles := []*Article{}
	next, err := fd.DoList(ctx, url, lao, &articles)
	return articles, next, err
}

func (fd *Freshdesk) IterFolderArticlesTranslated(ctx context.Context, fid int64, lang string, lao *ListArticlesOption, iaf func(*Article) error) error {
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
		articles, next, err := fd.ListFolderArticlesTranslated(ctx, fid, lang, lao)
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

func (fd *Freshdesk) DeleteArticle(ctx context.Context, aid int64) error {
	url := fd.Endpoint("/solutions/articles/%d", aid)
	return fd.DoDelete(ctx, url)
}

func (fd *Freshdesk) SearchArticles(ctx context.Context, keyword string) ([]*ArticleEx, error) {
	url := fd.Endpoint("/search/solutions?term=%s", url.QueryEscape(keyword))
	articles := []*ArticleEx{}
	err := fd.DoGet(ctx, url, &articles)
	return articles, err
}
