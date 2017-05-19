package design

import (
	"path/filepath"

	. "github.com/goadesign/goa/design/apidsl"
)

var path = "./example"
var indexHTML = filepath.Join(path, "index.html")
var assetsPath = filepath.Join(path, "assets")

var _ = Resource("files", func() {
	Files("/", indexHTML, func() {
	})
	Files("/assets/*filepath", assetsPath, func() {
	})
})
