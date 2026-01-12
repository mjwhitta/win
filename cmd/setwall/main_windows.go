//go:build windows

package main

import (
	"bytes"
	"fmt"
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
	var buf *bytes.Buffer = &bytes.Buffer{}
	var f *os.File
	var fn string = filepath.Join("c:/", "windows", "temp", "swframe")
	var frame string
	var g *gif.GIF
	var sig chan os.Signal = make(chan os.Signal, 1)

	if f, e = os.Open(filepath.Clean(img)); e != nil {
		return e
	}
	defer func() {
		if e2 := f.Close(); (e == nil) && (e2 != nil) {
			e = e2
		}
	}()

	if g, e = gif.DecodeAll(f); e != nil {
		return e
	}

	// Setup SIGINT for stopping
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for i := range g.Image {
		buf.Reset()

		e = png.Encode(buf, g.Image[i].SubImage(g.Image[i].Bounds()))
		if e != nil {
			return e
		}

		frame = fmt.Sprintf("%s%d.png", fn, i)

		// Write frame to c:/windows/temp/swframe#.png
		//nolint:mnd // u=rw,go=-
		e = os.WriteFile(filepath.Clean(frame), buf.Bytes(), 0o600)
		if e != nil {
			return e
		}
	}

	for {
		for i := range g.Image {
			select {
			case <-sig:
				signal.Stop(sig)

				for i := range g.Image {
					frame = fmt.Sprintf("%s%d.png", fn, i)
					_ = os.Remove(frame)
				}

				return nil
			default:
				frame = fmt.Sprintf("%s%d.png", fn, i)

				// Set wallpaper to frame
				e = desktop.SetWallpaper(frame, styles[flags.style])
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
