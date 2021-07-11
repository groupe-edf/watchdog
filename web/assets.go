package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
)

var Routes = []string{
	"/issues",
	"/policies",
	"/rules",
	"/repositories",
	"/settings",
}

//go:embed build/*
var assets embed.FS

type serveFile func(name string) (fs.File, error)

func (callback serveFile) Open(name string) (fs.File, error) {
	return callback(name)
}

func AssetsHandler(prefix, root string) http.Handler {
	handler := serveFile(func(name string) (fs.File, error) {
		assetPath := path.Join(root, name)
		assetFile, err := assets.Open(assetPath)
		if os.IsNotExist(err) {
			return assets.Open("build/index.html")
		}
		return assetFile, err
	})
	return http.StripPrefix(prefix, http.FileServer(http.FS(handler)))
}

func init() {
	if os.Getenv("BUILD") != "" {
		_, err := assets.Open("build/index.html")
		if err != nil {
			panic("failed to open build/index.html: have we built assets with 'yarn build'?")
		}
	}
}
