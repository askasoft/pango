fswatch
====================

recursive directory notifications built as a wrapper around fsnotify (golang)

This is a wrapper around https://github.com/fsnotify/fsnotify instead of only monitoring a top level folder,
it allows you to monitor all folders underneath the folder you specify.

Example:
--------
(error handling omitted to improve readability)
```golang
	import "github.com/panafw/pango/iox/fswatch"

	// works exactly like fsnotify and implements the same API.
	watcher, err := fswatch.NewFileWatcher()

	// watch recursive and recieve events with callback function
	watcher.AddRecursive("watchdir", fswatch.OpALL, "", func(path string, op fswatch.Op) {
		fmt.Printf("%s %s\n", path, op)
	})

```
