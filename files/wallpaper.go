package files

import (
	"path/filepath"

	"github.com/davenicholson-xyz/go-setwallpaper/wallpaper"
	"github.com/davenicholson-xyz/wallmancer/download"
)

func ApplyWallpaper(file string, provider string) (string, error) {
	filename := filepath.Base(file)
	cache_dir, _ := GetCacheDir()
	output := filepath.Join(cache_dir, provider, filename)

	_ = download.DownloadImage(file, output)
	wallpaper.Set(output)

	return output, nil
}
