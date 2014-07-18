package main

import (
	"log"
	"net/http"
	"os"

	"github.com/CoreyKaylor/gonocular"
)

var (
	template = gonocular.TemplateFiles("notfound.html").Template()
)

func handler(w http.ResponseWriter, r *http.Request) {
	template.RenderHtml(w, nil)
}

func main() {
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
