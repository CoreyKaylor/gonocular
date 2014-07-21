package gonocular

import (
	"bufio"
	"html/template"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

//Renderer is the template rendering interface
type Renderer interface {
	RenderHTML(wr http.ResponseWriter, data interface{})
}

var templateFactory = newDevRenderer

//DevModeRenderer is responsible for rendering templates and when an error occurs
//will output the error details
type DevModeRenderer struct {
	templateFiles []string
}

//ProductionRenderer is responsible for only rendering the template and will
//panic in the occurrence of an error. This is equivalent to template.Must(...
type ProductionRenderer struct {
	template *template.Template
}

//TemplateBuilder is responsible for building the collection of template files.
type TemplateBuilder struct {
	templateFiles []string
}

//DevMode changes the template behavior to include error details.
//This is the default behavior.
func DevMode() {
	templateFactory = newDevRenderer
}

//ProductionMode makes the template behavior panic only when there is
//an error. This is equivalent to calling template.Must(...
func ProductionMode() {
	templateFactory = newProductionRenderer
}

//TemplateFiles is the list of file paths relative from the callers code
//file.
func TemplateFiles(filenames ...string) *TemplateBuilder {
	templates := make([]string, 0, 10)
	for _, file := range filenames {
		file := filePathRelativeFromCaller(2, file)
		templates = append(templates, file)
	}
	builder := &TemplateBuilder{templates}
	return builder
}

func filePathRelativeFromCaller(skip int, file string) string {
	_, filename, _, _ := runtime.Caller(skip)
	dir := path.Dir(filename)
	file = path.Join(dir, file)
	return file
}

//Template will return the Dev or Production Renderer implementation based on
//current settings
func (v *TemplateBuilder) Template() Renderer {
	return templateFactory(v)
}

func newDevRenderer(tb *TemplateBuilder) Renderer {
	return &DevModeRenderer{tb.templateFiles}
}

func newProductionRenderer(tb *TemplateBuilder) Renderer {
	t := template.Must(template.ParseFiles(tb.templateFiles...))
	return &ProductionRenderer{t}
}

//RenderHTML will render an html/template file and panic if there is an error
func (r *ProductionRenderer) RenderHTML(wr http.ResponseWriter, data interface{}) {
	renderTemplate(r.template, wr, data)
}

func renderTemplate(t *template.Template, wr http.ResponseWriter, data interface{}) {
	wr.Header().Set("Content-Type", "text/html")
	t.Execute(wr, data)
}

//ParseError is details the error template uses to render the error page
type ParseError struct {
	FileName     string
	ErrorMessage string
	LineNumber   int
}

//SourceLine is the detail about a single line in a template file used for
//rendering the error page.
type SourceLine struct {
	Line     int
	HasError bool
	Text     string
}

//SourceLines will parse out the error details of a template error.
func (p *ParseError) SourceLines() []*SourceLine {
	lines := make([]*SourceLine, 0, 10)
	file, err := os.Open(p.FileName)
	if err != nil {
		return lines
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		hasError := lineCount == p.LineNumber
		sl := &SourceLine{lineCount, hasError, scanner.Text()}
		lines = append(lines, sl)
		lineCount++
	}
	return lines
}

//RenderHTML will render an html/template or in the case of an error will display
//an errorpage with the error details.
func (r *DevModeRenderer) RenderHTML(wr http.ResponseWriter, data interface{}) {
	t, err := template.ParseFiles(r.templateFiles...)
	if err == nil {
		renderTemplate(t, wr, data)
	} else {
		errorFile := filePathRelativeFromCaller(1, "error.html")
		et := template.Must(template.ParseFiles(errorFile))
		errorMessage := err.Error()
		errorParts := strings.Split(errorMessage, ":")
		parseError := &ParseError{}
		fileName := strings.TrimSpace(errorParts[1])
		for _, file := range r.templateFiles {
			if strings.HasSuffix(file, fileName) {
				parseError.FileName = file
			}
		}
		parseError.ErrorMessage = errorMessage
		if len(errorParts) > 2 {
			parseError.LineNumber, _ = strconv.Atoi(errorParts[2])
		}
		wr.WriteHeader(http.StatusInternalServerError)
		renderTemplate(et, wr, parseError)
	}
}
