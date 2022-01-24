package gmw

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/pandafw/pango/gin"
)

func ExampleHTTPGziper() {
	router := gin.Default()

	router.Use(DefaultHTTPGziper().Handler())
	router.GET("/", func(c *gin.Context) {
		c.String(200, strings.Repeat("This is a Test!\n", 1000))
	})

	server := &http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: router,
	}

	go func() {
		server.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 100)

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8888/", nil)
	client := &http.Client{Timeout: time.Second * 1}
	client.Do(req)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
