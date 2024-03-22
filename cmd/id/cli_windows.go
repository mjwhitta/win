//go:build windows

package main

import (
	"os"

	"github.com/mjwhitta/cli"
	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/win"
)

// Exit status
const (
	Good = iota
	InvalidOption
	MissingOption
	InvalidArgument
	MissingArgument
	ExtraArgument
	Exception
)

// Flags
var flags struct {
	all     bool
	groups  bool
	nocolor bool
	privs   bool
	user    bool
	verbose bool
	version bool
}

func init() {
	// Configure cli package
	cli.Align = true // Defaults to false
	cli.Authors = []string{"Miles Whittaker <mj@whitta.dev>"}
	cli.Banner = hl.Sprintf("%s [OPTIONS]", os.Args[0])
	cli.BugEmail = "id.bugs@whitta.dev"
	cli.ExitStatus(
		"Normally the exit status is 0. In the event of an error the",
		"exit status will be one of the below:\n\n",
		hl.Sprintf("  %d: Invalid option\n", InvalidOption),
		hl.Sprintf("  %d: Missing option\n", MissingOption),
		hl.Sprintf("  %d: Invalid argument\n", InvalidArgument),
		hl.Sprintf("  %d: Missing argument\n", MissingArgument),
		hl.Sprintf("  %d: Extra argument\n", ExtraArgument),
		hl.Sprintf("  %d: Exception", Exception),
	)
	cli.Info("Returns info very similar to the whoami command.")
	// cli.MaxWidth = 80 // Defaults to 80
	cli.SeeAlso = []string{"whoami"}
	// cli.TabWidth = 4 // Defaults to 4
	cli.Title = "ID"

	// Parse cli flags
	cli.Flag(
		&flags.all,
		"a",
		"all",
		false,
		"Display user info, groups info, and privs info.",
	)
	cli.Flag(
		&flags.groups,
		"g",
		"groups",
		false,
		"Display groups info.",
	)
	cli.Flag(
		&flags.nocolor,
		"no-color",
		false,
		"Disable colorized output.",
	)
	cli.Flag(
		&flags.privs,
		"p",
		"privs",
		false,
		"Display privs info.",
	)
	cli.Flag(
		&flags.user,
		"u",
		"user",
		false,
		"Display user info.",
	)
	cli.Flag(
		&flags.verbose,
		"v",
		"verbose",
		false,
		"Show stacktrace, if error.",
	)
	cli.Flag(&flags.version, "V", "version", false, "Show version.")
	cli.Parse()
}

// Process cli flags and ensure no issues
func validate() {
	hl.Disable(flags.nocolor)

	// Short circuit, if version was requested
	if flags.version {
		hl.Printf("id version %s\n", win.Version)
		os.Exit(Good)
	}

	// Validate cli flags
	if cli.NArg() > 0 {
		cli.Usage(ExtraArgument)
	}

	if flags.all {
		flags.groups = true
		flags.privs = true
		flags.user = true
	}
}
