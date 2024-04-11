package srv

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func PrintDefaultCommand(out io.Writer) {
	fmt.Fprintln(out, "    version             print the version information.")
	fmt.Fprintln(out, "    help | usage        print the usage information.")
}

func PrintDefaultOptions() {
	flag.PrintDefaults()
}

func PrintUsage(app App) {
	out := flag.CommandLine.Output()

	fmt.Fprintln(out, "Usage: "+app.Name()+" <command> [options]")
	fmt.Fprintln(out, "  <command>:")
	if cmd, ok := app.(Cmd); ok {
		cmd.PrintCommand(out)
	} else {
		PrintDefaultCommand(out)
	}

	fmt.Fprintln(out, "  <options>:")
	PrintDefaultOptions()
}

// Main server main
func Main(app App) {
	flag.CommandLine.Usage = app.Usage

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
		app.Usage()
	case "version":
		fmt.Println(app.Version())
	case "":
		runStandalone(app)
	default:
		if cmd, ok := app.(Cmd); ok {
			cmd.Exec(arg)
		} else {
			flag.CommandLine.SetOutput(os.Stdout)
			fmt.Fprintf(os.Stderr, "Invalid command %q\n\n", arg)
			app.Usage()
		}
	}
}
