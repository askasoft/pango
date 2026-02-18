package xmlrpc

import (
	"context"
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	client := &Client{
		Endpoint: "https://bugzilla.mozilla.org/xmlrpc.cgi",
	}

	result := struct {
		Version string `xmlrpc:"version"`
	}{}

	err := client.Call(context.Background(), "Bugzilla.version", &result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Version == "" {
		t.Fatal("empty result")
	}
	fmt.Printf("Version: %s\n", result.Version)
}
