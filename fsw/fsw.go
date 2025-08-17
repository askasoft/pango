package fsw

//--------------------------------------------------------------------
// package functions
//

// default watcher instance
var _fsw = NewFileWatcher()

// Default returns the default FileWatcher instance used by the package-level functions.
func Default() *FileWatcher {
	return _fsw
}

// SetDefault set the default FileWatcher instance used by the package-level functions.
func SetDefault(fw *FileWatcher) {
	_fsw = fw
}

// Start start file watching go-routine
func Start() error {
	return _fsw.Start()
}

// Stop stop file watching go-routine
func Stop() error {
	return _fsw.Stop()
}

// Add add a file to watch on specified operation op occurred
func Add(path string, op Op, callback func(string, Op)) error {
	return _fsw.Add(path, op, callback)
}

// Remove stop watching the file
func Remove(path string) error {
	return _fsw.Remove(path)
}

// AddRecursive add files and all sub-directories under the path to watch
func AddRecursive(path string, op Op, cb func(string, Op)) error {
	return _fsw.AddRecursive(path, op, cb)
}

// RemoveRecursive stops watching the directory and all sub-directories.
func RemoveRecursive(path string) error {
	return _fsw.RemoveRecursive(path)
}
