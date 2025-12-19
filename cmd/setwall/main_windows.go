//go:build windows

package main

import (
	"github.com/mjwhitta/cli"
	"github.com/mjwhitta/log"
	"github.com/mjwhitta/win/desktop"
)

var styles map[string]uint = map[string]uint{
	"center":  desktop.WallpaperStyleCenter,
	"fill":    desktop.WallpaperStyleFill,
	"fit":     desktop.WallpaperStyleFit,
	"span":    desktop.WallpaperStyleSpan,
	"stretch": desktop.WallpaperStyleStretch,
	"tile":    desktop.WallpaperStyleTile,
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			if flags.verbose {
				panic(r)
			}

			switch r := r.(type) {
			case error:
				log.ErrX(Exception, r.Error())
			case string:
				log.ErrX(Exception, r)
			}
		}
	}()

	var e error
	var img string

	validate()

	if cli.NArg() > 0 {
		img = cli.Arg(0)
	}

	if e = desktop.SetWallpaper(img, styles[flags.style]); e != nil {
		panic(e)
	}
}
