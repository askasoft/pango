package srv

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func PrintDefaultCommand(out io.Writer) {
}

func PrintDefaultOptions(out io.Writer) {
	fmt.Fprintln(out, "    -h | -help          print the help message.")
	fmt.Fprintln(out, "    -v | -version       print the version message.")
	fmt.Fprintln(out, "    -dir                set the working directory.")
}

func PrintDefaultUsage(app App) {
	out := os.Stdout

	fmt.Fprintln(out, "Usage: "+app.Name()+" [options]")
	fmt.Fprintln(out, "  <options>:")
	PrintDefaultOptions(out)
}

// Main server main
func Main(app App) {
	var (
		version bool
		workdir string
	)

	flag.BoolVar(&version, "v", false, "print version message.")
	flag.BoolVar(&version, "version", false, "print version message.")
	flag.StringVar(&workdir, "dir", "", "set the working directory.")

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

	arg := flag.Arg(0)
	switch arg {
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
