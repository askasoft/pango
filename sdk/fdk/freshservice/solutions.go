package freshservice

import "context"

// ---------------------------------------------------
// Solutions

// PerPage: 1 ~ 100, default: 30
type ListCategoriesOption = PageOption

// PerPage: 1 ~ 100, default: 30
type ListFoldersOption struct {
	categoryID int64
	Page       int
	PerPage    int
}

func (lfo *ListFoldersOption) IsNil() bool {
	return lfo == nil
}

func (lfo *ListFoldersOption) Values() Values {
	q := Values{}
	q.SetInt64("category_id", lfo.categoryID)
	q.SetInt("page", lfo.Page)
	q.SetInt("per_page", lfo.PerPage)
	return q
}

// PerPage: 1 ~ 100, default: 30
type ListArticlesOption struct {
	folderID int64
	Page     int
	PerPage  int
}

func (lao *ListArticlesOption) IsNil() bool {
	return lao == nil
}

func (lao *ListArticlesOption) Values() Values {
	q := Values{}
	q.SetInt64("folder_id", lao.folderID)
	q.SetInt("page", lao.Page)
	q.SetInt("per_page", lao.PerPage)
	return q
}

type SearchArticlesOption struct {
	SearchTerm string // The keywords for which the solution articles have to be searched.
	UserEmail  string // By default, the API will search the articles for the user whose API key is provided. If you want to search articles for a different user, please provide their user_email.
	Page       int
	PerPage    int
}

func (sao *SearchArticlesOption) IsNil() bool {
	return sao == nil
}

func (sao *SearchArticlesOption) Values() Values {
	q := Values{}
	q.SetString("search_term", sao.SearchTerm)
	q.SetString("user_email", sao.UserEmail)
	q.SetInt("page", sao.Page)
	q.SetInt("per_page", sao.PerPage)
	return q
}

func (fs *Freshservice) CreateCategory(ctx context.Context, category *CategoryCreate) (*Category, error) {
	url := fs.Endpoint("/solutions/categories")
	result := &categoryResult{}
	if err := fs.DoPost(ctx, url, category, result); err != nil {
		return nil, err
	}
	return result.Category, nil
}

func (fs *Freshservice) UpdateCategory(ctx context.Context, cid int64, category *CategoryUpdate) (*Category, error) {
	url := fs.Endpoint("/solutions/categories/%d", cid)
	result := &categoryResult{}
	if err := fs.DoPut(ctx, url, category, result); err != nil {
		return nil, err
	}
	return result.Category, nil
}

func (fs *Freshservice) GetCategory(ctx context.Context, cid int64) (*Category, error) {
	url := fs.Endpoint("/solutions/categories/%d", cid)
	result := &categoryResult{}
	err := fs.DoGet(ctx, url, result)
	return result.Category, err
}

func (fs *Freshservice) ListCategories(ctx context.Context, lco *ListCategoriesOption) ([]*Category, bool, error) {
	url := fs.Endpoint("/solutions/categories")
	result := &categoriesResult{}
	next, err := fs.DoList(ctx, url, lco, result)
	return result.Categories, next, err
}

func (fs *Freshservice) IterCategories(ctx context.Context, lco *ListCategoriesOption, icf func(*Category) error) error {
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
		categories, next, err := fs.ListCategories(ctx, lco)
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

func (fs *Freshservice) DeleteCategory(ctx context.Context, cid int64) error {
	url := fs.Endpoint("/solutions/categories/%d", cid)
	return fs.DoDelete(ctx, url)
}

func (fs *Freshservice) CreateFolder(ctx context.Context, folder *FolderCreate) (*Folder, error) {
	url := fs.Endpoint("/solutions/folders")
	result := &folderResult{}
	if err := fs.DoPost(ctx, url, folder, result); err != nil {
		return nil, err
	}
	return result.Foler, nil
}

func (fs *Freshservice) UpdateFolder(ctx context.Context, fid int64, folder *FolderUpdate) (*Folder, error) {
	url := fs.Endpoint("/solutions/folders/%d", fid)
	result := &folderResult{}
	if err := fs.DoPut(ctx, url, folder, result); err != nil {
		return nil, err
	}
	return result.Foler, nil
}

func (fs *Freshservice) GetFolder(ctx context.Context, fid int64) (*Folder, error) {
	url := fs.Endpoint("/solutions/folders/%d", fid)
	result := &folderResult{}
	err := fs.DoGet(ctx, url, result)
	return result.Foler, err
}

func (fs *Freshservice) ListCategoryFolders(ctx context.Context, cid int64, lfo *ListFoldersOption) ([]*Folder, bool, error) {
	if lfo == nil {
		lfo = &ListFoldersOption{}
	}
	lfo.categoryID = cid

	url := fs.Endpoint("/solutions/folders")
	result := &foldersResult{}
	next, err := fs.DoList(ctx, url, lfo, result)
	return result.Folders, next, err
}

func (fs *Freshservice) IterCategoryFolders(ctx context.Context, cid int64, lfo *ListFoldersOption, iff func(*Folder) error) error {
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
		folders, next, err := fs.ListCategoryFolders(ctx, cid, lfo)
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

func (fs *Freshservice) DeleteFolder(ctx context.Context, fid int64) error {
	url := fs.Endpoint("/solutions/folders/%d", fid)
	return fs.DoDelete(ctx, url)
}

func (fs *Freshservice) CreateArticle(ctx context.Context, article *ArticleCreate) (*Article, error) {
	url := fs.Endpoint("/solutions/articles")
	result := &articleResult{}
	if err := fs.DoPost(ctx, url, article, result); err != nil {
		return nil, err
	}
	return result.Article, nil
}

func (fs *Freshservice) SendArticleToApproval(ctx context.Context, aid int64) (*Article, error) {
	url := fs.Endpoint("/solutions/articles/%d/send_for_approval", aid)
	result := &articleResult{}
	if err := fs.DoPut(ctx, url, nil, result); err != nil {
		return nil, err
	}
	return result.Article, nil
}

func (fs *Freshservice) UpdateArticle(ctx context.Context, aid int64, article *ArticleUpdate) (*Article, error) {
	url := fs.Endpoint("/solutions/articles/%d", aid)
	result := &articleResult{}
	if err := fs.DoPut(ctx, url, article, result); err != nil {
		return nil, err
	}
	return result.Article, nil
}

func (fs *Freshservice) GetArticle(ctx context.Context, aid int64) (*Article, error) {
	url := fs.Endpoint("/solutions/articles/%d", aid)
	result := &articleResult{}
	err := fs.DoGet(ctx, url, result)
	return result.Article, err
}

func (fs *Freshservice) ListFolderArticles(ctx context.Context, fid int64, lao *ListArticlesOption) ([]*ArticleInfo, bool, error) {
	if lao == nil {
		lao = &ListArticlesOption{}
	}
	lao.folderID = fid

	url := fs.Endpoint("/solutions/articles")
	result := &articlesResult{}
	next, err := fs.DoList(ctx, url, lao, result)
	for _, ai := range result.Articles {
		ai.normalize()
	}
	return result.Articles, next, err
}

func (fs *Freshservice) IterFolderArticles(ctx context.Context, fid int64, lao *ListArticlesOption, iaf func(*ArticleInfo) error) error {
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
		articles, next, err := fs.ListFolderArticles(ctx, fid, lao)
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

func (fs *Freshservice) DeleteArticle(ctx context.Context, aid int64) error {
	url := fs.Endpoint("/solutions/articles/%d", aid)
	return fs.DoDelete(ctx, url)
}

func (fs *Freshservice) SearchArticles(ctx context.Context, sao *SearchArticlesOption) ([]*ArticleInfo, bool, error) {
	url := fs.Endpoint("/solutions/articles/search")
	result := &articlesResult{}
	next, err := fs.DoList(ctx, url, sao, result)
	for _, ai := range result.Articles {
		ai.normalize()
	}
	return result.Articles, next, err
}

// func (fs *Freshservice) DeleteArticleAttachment(aid, tid int64) error {
// 	url := fs.Endpoint("/solutions/articles/%d/attachments/%d", aid, tid)
// 	return fs.DoDelete(ctx, url)
// }
