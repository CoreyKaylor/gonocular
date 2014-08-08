package main

import (
	"log"
	"net/http"
	"os"

	"github.com/CoreyKaylor/gonocular"
	"github.com/julienschmidt/httprouter"
)

var (
	template = gonocular.TemplateFiles("notfound.html").Template()
)

func handler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	template.RenderHtml(w, nil)
}

func main() {
	router := httprouter.New()
	router.GET("/", handler)
	router.ServeFiles("/public/*filepath", http.Dir("public"))
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}
