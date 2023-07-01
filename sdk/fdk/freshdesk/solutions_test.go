package freshdesk

import (
	"fmt"
	"testing"

	"github.com/askasoft/pango/str"
)

func TestSolutionAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	cc := &Category{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fd.CreateCategory(cc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteCategory(cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	cf := &Folder{
		Name:        "Test Folder",
		Description: "Test Folder For API Test",
		Visibility:  FolderVisibilityAgents,
	}
	fol, err := fd.CreateFolder(cat.ID, cf)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteFolder(fol.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	ca := &Article{
		Title:       "Test Article",
		Description: "Test Article for API Test",
		Status:      ArticleStatusDraft,
	}
	art, err := fd.CreateArticle(fol.ID, ca)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteArticle(art.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	art.AddAttachment("./any.go")
	_, err = fd.UpdateArticle(art.ID, art)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	cats, _, err := fd.ListCategories(nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(cats) < 1 {
		t.Fatalf("ERROR: categories=%d", len(cats))
	}

	arts, _, err := fd.ListFolderArticles(fol.ID, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(arts) != 1 {
		t.Fatalf("ERROR: articles=%d", len(arts))
	}

	fols, _, err := fd.ListCategoryFolders(cat.ID, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(fols) != 1 {
		t.Fatalf("ERROR: folders=%d", len(fols))
	}
}

func TestSolutionManyCategories(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	for i := 1; i <= 101; i++ {
		cc := &Category{
			Name:        fmt.Sprintf("Test Category %d", i),
			Description: fmt.Sprintf("Test Category For API Test %d", i),
		}
		_, err := fd.CreateCategory(cc)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
	}

	cids := make([]int64, 0, 101)
	err := fd.IterCategories(nil, func(c *Category) error {
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
		fd.DeleteCategory(cid)
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

	cc := &Category{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fd.CreateCategory(cc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteCategory(cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	for i := 1; i <= 101; i++ {
		cf := &Folder{
			Name:        fmt.Sprintf("Test Folder %d", i),
			Description: fmt.Sprintf("Test Folder For API Test %d", i),
			Visibility:  FolderVisibilityAgents,
		}
		_, err := fd.CreateFolder(cat.ID, cf)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
	}

	fids := make([]int64, 0, 101)
	err = fd.IterCategoryFolders(cat.ID, nil, func(f *Folder) error {
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
		fd.DeleteFolder(fid)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
	}
}

func TestSolutionManyArticles(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	cc := &Category{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fd.CreateCategory(cc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteCategory(cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	cf := &Folder{
		Name:        "Test Folder",
		Description: "Test Folder For API Test",
		Visibility:  FolderVisibilityAgents,
	}
	fol, err := fd.CreateFolder(cat.ID, cf)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteFolder(fol.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	for i := 1; i <= 101; i++ {
		ca := &Article{
			Title:       fmt.Sprintf("Test Article %d", i),
			Description: fmt.Sprintf("Test Article for API Test %d", i),
			Status:      ArticleStatusDraft,
		}

		_, err := fd.CreateArticle(fol.ID, ca)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
	}

	aids := make([]int64, 0, 101)
	err = fd.IterFolderArticles(fol.ID, nil, func(a *Article) error {
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
		fd.DeleteArticle(aid)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
	}
}
