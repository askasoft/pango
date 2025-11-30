package srv

import (
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/windows/svc"
)

func PrintDefaultCommand(out io.Writer) {
	fmt.Fprintln(out, "    service <command>")
	fmt.Fprintln(out, "      command=install    install as windows service.")
	fmt.Fprintln(out, "      command=remove     remove installed windows service.")
	fmt.Fprintln(out, "      command=start      start the windows service.")
	fmt.Fprintln(out, "      command=stop       stop the windows service.")
}

func PrintDefaultOptions(out io.Writer) {
	fmt.Fprintln(out, "    -h | -help           print the help message.")
	fmt.Fprintln(out, "    -v | -version        print the version message.")
	fmt.Fprintln(out, "    -dir                 set the working directory.")
	fmt.Fprintln(out, "    -name                set the service name.")
}

func PrintDefaultUsage(app App) {
	out := os.Stdout

	fmt.Fprintln(out, "Usage: "+app.Name()+".exe <command> [options]")

	fmt.Fprintln(out, "  <command>:")
	PrintDefaultCommand(out)

	fmt.Fprintln(out, "  <options>:")
	PrintDefaultOptions(out)
}

// Main server main
func Main(app App) {
	var (
		version bool
		workdir string
		svcname string
	)

	flag.BoolVar(&version, "v", false, "print version message.")
	flag.BoolVar(&version, "version", false, "print version message.")
	flag.StringVar(&workdir, "dir", "", "set the working directory.")
	flag.StringVar(&svcname, "name", app.Name(), "set the service name.")

	flag.CommandLine.Usage = app.Usage

	if cmd, ok := app.(Cmd); ok {
		cmd.Flag()
	}

	flag.Parse()

	chdir(workdir)

	if version {
		fmt.Println(app.Version())
		os.Exit(0)
	}

	inService, err := svc.IsWindowsService()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to determine if we are running in service: %v\n", err)
		os.Exit(1)
	}

	if inService {
		runService(app, svcname, false)
		return
	}

	arg := flag.Arg(0)
	switch arg {
	case "service":
		cmd := flag.Arg(1)
		switch cmd {
		case "install":
			err = installService(svcname, app.DisplayName(), app.Description())
			if err == nil {
				fmt.Printf("service %q installed\n", svcname)
			} else {
				fmt.Fprintf(os.Stderr, "Failed to %s service %q: %v\n", arg, svcname, err)
				os.Exit(2)
			}
		case "remove":
			err = removeService(svcname)
			if err == nil {
				fmt.Printf("service %q removed\n", svcname)
			} else {
				fmt.Fprintf(os.Stderr, "Failed to %s service %q: %v\n", arg, svcname, err)
				os.Exit(2)
			}
		case "start":
			err = startService(svcname)
			if err == nil {
				fmt.Printf("service %q started\n", svcname)
			} else {
				fmt.Fprintf(os.Stderr, "Failed to %s service %q: %v\n", arg, svcname, err)
				os.Exit(2)
			}
		case "stop":
			err = controlService(svcname, svc.Stop, svc.Stopped)
			if err == nil {
				fmt.Printf("service %q stoped\n", svcname)
			} else {
				fmt.Fprintf(os.Stderr, "Failed to %s service %q: %v\n", arg, svcname, err)
				os.Exit(2)
			}
		case "debug":
			runService(app, svcname, true)
		default:
			flag.CommandLine.SetOutput(os.Stdout)
			fmt.Fprintf(os.Stderr, "Invalid service command %q\n\n", cmd)
			app.Usage()
			os.Exit(1)
		}
	case "":
		runStandalone(app)
	default:
		if cmd, ok := app.(Cmd); ok {
			cmd.Exec(arg)
		} else {
			flag.CommandLine.SetOutput(os.Stdout)
			fmt.Fprintf(os.Stderr, "Invalid command %q\n\n", arg)
			app.Usage()
			os.Exit(1)
		}
	}

}
