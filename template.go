package genus

import (
	"bytes"
	"errors"
	"fmt"
	gofmt "go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	goimports "golang.org/x/tools/imports"
)

var defaultHeader = []byte(`package {{ .Package }}

{{ with .Imports }}
import (
{{- range $k, $v := . }}
 {{- $v -}} "{{- $k -}}"
{{ end -}}
)
{{ end }}
`)

// Represents single Go file
type Template struct {
	Name           string // template name
	Source         string // source path
	TargetDir      string // target directory
	Filename       string // filename of generated code
	SkipExists     bool   // skip generation if exist
	SkipFormat     bool   // skip go format
	SkipFixImports bool   // skip generation if exist
	header         []byte
	rawTemplate    []byte // rawTemplate data in bytes
	rawResult      []byte
}

// Set raw template data
func (tmpl *Template) SetRawTemplate(raw []byte) (data []byte) {
	tmpl.rawTemplate = raw
	return raw
}

func (tmpl *Template) RenderPartial(data interface{}) (result []byte, err error) {
	_, err = tmpl.load()
	if err != nil {
		return nil, err
	}

	return tmpl.renderPartial(data)
}

// Render template by context
func (tmpl *Template) renderPartial(context interface{}) (data []byte, err error) {
	rawTemplate := fmt.Sprintf("{{ with .Data }}%s{{ end -}}", string(tmpl.rawTemplate))
	parsed, parsedErr := template.New(tmpl.Name).Funcs(GoHelperFuncs).Parse(string(rawTemplate))
	if parsedErr != nil {
		return nil, parsedErr
	}

	buf := bytes.NewBuffer([]byte{})
	if execErr := parsed.Execute(buf, context); execErr != nil {
		return nil, execErr
	}
	data = buf.Bytes()
	tmpl.rawResult = data
	return
}

// Render template by given data
func (tmpl *Template) Render(data interface{}) (result []byte, err error) {
	_, err = tmpl.load()
	if err != nil {
		return nil, err
	}

	result, err = tmpl.render(data)
	if err != nil {
		return
	}

	result, err = tmpl.format()
	if err != nil {
		return
	}

	err = tmpl.write()
	if err != nil {
		return nil, err
	}

	result, err = tmpl.fixImports()
	if err != nil {
		return nil, err
	}

	return
}

// Fix format of rawResult
func (tmpl *Template) format() (data []byte, err error) {
	if tmpl.SkipFormat {
		return tmpl.rawResult, nil
	}

	data, err = gofmt.Source(tmpl.rawResult)
	if err != nil {
		return nil, err
	}

	tmpl.rawResult = data

	return
}

func (tmpl *Template) fixImports() (data []byte, err error) {
	if tmpl.SkipFixImports {
		return tmpl.rawResult, nil
	}

	data, err = goimports.Process(filepath.Join(tmpl.TargetDir, tmpl.Filename), tmpl.rawResult, &goimports.Options{Comments: true})
	if err != nil {
		return nil, err
	}

	tmpl.rawResult = data
	if err = tmpl.write(); err != nil {
		return nil, err
	}
	return
}

// Load data from file if rawTemplate is not set
func (tmpl *Template) load() (data []byte, err error) {
	if len(tmpl.rawTemplate) <= 0 {
		log.Printf("raw template %s not set, loading from file %s", tmpl.Name, tmpl.Source)
		return tmpl.loadFile()
	}

	return tmpl.rawTemplate, nil
}

// Load raw template data from file
func (tmpl *Template) loadFile() (data []byte, err error) {
	if tmpl.Source == "" {
		return nil, errors.New("Empty source path")
	}

	data, err = ioutil.ReadFile(tmpl.Source)
	if err != nil {
		return nil, err
	}

	tmpl.rawTemplate = data
	return
}

// Render template by context
func (tmpl *Template) render(context interface{}) (data []byte, err error) {
	if len(tmpl.header) <= 0 {
		tmpl.header = defaultHeader
	}

	withHeader := fmt.Sprintf("%s{{ with .Data }}%s{{ end -}}", string(tmpl.header), string(tmpl.rawTemplate))
	parsed, parsedErr := template.New(tmpl.Name).Funcs(GoHelperFuncs).Parse(withHeader)
	if parsedErr != nil {
		return nil, parsedErr
	}

	buf := bytes.NewBuffer([]byte{})
	if execErr := parsed.Execute(buf, context); execErr != nil {
		return nil, execErr
	}
	data = buf.Bytes()

	fbuf := bytes.NewBuffer([]byte{})
	err = template.Must(template.New("filename").Funcs(GoHelperFuncs).Parse(fmt.Sprintf("{{- with .Data -}} %s {{- end -}}", tmpl.Filename))).Execute(fbuf, context)
	if err != nil {
		return nil, err
	}

	if len(fbuf.Bytes()) > 0 {
		tmpl.Filename = string(fbuf.Bytes())
	}

	tmpl.rawResult = data
	return
}

// Write rawResult to file
func (tmpl *Template) write() (err error) {
	if tmpl.TargetDir == "" {
		return
	}

	path := filepath.Join(tmpl.TargetDir, tmpl.Filename)

	if _, err := os.Stat(path); err == nil && tmpl.SkipExists {
		return nil
	}

	log.Printf("Creating directory %s", tmpl.TargetDir)
	err = os.MkdirAll(tmpl.TargetDir, 0777)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, tmpl.rawResult, 0666)
}
