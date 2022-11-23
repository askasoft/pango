package srv

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/sys/windows/svc"
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: %s [options] <command>\n"+
			"  <command>:\n"+
			"     install      install as windows service.\n"+
			"     remove       remove installed windows service.\n"+
			"     start        start the windows service.\n"+
			"     stop         stop the windows service.\n"+
			"     usage        print the usage information.\n"+
			"     version      print the version information.\n"+
			"  <options>:\n",
		os.Args[0])

	flag.PrintDefaults()
}

// Main server main
func Main(app App) {
	workdir := flag.String("d", "", "working directory")
	svcname := flag.String("name", app.Name(), "service name")

	flag.CommandLine.Usage = usage
	flag.Parse()

	if *workdir != "" {
		if err := os.Chdir(*workdir); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to change directory: %v\n", err)
			os.Exit(1)
		}
	}

	inService, err := svc.IsWindowsService()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to determine if we are running in service: %v\n", err)
		os.Exit(1)
	}
	if inService {
		runService(app, *svcname, false)
		return
	}

	cmd := flag.Arg(0)
	switch cmd {
	case "usage":
		flag.CommandLine.SetOutput(os.Stdout)
		usage()
		os.Exit(0)
	case "version":
		fmt.Printf("%s.%s (%s)\n", app.Version(), app.Revision(), app.BuildTime())
		os.Exit(0)
	case "install":
		err = installService(*svcname, app.DispName(), app.Description())
		if err == nil {
			fmt.Printf("service %s installed\n", *svcname)
		}
	case "remove":
		err = removeService(*svcname)
		if err == nil {
			fmt.Printf("service %s removed\n", *svcname)
		}
	case "start":
		err = startService(*svcname)
		if err == nil {
			fmt.Printf("service %s started\n", *svcname)
		}
	case "stop":
		err = controlService(*svcname, svc.Stop, svc.Stopped)
		if err == nil {
			fmt.Printf("service %s stoped\n", *svcname)
		}
	case "debug":
		runService(app, *svcname, true)
	case "":
		runStandalone(app)
	default:
		fmt.Fprintf(os.Stderr, "Invalid command %s\n\n", cmd)
		usage()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to %s %s: %v\n", cmd, *svcname, err)
		os.Exit(2)
	}
}

func runStandalone(app App) {
	app.Init()

	app.Run()

	wait(app)
}
