package freshservice

import (
	"fmt"
	"testing"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/iox/fsu"
	"github.com/askasoft/pango/str"
)

var (
	_ WithFiles = &ArticleCreate{}
	_ WithFiles = &ArticleUpdate{}
)

func TestSolutionAPIs(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	cc := &CategoryCreate{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fs.CreateCategory(ctxbg, cc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fs.DeleteCategory(ctxbg, cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	fc := &FolderCreate{
		CategoryID:  cat.ID,
		Name:        "Test Folder",
		Description: "Test Folder For API Test",
		Visibility:  FolderVisibilityAgents,
	}
	fol, err := fs.CreateFolder(ctxbg, fc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fs.DeleteFolder(ctxbg, fol.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	ac := &ArticleCreate{
		FolderID:    fol.ID,
		Title:       "Test Article",
		Description: "Test Article for API Test",
		Status:      ArticleStatusDraft,
		Tags:        &[]string{"リンゴ"},
	}
	art, err := fs.CreateArticle(ctxbg, ac)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fs.DeleteArticle(ctxbg, art.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	aad, err := fsu.ReadFile("./agent.go")
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	au := &ArticleUpdate{}
	au.AddAttachment("./agent.go")
	ua, err := fs.UpdateArticle(ctxbg, art.ID, au)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	uad, err := fs.DownloadNoAuth(ctxbg, ua.Attachments[0].AttachmentURL)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if !bye.Equal(aad, uad) {
		t.Fatal("Attachment content not equal")
	}

	uad, err = fs.DownloadAttachment(ctxbg, ua.Attachments[0].ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if !bye.Equal(aad, uad) {
		t.Fatal("Attachment content not equal")
	}

	au2 := &ArticleUpdate{}
	au2.AddAttachment("./article.go")
	_, err = fs.UpdateArticle(ctxbg, art.ID, au2)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	au3 := &ArticleUpdate{Tags: &[]string{}}
	au3.AddAttachment("./agent.go", []byte("agent.go"))
	_, err = fs.UpdateArticle(ctxbg, art.ID, au3)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	ga, err := fs.GetArticle(ctxbg, art.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fmt.Println(ga.Tags)
	fmt.Println(ga.Attachments)

	// no article attachment delete method
	// for _, at := range ga.Attachments {
	// 	err = fs.DeleteArticleAttachment(ga.ID, at.ID)
	// 	if err != nil {
	// 		t.Fatalf("ERROR: %v", err)
	// 	}
	// }

	cats, _, err := fs.ListCategories(ctxbg, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(cats) < 1 {
		t.Fatalf("ERROR: categories=%d", len(cats))
	}

	fols, _, err := fs.ListCategoryFolders(ctxbg, cat.ID, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(fols) != 1 {
		t.Fatalf("ERROR: folders=%d", len(fols))
	}

	arts, _, err := fs.ListFolderArticles(ctxbg, fol.ID, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(arts) != 1 {
		t.Fatalf("ERROR: articles=%d", len(arts))
	}

	err = fs.IterFolderArticles(ctxbg, fol.ID, nil, func(ai *ArticleInfo) error {
		fs.Logger.Debug(ai.String())
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}

func TestSolutionIterAllArticles(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	err := fs.IterCategories(ctxbg, nil, func(c *Category) error {
		fs.Logger.Debugf("Enter Category #%d - %s", c.ID, c.Name)

		return fs.IterCategoryFolders(ctxbg, c.ID, nil, func(f *Folder) error {
			fs.Logger.Debugf("Enter Folder #%d - %s", f.ID, f)

			return fs.IterFolderArticles(ctxbg, f.ID, nil, func(ai *ArticleInfo) error {
				fs.Logger.Debugf("Article #%d - %s", ai.ID, ai)
				return nil
			})
		})
	})
	if err != nil {
		t.Error(err)
	}
}

func TestSolutionManyCategories(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	for i := 1; i <= 101; i++ {
		cc := &CategoryCreate{
			Name:        fmt.Sprintf("Test Category %d", i),
			Description: fmt.Sprintf("Test Category For API Test %d", i),
		}
		_, err := fs.CreateCategory(ctxbg, cc)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}

	cids := make([]int64, 0, 101)
	err := fs.IterCategories(ctxbg, nil, func(c *Category) error {
		if str.StartsWith(c.Name, "Test Category") {
			cids = append(cids, c.ID)
		}
		return nil
	})
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	if len(cids) != 101 {
		t.Errorf("ERROR: categories=%d", len(cids))
	}

	for _, cid := range cids {
		fs.DeleteCategory(ctxbg, cid)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}
}

func TestSolutionManyFolders(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	cc := &CategoryCreate{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fs.CreateCategory(ctxbg, cc)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	defer func() {
		err = fs.DeleteCategory(ctxbg, cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	for i := 1; i <= 101; i++ {
		fc := &FolderCreate{
			CategoryID:  cat.ID,
			Name:        fmt.Sprintf("Test Folder %d", i),
			Description: fmt.Sprintf("Test Folder For API Test %d", i),
			Visibility:  FolderVisibilityAgents,
		}
		_, err := fs.CreateFolder(ctxbg, fc)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}

	fids := make([]int64, 0, 101)
	err = fs.IterCategoryFolders(ctxbg, cat.ID, nil, func(f *Folder) error {
		fids = append(fids, f.ID)
		return nil
	})
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	if len(fids) != 101 {
		t.Errorf("ERROR: articles=%d", len(fids))
	}

	for _, fid := range fids {
		fs.DeleteFolder(ctxbg, fid)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}
}

func TestSolutionManyArticles(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	cc := &CategoryCreate{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fs.CreateCategory(ctxbg, cc)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	defer func() {
		err = fs.DeleteCategory(ctxbg, cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	fc := &FolderCreate{
		CategoryID:  cat.ID,
		Name:        "Test Folder",
		Description: "Test Folder For API Test",
		Visibility:  FolderVisibilityAgents,
	}
	fol, err := fs.CreateFolder(ctxbg, fc)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	defer func() {
		err = fs.DeleteFolder(ctxbg, fol.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	for i := 1; i <= 101; i++ {
		ac := &ArticleCreate{
			FolderID:    fol.ID,
			Title:       fmt.Sprintf("Test Article %d", i),
			Description: fmt.Sprintf("Test Article for API Test %d", i),
			Status:      ArticleStatusDraft,
		}

		_, err := fs.CreateArticle(ctxbg, ac)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}

	aids := make([]int64, 0, 101)
	err = fs.IterFolderArticles(ctxbg, fol.ID, nil, func(a *ArticleInfo) error {
		aids = append(aids, a.ID)
		return nil
	})
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	if len(aids) != 101 {
		t.Errorf("ERROR: articles=%d", len(aids))
	}

	for _, aid := range aids {
		fs.DeleteArticle(ctxbg, aid)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}
}
