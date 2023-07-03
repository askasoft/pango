package freshservice

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

func (fs *Freshservice) CreateCategory(category *Category) (*Category, error) {
	url := fs.endpoint("/solutions/categories")
	result := &categoryResult{}
	err := fs.doPost(url, category, result)
	return result.Category, err
}

func (fs *Freshservice) UpdateCategory(cid int64, category *Category) (*Category, error) {
	url := fs.endpoint("/solutions/categories/%d", cid)
	result := &categoryResult{}
	err := fs.doPut(url, category, result)
	return result.Category, err
}

func (fs *Freshservice) GetCategory(cid int64) (*Category, error) {
	url := fs.endpoint("/solutions/categories/%d", cid)
	result := &categoryResult{}
	err := fs.doGet(url, result)
	return result.Category, err
}

func (fs *Freshservice) ListCategories(lco *ListCategoriesOption) ([]*Category, bool, error) {
	url := fs.endpoint("/solutions/categories")
	result := &categoryResult{}
	next, err := fs.doList(url, lco, result)
	return result.Categories, next, err
}

func (fs *Freshservice) IterCategories(lco *ListCategoriesOption, icf func(*Category) error) error {
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
		categories, next, err := fs.ListCategories(lco)
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

func (fs *Freshservice) DeleteCategory(cid int64) error {
	url := fs.endpoint("/solutions/categories/%d", cid)
	return fs.doDelete(url)
}

func (fs *Freshservice) CreateFolder(folder *Folder) (*Folder, error) {
	url := fs.endpoint("/solutions/folders")
	result := &folderResult{}
	err := fs.doPost(url, folder, result)
	return result.Foler, err
}

func (fs *Freshservice) UpdateFolder(fid int64, folder *Folder) (*Folder, error) {
	url := fs.endpoint("/solutions/folders/%d", fid)
	result := &folderResult{}
	err := fs.doPut(url, folder, result)
	return result.Foler, err
}

func (fs *Freshservice) GetFolder(fid int64) (*Folder, error) {
	url := fs.endpoint("/solutions/folders/%d", fid)
	result := &folderResult{}
	err := fs.doGet(url, result)
	return result.Foler, err
}

func (fs *Freshservice) ListCategoryFolders(cid int64, lfo *ListFoldersOption) ([]*Folder, bool, error) {
	if lfo == nil {
		lfo = &ListFoldersOption{}
	}
	lfo.categoryID = cid

	url := fs.endpoint("/solutions/folders")
	result := &folderResult{}
	next, err := fs.doList(url, lfo, result)
	return result.Folders, next, err
}

func (fs *Freshservice) IterCategoryFolders(cid int64, lfo *ListFoldersOption, iff func(*Folder) error) error {
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
		folders, next, err := fs.ListCategoryFolders(cid, lfo)
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

func (fs *Freshservice) DeleteFolder(fid int64) error {
	url := fs.endpoint("/solutions/folders/%d", fid)
	return fs.doDelete(url)
}

func (fs *Freshservice) CreateArticle(article *Article) (*Article, error) {
	url := fs.endpoint("/solutions/articles")
	result := &articleResult{}
	err := fs.doPost(url, article, result)
	return result.Article, err
}

func (fs *Freshservice) SendArticleToApproval(aid int64) (*Article, error) {
	url := fs.endpoint("/solutions/articles/%d/send_for_approval", aid)
	result := &articleResult{}
	err := fs.doPut(url, nil, result)
	return result.Article, err
}

func (fs *Freshservice) UpdateArticle(aid int64, article *Article) (*Article, error) {
	url := fs.endpoint("/solutions/articles/%d", aid)
	result := &articleResult{}
	err := fs.doPut(url, article, result)
	return result.Article, err
}

func (fs *Freshservice) GetArticle(aid int64) (*Article, error) {
	url := fs.endpoint("/solutions/articles/%d", aid)
	result := &articleResult{}
	err := fs.doGet(url, result)
	return result.Article, err
}

func (fs *Freshservice) ListFolderArticles(fid int64, lao *ListArticlesOption) ([]*Article, bool, error) {
	if lao == nil {
		lao = &ListArticlesOption{}
	}
	lao.folderID = fid

	url := fs.endpoint("/solutions/articles")
	result := &articleResult{}
	next, err := fs.doList(url, lao, result)
	return result.Articles, next, err
}

func (fs *Freshservice) IterFolderArticles(fid int64, lao *ListArticlesOption, iaf func(*Article) error) error {
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
		articles, next, err := fs.ListFolderArticles(fid, lao)
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

func (fs *Freshservice) DeleteArticle(aid int64) error {
	url := fs.endpoint("/solutions/articles/%d", aid)
	return fs.doDelete(url)
}

func (fs *Freshservice) SearchArticles(sao *SearchArticlesOption) ([]*Article, bool, error) {
	url := fs.endpoint("/solutions/articles/search")
	result := &articleResult{}
	next, err := fs.doList(url, sao, result)
	return result.Articles, next, err
}

// func (fs *Freshservice) DeleteArticleAttachment(aid, tid int64) error {
// 	url := fs.endpoint("/solutions/articles/%d/attachments/%d", aid, tid)
// 	return fs.doDelete(url)
// }
