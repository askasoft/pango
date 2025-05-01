package freshdesk

import (
	"os"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
)

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
	logs.SetLevel(log.LevelInfo)
	fd := &Freshdesk{
		Domain:     domain,
		Apikey:     apikey,
		Logger:     logs.GetLogger("FDK"),
		MaxRetries: 1,
		RetryAfter: time.Second * 3,
	}

	return fd
}
