package job

type Job interface {
	IsRunning() bool
	IsAborted() bool
	Abort()
	Start()
}
