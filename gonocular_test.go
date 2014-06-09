package gonocular

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestComposingViews(t *testing.T) {
	Convey("When composing views", t, func() {
		Convey("During dev mode", func() {
			DevMode()
			Convey("With template parse errors", func() {
				template := TemplateFiles("fortests/error.html").Template()
				rw := httptest.NewRecorder()
				template.RenderHtml(rw, nil)
				Convey("Error page should include template source", func() {
					body := rw.Body.String()
					containsSource := strings.Contains(body, "{{ blah() }}")
					So(containsSource, ShouldBeTrue)
				})
				Convey("Content-Type should be text/html", func() {
					contentType := rw.Header().Get("Content-Type")
					So(contentType, ShouldEqual, "text/html")
				})
				Convey("Status is 500 InternalServerError", func() {
					So(rw.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
			Convey("Without template parse errors", func() {
				template := TemplateFiles("fortests/noerror.html").Template()
				rw := httptest.NewRecorder()
				template.RenderHtml(rw, nil)
				Convey("Page should render content", func() {
					body := rw.Body.String()
					containsSource := strings.Contains(body, "<h1>Hello World!</h1>")
					So(containsSource, ShouldBeTrue)
				})
				Convey("Content-Type should be text/html", func() {
					contentType := rw.Header().Get("Content-Type")
					So(contentType, ShouldEqual, "text/html")
				})
				Convey("Status is 200 OK", func() {
					So(rw.Code, ShouldEqual, http.StatusOK)
				})
			})
		})
		Convey("During production mode", func() {
			ProductionMode()
			Convey("With template parse errors", func() {
				Convey("Getting the template should panic", func() {
					templateFunc := func() {
						TemplateFiles("fortests/error.html").Template()
					}
					So(templateFunc, ShouldPanic)
				})
			})
			Convey("Without template parse errors", func() {
				template := TemplateFiles("fortests/noerror.html").Template()
				rw := httptest.NewRecorder()
				template.RenderHtml(rw, nil)
				Convey("Page should render content", func() {
					body := rw.Body.String()
					containsSource := strings.Contains(body, "<h1>Hello World!</h1>")
					So(containsSource, ShouldBeTrue)
				})
				Convey("Content-Type should be text/html", func() {
					contentType := rw.Header().Get("Content-Type")
					So(contentType, ShouldEqual, "text/html")
				})
				Convey("Status is 200 OK", func() {
					So(rw.Code, ShouldEqual, http.StatusOK)
				})
			})
		})
	})
}
