# ginhtml

support for gin html template

## Install
```bash
go get -u github.com/pandafw/ginx
```

### Example

```go
package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/pandafw/pango/ginx/ginhtml"
)

func main() {
	// new template engine
	ghe := ginhtml.NewEngine()
	if err = ghe.Load("templates"); err != nil {
		panic(err)
	}

	router := gin.Default()

	// customize gin html render
	router.HTMLRender = ghe

	router.GET("/", func(ctx *gin.Context) {
		// render
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
