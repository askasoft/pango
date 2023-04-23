package freshservice

import (
	"fmt"
	"os"
	"testing"

	"github.com/askasoft/pango/log"
)

func TestWithFiles(t *testing.T) {
	var (
		// tt WithFiles = &Ticket{}
		// tc WithFiles = &Conversation{}
		at WithFiles = &Article{}
		// ac WithFiles = &Contact{}
		//ag WithFiles = &Agent{}
	)
	fmt.Println(at)
}

func testNewFreshservice(t *testing.T) *FreshService {
	apikey := os.Getenv("FSE_APIKEY")
	if apikey == "" {
		t.Skip("FSE_APIKEY not set")
		return nil
	}

	domain := os.Getenv("FSE_DOMAIN")
	if domain == "" {
		t.Skip("FSE_DOMAIN not set")
		return nil
	}

	logs := log.NewLog()
	//logs.SetLevel(log.LevelDebug)
	fd := &FreshService{
		Domain:             domain,
		Apikey:             apikey,
		Logger:             logs.GetLogger("FSE"),
		RetryOnRateLimited: 1,
	}

	return fd
}

func TestSolutionAPIs(t *testing.T) {
	fd := testNewFreshservice(t)
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
		CategoryID:  cat.ID,
		Name:        "Test Folder",
		Description: "Test Folder For API Test",
		Visibility:  FolderVisibilityAgentsOnly,
	}
	fol, err := fd.CreateFolder(cf)
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

	arts, err := fd.ListFolderArticles(fol.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(arts) != 1 {
		t.Fatalf("ERROR: articles=%d", len(arts))
	}

	fols, err := fd.ListCategoryFolders(cat.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(fols) != 1 {
		t.Fatalf("ERROR: folders=%d", len(fols))
	}
}
