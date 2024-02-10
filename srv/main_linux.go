package srv

import (
	"flag"
	"fmt"
	"os"
)

// Main server main
func Main(app App) {
	Usage = func() {
		out := flag.CommandLine.Output()

		fmt.Fprintln(out, "Usage: "+app.Name()+" <command> [options]")
		fmt.Fprintln(out, "  <command>:")
		if cmd, ok := app.(Cmd); ok {
			cmd.CmdHelp(out)
		}
		fmt.Fprintln(out, "    version         print the version information.")
		fmt.Fprintln(out, "    help | usage    print the usage information.")
		fmt.Fprintln(out, "  <options>:")

		flag.PrintDefaults()
	}

	flag.CommandLine.Usage = Usage

	workdir := flag.String("d", "", "set the working directory.")

	if cmd, ok := app.(Cmd); ok {
		cmd.Flag()
	}

	flag.Parse()

	chdir(*workdir)

	arg := flag.Arg(0)
	switch arg {
	case "help", "usage":
		flag.CommandLine.SetOutput(os.Stdout)
		Help()
	case "version":
		fmt.Printf("%s.%s (%s)\n", app.Version(), app.Revision(), app.BuildTime().Local())
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
