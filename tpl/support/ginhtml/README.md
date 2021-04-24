# GinView

support for gin html template

## Install
```bash
go get -u github.com/pandafw/pango/tpl/support/ginhtml
```

### Example

```go
package main

import (
	"github.com/pandafw/pango/tpl"
	"github.com/pandafw/pango/tpl/supports/ginhtml"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	ht := tpl.NewHTMLTemplate()
	if err = ht.Load("templates"); err != nil {
		panic(err)
	}

	router := gin.Default()

	//new template engine
	router.HTMLRender = ginhtml.NewEngine(ht)

	router.GET("/", func(ctx *gin.Context) {
		//render with master
		ctx.HTML(http.StatusOK, "index", gin.H{
			"title": "Index title!",
			"add": func(a int, b int) int {
				return a + b
			},
		})
	})

	router.Run(":9090")
}
```
