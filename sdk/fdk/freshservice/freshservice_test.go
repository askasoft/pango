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

func testNewFreshservice(t *testing.T) *Freshservice {
	apikey := os.Getenv("FSV_APIKEY")
	if apikey == "" {
		t.Skip("FSV_APIKEY not set")
		return nil
	}

	domain := os.Getenv("FSV_DOMAIN")
	if domain == "" {
		t.Skip("FSV_DOMAIN not set")
		return nil
	}

	logs := log.NewLog()
	//logs.SetLevel(log.LevelDebug)
	fd := &Freshservice{
		Domain:             domain,
		Apikey:             apikey,
		Logger:             logs.GetLogger("FSV"),
		RetryOnRateLimited: 1,
	}

	return fd
}
