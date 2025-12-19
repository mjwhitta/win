//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mjwhitta/cli"
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
	style   string
	verbose bool
	version bool
}

func init() {
	// Configure cli package
	cli.Align = true // Defaults to false
	cli.Authors = []string{"Miles Whittaker <mj@whitta.dev>"}
	cli.Banner = filepath.Base(os.Args[0]) + " [OPTIONS] [img]"
	cli.BugEmail = "setwall.bugs@whitta.dev"

	cli.ExitStatus(
		"Normally the exit status is 0. In the event of an error the",
		"exit status will be one of the below:\n\n",
		fmt.Sprintf("  %d: Invalid option\n", InvalidOption),
		fmt.Sprintf("  %d: Missing option\n", MissingOption),
		fmt.Sprintf("  %d: Invalid argument\n", InvalidArgument),
		fmt.Sprintf("  %d: Missing argument\n", MissingArgument),
		fmt.Sprintf("  %d: Extra argument\n", ExtraArgument),
		fmt.Sprintf("  %d: Exception", Exception),
	)
	cli.Info("Sets the user's desktop wallpaper.")
	cli.Section(
		"WALLPAPER STYLES",
		"Supported wallpapers styles include: center, fill, fit,",
		"span, stretch, and tile. The default is stretch.",
	)

	cli.Title = "Set Desktop Wallpaper"

	// Parse cli flags
	cli.Flag(
		&flags.style,
		"s",
		"style",
		"stretch",
		"Set wallpaper style (default: stretch).",
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
	// Short circuit, if version was requested
	if flags.version {
		fmt.Println(
			filepath.Base(os.Args[0]) + " version " + win.Version,
		)
		os.Exit(Good)
	}

	if _, ok := styles[flags.style]; !ok {
		cli.Usage(InvalidOption)
	}

	// Validate cli flags
	if cli.NArg() > 1 {
		cli.Usage(ExtraArgument)
	}
}
