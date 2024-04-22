package freshservice

import (
	"fmt"
	"testing"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/str"
)

func TestSolutionAPIs(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	cc := &Category{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fs.CreateCategory(cc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fs.DeleteCategory(cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	cf := &Folder{
		CategoryID:  cat.ID,
		Name:        "Test Folder",
		Description: "Test Folder For API Test",
		Visibility:  FolderVisibilityAgents,
	}
	fol, err := fs.CreateFolder(cf)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fs.DeleteFolder(fol.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	ca := &Article{
		FolderID:    fol.ID,
		Title:       "Test Article",
		Description: "Test Article for API Test",
		Status:      ArticleStatusDraft,
	}
	art, err := fs.CreateArticle(ca)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fs.DeleteArticle(art.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	aad, err := fsu.ReadFile("./any.go")
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	ua := &Article{}
	ua.AddAttachment("./any.go")
	ua, err = fs.UpdateArticle(art.ID, ua)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	uad, err := fs.DownloadNoAuth(ua.Attachments[0].AttachmentURL)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if !bye.Equal(aad, uad) {
		t.Fatal("Attachment content not equal")
	}

	uad, err = fs.DownloadAttachment(ua.Attachments[0].ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if !bye.Equal(aad, uad) {
		t.Fatal("Attachment content not equal")
	}

	ua2 := &Article{}
	ua2.AddAttachment("./article.go")
	ua2, err = fs.UpdateArticle(art.ID, ua2)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	ua3 := &Article{}
	ua3.AddAttachment("./any.go", []byte("any.go"))
	ua3, err = fs.UpdateArticle(art.ID, ua3)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	ga, err := fs.GetArticle(art.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fmt.Print(ga.Attachments)

	// no article attachment delete method
	// for _, at := range ga.Attachments {
	// 	err = fs.DeleteArticleAttachment(ga.ID, at.ID)
	// 	if err != nil {
	// 		t.Fatalf("ERROR: %v", err)
	// 	}
	// }

	cats, _, err := fs.ListCategories(nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(cats) < 1 {
		t.Fatalf("ERROR: categories=%d", len(cats))
	}

	fols, _, err := fs.ListCategoryFolders(cat.ID, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(fols) != 1 {
		t.Fatalf("ERROR: folders=%d", len(fols))
	}

	arts, _, err := fs.ListFolderArticles(fol.ID, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(arts) != 1 {
		t.Fatalf("ERROR: articles=%d", len(arts))
	}

	err = fs.IterFolderArticles(fol.ID, nil, func(ai *ArticleInfo) error {
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

	err := fs.IterCategories(nil, func(c *Category) error {
		fs.Logger.Debugf("Enter Category #%d - %s", c.ID, c.Name)

		return fs.IterCategoryFolders(c.ID, nil, func(f *Folder) error {
			fs.Logger.Debugf("Enter Folder #%d - %s", f.ID, f)

			return fs.IterFolderArticles(f.ID, nil, func(ai *ArticleInfo) error {
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
		cc := &Category{
			Name:        fmt.Sprintf("Test Category %d", i),
			Description: fmt.Sprintf("Test Category For API Test %d", i),
		}
		_, err := fs.CreateCategory(cc)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}

	cids := make([]int64, 0, 101)
	err := fs.IterCategories(nil, func(c *Category) error {
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
		fs.DeleteCategory(cid)
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

	cc := &Category{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fs.CreateCategory(cc)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	defer func() {
		err = fs.DeleteCategory(cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	for i := 1; i <= 101; i++ {
		cf := &Folder{
			CategoryID:  cat.ID,
			Name:        fmt.Sprintf("Test Folder %d", i),
			Description: fmt.Sprintf("Test Folder For API Test %d", i),
			Visibility:  FolderVisibilityAgents,
		}
		_, err := fs.CreateFolder(cf)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}

	fids := make([]int64, 0, 101)
	err = fs.IterCategoryFolders(cat.ID, nil, func(f *Folder) error {
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
		fs.DeleteFolder(fid)
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

	cc := &Category{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fs.CreateCategory(cc)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	defer func() {
		err = fs.DeleteCategory(cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	cf := &Folder{
		CategoryID:  cat.ID,
		Name:        "Test Folder",
		Description: "Test Folder For API Test",
		Visibility:  FolderVisibilityAgents,
	}
	fol, err := fs.CreateFolder(cf)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	defer func() {
		err = fs.DeleteFolder(fol.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	for i := 1; i <= 101; i++ {
		ca := &Article{
			FolderID:    fol.ID,
			Title:       fmt.Sprintf("Test Article %d", i),
			Description: fmt.Sprintf("Test Article for API Test %d", i),
			Status:      ArticleStatusDraft,
		}

		_, err := fs.CreateArticle(ca)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}

	aids := make([]int64, 0, 101)
	err = fs.IterFolderArticles(fol.ID, nil, func(a *ArticleInfo) error {
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
		fs.DeleteArticle(aid)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}
}
