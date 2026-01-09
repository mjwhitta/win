//go:build windows

package desktop

import (
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows/registry"

	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/pathname"
	w32 "github.com/mjwhitta/win/api"
)

// Wallpaper style consts
const (
	WallpaperStyleCenter uint = iota
	WallpaperStyleFill
	WallpaperStyleFit
	WallpaperStyleSpan
	WallpaperStyleStretch
	WallpaperStyleTile
)

// Wallpaper styles lookup map
var styles map[uint]string = map[uint]string{
	WallpaperStyleCenter:  "0",
	WallpaperStyleFill:    "10",
	WallpaperStyleFit:     "6",
	WallpaperStyleSpan:    "22",
	WallpaperStyleStretch: "2",
	WallpaperStyleTile:    "0",
}

// SetWallpaper will change the wallpaper to the specified image
// filepath.
func SetWallpaper(img string, style uint) error {
	var e error
	var fn uintptr
	var k registry.Key
	var path string

	if _, ok := styles[style]; !ok {
		return errors.Newf("unsupported wallpaper style: %d", style)
	}

	if img = strings.TrimSpace(img); img != "" {
		if ok, e := pathname.DoesExist(img); e != nil {
			return errors.Newf("file %s not accessible: %w", img, e)
		} else if !ok {
			return errors.Newf("file %s not found", img)
		}

		if path, e = filepath.Abs(img); e != nil {
			return errors.Newf(
				"failed to find absolute path for %s: %w",
				img,
				e,
			)
		}

		fn = uintptr(unsafe.Pointer(&[]byte(path)[0]))
	}

	// Get key
	k, _, e = registry.CreateKey(
		registry.CURRENT_USER,
		filepath.Join("Control Panel", "Desktop"),
		registry.SET_VALUE,
	)
	if e != nil {
		return errors.Newf("failed to get registry key: %w", e)
	}

	// Set wallpaper
	if img == "" {
		if e = k.SetStringValue("WallPaper", ""); e != nil {
			return errors.Newf("failed to set WallPaper value: %w", e)
		}
	}

	// Set wallpaper style
	e = k.SetStringValue("WallpaperStyle", styles[style])
	if e != nil {
		e = errors.Newf("failed to set WallpaperStyle value: %w", e)
		return e
	}

	// Set tiling
	switch style {
	case WallpaperStyleTile:
		e = k.SetStringValue("TileWallpaper", "1")
	default:
		e = k.SetStringValue("TileWallpaper", "0")
	}

	if e != nil {
		return errors.Newf("failed to set TileWallpaper value: %w", e)
	}

	// If the file doesn't exist, SystemParametersInfoA is supposed to
	// return an error. That does not, however, appear to be the case.
	// It never returns an error.
	e = w32.SystemParametersInfoA(
		w32.Winuser.SpiSetdeskwallpaper,
		0,
		fn,
		w32.Winuser.SpifSendchange|w32.Winuser.SpifUpdateinifile,
	)
	if e != nil {
		return errors.Newf(
			"failed to set desktop wallpaper to %s: %w",
			img,
			e,
		)
	}

	return nil
}
