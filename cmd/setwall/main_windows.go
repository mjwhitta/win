//go:build windows

package main

import (
	"image/gif"
	"image/png"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

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

	if strings.HasSuffix(img, ".gif") {
		if e = wallpaperSlideshow(img); e != nil {
			panic(e)
		}

		img = ""
	}

	if e = desktop.SetWallpaper(img, styles[flags.style]); e != nil {
		panic(e)
	}
}

//nolint:wrapcheck // Not wrapping errors in a main package
func wallpaperSlideshow(img string) (e error) {
	var fn string = filepath.Join("c:/", "windows", "temp", "mw.png")
	var fr *os.File
	var fw *os.File
	var g *gif.GIF
	var sig chan os.Signal = make(chan os.Signal, 1)

	if fr, e = os.Open(filepath.Clean(img)); e != nil {
		return e
	}
	defer func() {
		if e2 := fr.Close(); (e == nil) && (e2 != nil) {
			e = e2
		}
	}()

	if g, e = gif.DecodeAll(fr); e != nil {
		return e
	}

	// Setup SIGINT for stopping
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		for i := range g.Image {
			select {
			case <-sig:
				signal.Stop(sig)

				_ = os.Remove(fn)

				return nil
			default:
				// Write frame to c:/windows/temp
				if fw, e = os.Create(filepath.Clean(fn)); e != nil {
					return e
				}

				e = png.Encode(
					fw,
					g.Image[i].SubImage(g.Image[i].Bounds()),
				)
				if e != nil {
					_ = fw.Close()
					return e
				}

				if e = fw.Close(); e != nil {
					return e
				}

				// Set wallpaper to frame
				e = desktop.SetWallpaper(fn, styles[flags.style])
				if e != nil {
					panic(e)
				}

				// Sleep for the specified delay
				//
				// NOTE: This is supposed to be in 100ths of a second,
				// but looks to actually be milliseconds.
				time.Sleep(
					time.Duration(g.Delay[i]) * time.Millisecond,
				)
			}
		}
	}
}
