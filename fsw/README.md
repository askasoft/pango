 Pango FSWatch
=====================================================================


This is a wrapper around https://github.com/fsnotify/fsnotify instead of only monitoring a top level folder,
it allows you to monitor all folders underneath the folder you specify.


### Example:

(error handling omitted to improve readability)

```golang
import (
	"fmt"
	github.com/askasoft/pango/fsw
)

// works exactly like fsnotify and implements the same API.
watcher, err := fsw.NewFileWatcher()

// watch recursive and recieve events with callback function
watcher.AddRecursive("watchdir", fsw.OpALL, func(path string, op fsw.Op) {
	fmt.Printf("%s %s\n", path, op)
})

// start watch go-routine
watcher.Start()
```

