package srv

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: %s <command> [options]\n"+
			"  <command>:\n"+
			"    help | usage    print the usage information.\n"+
			"    version         print the version information.\n"+
			"  <options>:\n",
		os.Args[0])

	flag.PrintDefaults()
}

// Main server main
func Main(app App) {
	workdir := flag.String("d", "", "set the working directory.")

	flag.CommandLine.Usage = usage
	flag.Parse()

	if *workdir != "" {
		if err := os.Chdir(*workdir); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to change directory: %v\n", err)
			os.Exit(1)
		}
	}

	cmd := flag.Arg(0)
	switch cmd {
	case "help", "usage":
		flag.CommandLine.SetOutput(os.Stdout)
		usage()
	case "version":
		fmt.Printf("%s.%s (%s)\n", app.Version(), app.Revision(), app.BuildTime().Local())
	case "":
		runStandalone(app)
	default:
		fmt.Fprintf(os.Stderr, "Invalid command %q\n\n", cmd)
		usage()
	}
}
