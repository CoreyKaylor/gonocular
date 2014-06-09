gonocular
=========

Extremely lightweight wrapper for go's html/template package with dev-time friendly error messages (inspired by revel's template support). If the environment variable OPTIMIZE=true is specified the package falls back to the standard behavior of causing a panic if the template cannot be parsed with no friendly error messages.

Example Usage
============

~~~ go
package home

import (
	"github.com/CoreyKaylor/gonocular"
	"net/http"
)

var (
	homeTemplate = gonocular.TemplateFiles("../templates/layout.html", "index.html").Template()
)

func Index(rw http.ResponseWriter, r *http.Request) {
	homeTemplate.RenderHtml(rw, nil)
}
~~~