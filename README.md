gonocular
=========
[![Build Status](https://travis-ci.org/CoreyKaylor/gonocular.png)](https://travis-ci.org/CoreyKaylor/gonocular)
[![Coverage Status](https://coveralls.io/repos/CoreyKaylor/gonocular/badge.png?branch=master)](https://coveralls.io/r/CoreyKaylor/gonocular?branch=master)
[![GoDoc](https://godoc.org/github.com/CoreyKaylor/gonocular?status.svg)](http://godoc.org/github.com/CoreyKaylor/gonocular)

![Gonocular](https://raw.githubusercontent.com/CoreyKaylor/gonocular/master/gopher-grave.png "gonocular")

Extremely lightweight wrapper for go's html/template package with dev-time friendly error messages (originally inspired by revel's template support). If gonocular.ProductionMode() is called the package falls back to the standard behavior of causing a panic if the template cannot be parsed with no friendly error messages.

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
	homeTemplate.RenderHTML(rw, nil)
}
~~~
