package srv

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/sys/windows/svc"
)

// Main server main
func Main(app App) {
	Usage = func() {
		out := flag.CommandLine.Output()

		fmt.Fprintln(out, "Usage: "+app.Name()+".exe <command> [options]")
		fmt.Fprintln(out, "  <command>:")
		if cmd, ok := app.(Cmd); ok {
			cmd.CmdHelp(out)
		}
		fmt.Fprintln(out, "    install         install as windows service.")
		fmt.Fprintln(out, "    remove          remove installed windows service.")
		fmt.Fprintln(out, "    start           start the windows service.")
		fmt.Fprintln(out, "    stop            stop the windows service.")
		fmt.Fprintln(out, "    version         print the version information.")
		fmt.Fprintln(out, "    help | usage    print the usage information.")
		fmt.Fprintln(out, "  <options>:")

		flag.PrintDefaults()
	}

	flag.CommandLine.Usage = Usage

	workdir := flag.String("d", "", "set the working directory.")
	svcname := flag.String("name", app.Name(), "set the service name.")

	if cmd, ok := app.(Cmd); ok {
		cmd.Flag()
	}

	flag.Parse()

	chdir(*workdir)

	inService, err := svc.IsWindowsService()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to determine if we are running in service: %v\n", err)
		os.Exit(1)
	}

	if inService {
		runService(app, *svcname, false)
		return
	}

	arg := flag.Arg(0)
	switch arg {
	case "help", "usage":
		flag.CommandLine.SetOutput(os.Stdout)
		Help()
	case "version":
		fmt.Printf("%s.%s (%s)\n", app.Version(), app.Revision(), app.BuildTime().Local())
	case "install":
		err = installService(*svcname, app.DisplayName(), app.Description())
		if err == nil {
			fmt.Printf("service %q installed\n", *svcname)
		} else {
			fmt.Fprintf(os.Stderr, "Failed to %s service %q: %v\n", arg, *svcname, err)
			os.Exit(2)
		}
	case "remove":
		err = removeService(*svcname)
		if err == nil {
			fmt.Printf("service %q removed\n", *svcname)
		} else {
			fmt.Fprintf(os.Stderr, "Failed to %s service %q: %v\n", arg, *svcname, err)
			os.Exit(2)
		}
	case "start":
		err = startService(*svcname)
		if err == nil {
			fmt.Printf("service %q started\n", *svcname)
		} else {
			fmt.Fprintf(os.Stderr, "Failed to %s service %q: %v\n", arg, *svcname, err)
			os.Exit(2)
		}
	case "stop":
		err = controlService(*svcname, svc.Stop, svc.Stopped)
		if err == nil {
			fmt.Printf("service %q stoped\n", *svcname)
		} else {
			fmt.Fprintf(os.Stderr, "Failed to %s service %q: %v\n", arg, *svcname, err)
			os.Exit(2)
		}
	case "debug":
		runService(app, *svcname, true)
	case "":
		runStandalone(app)
	default:
		if cmd, ok := app.(Cmd); ok {
			cmd.Exec(arg)
		} else {
			fmt.Fprintf(os.Stderr, "Invalid command %q\n\n", arg)
			Help()
		}
	}

}
