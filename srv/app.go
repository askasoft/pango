package srv

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App interface {
	// Name app/service name
	Name() string

	// DisplayName app/service display name
	DisplayName() string

	// Description app/service description
	Description() string

	// Version app version
	Version() string

	// Revision app revision
	Revision() string

	// BuildTime app build time
	BuildTime() time.Time

	// Init initialize the app
	Init()

	// Reload reload the app
	Reload()

	// Run run the app
	Run()

	// Shutdown shutdown the app
	Shutdown()

	// Wait wait signal for reload or shutdown
	Wait()
}

// Wait wait signal for reload or shutdown the app
func Wait(app App) {
	// signal channel
	sigChan := make(chan os.Signal, 1)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		sig := <-sigChan
		if sig == syscall.SIGHUP {
			app.Reload()
		} else {
			app.Shutdown()
			break
		}
	}
}

func runStandalone(app App) {
	app.Init()

	app.Run()

	app.Wait()
}
