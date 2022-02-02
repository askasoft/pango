package gmw

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/pandafw/pango/gin"
)

func Example() {
	router := gin.New()
	router.Use(NewAccessLogger(os.Stdout, DefaultTextLogFormat).Handler())

	router.Any("/example", func(c *gin.Context) {
		c.String(http.StatusOK, c.Request.URL.String())
	})

	server := &http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: router,
	}

	go func() {
		server.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 100)

	req, _ := http.NewRequest("GET", "http://127.0.0.1:8888/example?a=100", nil)
	client := &http.Client{Timeout: time.Second * 1}
	client.Do(req)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
