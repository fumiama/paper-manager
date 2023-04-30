package vben

import (
	"archive/zip"
	"bytes"
	_ "embed" // embed dist.zip
	"io/fs"
	"net/http"
)

//go:embed dist.zip
var distzipbytes []byte

// Distribution ...
var Distribution = func() http.FileSystem {
	distzip, err := zip.NewReader(bytes.NewReader(distzipbytes), int64(len(distzipbytes)))
	if err != nil {
		panic(err)
	}
	distfs, err := fs.Sub(fs.FS(distzip), "dist")
	if err != nil {
		panic(err)
	}
	return http.FS(distfs)
}()
