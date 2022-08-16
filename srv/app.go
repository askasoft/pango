package srv

import (
	"time"
)

type App interface {
	// Name app/service name
	Name() string

	// DispName app/service display name
	DispName() string

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

	// Run run the app
	Run()

	// Shutdown shutdown the app
	Shutdown()

	// Wait wait for server shutdown
	Wait()
}
