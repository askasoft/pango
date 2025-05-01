package freshdesk

import (
	"fmt"
	"testing"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/str"
)

var (
	_ WithFiles = &ArticleCreate{}
	_ WithFiles = &ArticleUpdate{}
)

func TestSolutionAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	cc := &CategoryCreate{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fd.CreateCategory(ctxbg, cc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteCategory(ctxbg, cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	fc := &FolderCreate{
		Name:        "Test Folder",
		Description: "Test Folder For API Test",
		Visibility:  FolderVisibilityAgents,
	}
	fol, err := fd.CreateFolder(ctxbg, cat.ID, fc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteFolder(ctxbg, fol.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	ac := &ArticleCreate{
		Title:       "Test Article",
		Description: "Test Article for API Test",
		Status:      ArticleStatusDraft,
		Tags:        &[]string{"リンゴ"},
	}
	art, err := fd.CreateArticle(ctxbg, fol.ID, ac)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteArticle(ctxbg, art.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	ac2 := &ArticleCreate{
		Title:       "Test Article2",
		Description: "Test Article2 for API Test",
		Status:      ArticleStatusDraft,
		Tags:        &[]string{"りんご"},
	}
	art2, err := fd.CreateArticle(ctxbg, fol.ID, ac2)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteArticle(ctxbg, art2.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	aad, err := fsu.ReadFile("./agent.go")
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	au := &ArticleUpdate{
		Tags: &[]string{"りんご"},
	}
	au.AddAttachment("./agent.go")
	ua, err := fd.UpdateArticle(ctxbg, art.ID, au)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	uad, err := fd.DownloadNoAuth(ctxbg, ua.Attachments[0].AttachmentURL)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if !bye.Equal(aad, uad) {
		t.Fatal("Attachment content not equal")
	}

	au2 := &ArticleUpdate{}
	au2.AddAttachment("./article.go")
	_, err = fd.UpdateArticle(ctxbg, art.ID, au2)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	au3 := &ArticleUpdate{Tags: &[]string{}}
	au3.AddAttachment("./agent.go", []byte("agent.go"))
	_, err = fd.UpdateArticle(ctxbg, art.ID, au3)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	ga, err := fd.GetArticle(ctxbg, art.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	if len(ga.Tags) > 0 {
		t.Fatalf("Tags exists: %v", ga.Tags)
	}
	for _, at := range ga.Attachments {
		err = fd.DeleteAttachment(ctxbg, at.ID)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
	}

	ga2, err := fd.GetArticle(ctxbg, art2.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fmt.Println(ga2.Tags)

	cats, _, err := fd.ListCategories(ctxbg, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(cats) < 1 {
		t.Fatalf("ERROR: categories=%d", len(cats))
	}

	arts, _, err := fd.ListFolderArticles(ctxbg, fol.ID, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(arts) != 2 {
		t.Fatalf("ERROR: articles=%d", len(arts))
	}

	fols, _, err := fd.ListCategoryFolders(ctxbg, cat.ID, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(fols) != 1 {
		t.Fatalf("ERROR: folders=%d", len(fols))
	}
}

func TestSolutionIterAllArticles(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	var itf func(f *Folder) error
	itf = func(f *Folder) error {
		fd.Logger.Debugf("Enter Folder #%d - %s", f.ID, f)

		fd.IterSubFolders(ctxbg, f.ID, nil, itf)

		return fd.IterFolderArticles(ctxbg, f.ID, nil, func(a *Article) error {
			fd.Logger.Debugf("Article #%d - %s", a.ID, a)
			return nil
		})
	}

	err := fd.IterCategories(ctxbg, nil, func(c *Category) error {
		fd.Logger.Debugf("Enter Category #%d - %s", c.ID, c.Name)
		return fd.IterCategoryFolders(ctxbg, c.ID, nil, itf)
	})
	if err != nil {
		t.Error(err)
	}
}

func TestSolutionManyCategories(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	for i := 1; i <= 101; i++ {
		cc := &CategoryCreate{
			Name:        fmt.Sprintf("Test Category %d", i),
			Description: fmt.Sprintf("Test Category For API Test %d", i),
		}
		_, err := fd.CreateCategory(ctxbg, cc)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
	}

	cids := make([]int64, 0, 101)
	err := fd.IterCategories(ctxbg, nil, func(c *Category) error {
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
		fd.DeleteCategory(ctxbg, cid)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
	}
}

func TestSolutionManyFolders(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	cc := &CategoryCreate{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fd.CreateCategory(ctxbg, cc)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteCategory(ctxbg, cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	for i := 1; i <= 101; i++ {
		fc := &FolderCreate{
			Name:        fmt.Sprintf("Test Folder %d", i),
			Description: fmt.Sprintf("Test Folder For API Test %d", i),
			Visibility:  FolderVisibilityAgents,
		}
		_, err := fd.CreateFolder(ctxbg, cat.ID, fc)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}

	fids := make([]int64, 0, 101)
	err = fd.IterCategoryFolders(ctxbg, cat.ID, nil, func(f *Folder) error {
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
		fd.DeleteFolder(ctxbg, fid)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}
}

func TestSolutionManyArticles(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	cc := &CategoryCreate{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fd.CreateCategory(ctxbg, cc)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteCategory(ctxbg, cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	fc := &FolderCreate{
		Name:        "Test Folder",
		Description: "Test Folder For API Test",
		Visibility:  FolderVisibilityAgents,
	}
	fol, err := fd.CreateFolder(ctxbg, cat.ID, fc)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteFolder(ctxbg, fol.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	for i := 1; i <= 101; i++ {
		ac := &ArticleCreate{
			Title:       fmt.Sprintf("Test Article %d", i),
			Description: fmt.Sprintf("Test Article for API Test %d", i),
			Status:      ArticleStatusDraft,
		}

		_, err := fd.CreateArticle(ctxbg, fol.ID, ac)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}

	aids := make([]int64, 0, 101)
	err = fd.IterFolderArticles(ctxbg, fol.ID, nil, func(a *Article) error {
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
		fd.DeleteArticle(ctxbg, aid)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}
}
