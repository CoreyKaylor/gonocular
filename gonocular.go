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

type Renderer interface {
	RenderHtml(wr http.ResponseWriter, data interface{})
}

var (
	//dev mode by default
	optimize = os.Getenv("OPTIMIZE") != "true" && os.Getenv("OPTIMIZE") != ""
)

type DevModeRenderer struct {
	templateFiles []string
}

type ProductionRenderer struct {
	template *template.Template
}

type TemplateBuilder struct {
	templateFiles []string
}

func DevMode() bool {
	optimize = false
	return optimize
}

func ProductionMode() {
	optimize = true
}

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

func (v *TemplateBuilder) WithLayout() *TemplateBuilder {
	return v
}

func (v *TemplateBuilder) Template() Renderer {
	if optimize {
		return newProductionRenderer(v)
	} else {
		return &DevModeRenderer{v.templateFiles}
	}
}

func newProductionRenderer(tb *TemplateBuilder) *ProductionRenderer {
	t := template.Must(template.ParseFiles(tb.templateFiles...))
	return &ProductionRenderer{t}
}

func (r *ProductionRenderer) RenderHtml(wr http.ResponseWriter, data interface{}) {
	renderTemplate(r.template, wr, data)
}

func renderTemplate(t *template.Template, wr http.ResponseWriter, data interface{}) {
	wr.Header().Set("Content-Type", "text/html")
	t.Execute(wr, data)
}

type ParseError struct {
	FileName     string
	ErrorMessage string
	LineNumber   int
}

type SourceLine struct {
	Line     int
	HasError bool
	Text     string
}

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

func (r *DevModeRenderer) RenderHtml(wr http.ResponseWriter, data interface{}) {
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
