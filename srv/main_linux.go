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
			"     usage        print the usage information.\n"+
			"     version      print the version information.\n"+
			"  <options>:\n",
		os.Args[0])

	flag.PrintDefaults()
}

// Main server main
func Main(app App) {
	workdir := flag.String("d", "", "working directory")

	flag.CommandLine.Usage = usage
	flag.Parse()

	switch flag.Arg(0) {
	case "usage":
		flag.CommandLine.SetOutput(os.Stdout)
		usage()
		os.Exit(0)
	case "version":
		fmt.Printf("%s.%s (%s)\n", app.Version(), app.Revision(), app.BuildTime())
		os.Exit(0)
	}

	if *workdir != "" {
		if err := os.Chdir(*workdir); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to change directory: %v\n", err)
			os.Exit(1)
		}
	}

	app.Init()

	app.Run()

	wait(app)
}
