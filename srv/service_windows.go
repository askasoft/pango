package srv

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

var elog debug.Log

type service struct {
	app App
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}
	s.app.Init()
	s.app.Run()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		c := <-r
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			changes <- svc.Status{State: svc.StopPending}
			s.app.Shutdown()
			break loop
		default:
			continue loop
		}
	}
	return
}

func runService(app App, name string, isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name)) //nolint: errcheck
	run := svc.Run
	if isDebug {
		run = debug.Run
	}

	srv := &service{app}
	err = run(name, srv)
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err)) //nolint: errcheck
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name)) //nolint: errcheck
}

func exePath() (path string, err error) {
	prog := os.Args[0]

	if path, err = filepath.Abs(prog); err != nil {
		return
	}

	fi, err := os.Stat(path)
	if err == nil {
		if !fi.Mode().IsDir() {
			return
		}
	}

	if filepath.Ext(path) == "" {
		path += ".exe"
		fi, err = os.Stat(path)
		if err == nil {
			if !fi.Mode().IsDir() {
				return
			}
		}
	}

	err = fmt.Errorf("%s is not a executable file", prog)
	return
}

func installService(name, display, description string) error {
	exepath, err := exePath()
	if err != nil {
		return err
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect() //nolint: errcheck

	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", name)
	}

	cfg := mgr.Config{
		DisplayName: display,
		Description: description,
		StartType:   mgr.StartAutomatic,
	}
	s, err = m.CreateService(name, exepath, cfg, "-d", filepath.Dir(exepath))
	if err != nil {
		return err
	}
	defer s.Close()

	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete() //nolint: errcheck
		return fmt.Errorf("SetupEventLogSource() failed: %w", err)
	}

	return nil
}

func removeService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect() //nolint: errcheck

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", name)
	}
	defer s.Close()

	err = s.Delete()
	if err != nil {
		return err
	}

	err = eventlog.Remove(name)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %w", err)
	}
	return nil
}

func startService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect() //nolint: errcheck

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %w", err)
	}
	defer s.Close()

	err = s.Start("is", "manual-started")
	if err != nil {
		return fmt.Errorf("could not start service: %w", err)
	}
	return nil
}

func controlService(name string, c svc.Cmd, to svc.State) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect() //nolint: errcheck

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %w", err)
	}
	defer s.Close()

	status, err := s.Control(c)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %w", c, err)
	}

	timeout := time.Now().Add(10 * time.Second)
	for status.State != to {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", to)
		}

		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %w", err)
		}
	}
	return nil
}
