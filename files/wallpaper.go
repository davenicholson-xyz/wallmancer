package files

import (
	"fmt"
	"path/filepath"

	"github.com/davenicholson-xyz/go-setwallpaper/wallpaper"
	"github.com/davenicholson-xyz/wallmancer/download"
)

func ApplyWallpaper(file string, provider string) error {
	filename := filepath.Base(file)
	cache_dir, _ := GetCacheDir()
	output := filepath.Join(cache_dir, provider, filename)
	fmt.Println(output)

	_ = download.DownloadImage(file, output)
	wallpaper.Set(output)
	return nil
}
