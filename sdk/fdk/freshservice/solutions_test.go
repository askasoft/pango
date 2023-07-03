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
		Visibility:  FolderVisibilityAgentsOnly,
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

	art.AddAttachment("./any.go")
	ua, err := fs.UpdateArticle(art.ID, art)
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

	art = &Article{}
	art.AddAttachment("./note.go")
	ua2, err := fs.UpdateArticle(ua.ID, art)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fmt.Print(ua2.Attachments)

	// no article attachment delete method
	// for _, at := range ua.Attachments {
	// 	err = fs.DeleteArticleAttachment(ua.ID, at.ID)
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

	arts, _, err := fs.ListFolderArticles(fol.ID, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(arts) != 1 {
		t.Fatalf("ERROR: articles=%d", len(arts))
	}

	fols, _, err := fs.ListCategoryFolders(cat.ID, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(fols) != 1 {
		t.Fatalf("ERROR: folders=%d", len(fols))
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
			t.Fatalf("ERROR: %v", err)
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
		t.Fatalf("ERROR: %v", err)
	}
	if len(cids) != 101 {
		t.Fatalf("ERROR: categories=%d", len(cids))
	}

	for _, cid := range cids {
		fs.DeleteCategory(cid)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
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
		t.Fatalf("ERROR: %v", err)
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
			Visibility:  FolderVisibilityAgentsOnly,
		}
		_, err := fs.CreateFolder(cf)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
	}

	fids := make([]int64, 0, 101)
	err = fs.IterCategoryFolders(cat.ID, nil, func(f *Folder) error {
		fids = append(fids, f.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(fids) != 101 {
		t.Fatalf("ERROR: articles=%d", len(fids))
	}

	for _, fid := range fids {
		fs.DeleteFolder(fid)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
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
		Visibility:  FolderVisibilityAgentsOnly,
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

	for i := 1; i <= 101; i++ {
		ca := &Article{
			FolderID:    fol.ID,
			Title:       fmt.Sprintf("Test Article %d", i),
			Description: fmt.Sprintf("Test Article for API Test %d", i),
			Status:      ArticleStatusDraft,
		}

		_, err := fs.CreateArticle(ca)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
	}

	aids := make([]int64, 0, 101)
	err = fs.IterFolderArticles(fol.ID, nil, func(a *Article) error {
		aids = append(aids, a.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(aids) != 101 {
		t.Fatalf("ERROR: articles=%d", len(aids))
	}

	for _, aid := range aids {
		fs.DeleteArticle(aid)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
	}
}
