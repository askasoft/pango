package xmw

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/askasoft/pango/xin"
)

func ExampleLocalizer() {
	router := xin.Default()

	router.Use(NewLocalizer("en", "ja", "zh").Handler())
	router.GET("/", func(c *xin.Context) {
		c.String(200, c.Locale)
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
	req.Header.Add("Accept-Languages", "ja;zh")

	client := &http.Client{Timeout: time.Second * 1}
	res, _ := client.Do(req)

	raw, _ := io.ReadAll(res.Body)
	fmt.Println(string(raw))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
