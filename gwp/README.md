 Pango Worker Pool
=====================================================================

Concurrency limiting go-routine pool. Limits the concurrency of task execution, not the number of tasks queued.

This implementation builds on ideas from the following:

- http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang
- http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html


## Example

```go
package main

import (
	"fmt"
	"github.com/askasoft/pango/gwp"
)

func main() {
	wp := gwp.NewWorkerPool(2, 100)
	requests := []string{"alpha", "beta", "gamma", "delta", "epsilon"}

	for _, r := range requests {
		r := r
		wp.Submit(func() {
			fmt.Println("Handling request:", r)
		})
	}

	wp.StopWait()
}
```

