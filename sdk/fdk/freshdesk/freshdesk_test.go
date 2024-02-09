package freshdesk

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
)

func TestWithFiles(t *testing.T) {
	var (
		tt WithFiles = &Ticket{}
		tc WithFiles = &Conversation{}
		at WithFiles = &Article{}
		ac WithFiles = &Contact{}
		//ag WithFiles = &Agent{}
	)
	fmt.Println(tt, tc, at, ac)
}

func testNewFreshdesk(t *testing.T) *Freshdesk {
	apikey := os.Getenv("FDK_APIKEY")
	if apikey == "" {
		t.Skip("FDK_APIKEY not set")
		return nil
	}

	domain := os.Getenv("FDK_DOMAIN")
	if domain == "" {
		t.Skip("FDK_DOMAIN not set")
		return nil
	}

	logs := log.NewLog()
	//logs.SetLevel(log.LevelDebug)
	fd := &Freshdesk{
		Domain:     domain,
		Apikey:     apikey,
		Logger:     logs.GetLogger("FDK"),
		MaxRetries: 1,
		RetryAfter: time.Second * 3,
	}

	return fd
}
