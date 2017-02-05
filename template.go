package genus

import (
	"bytes"
	"errors"
	gofmt "go/format"
	"io/ioutil"
	"text/template"
)

type Template struct {
	Name        string // template name
	Source      string // source path
	TargetDir   string // target directory
	Filename    string // filename of generated code
	SkipExists  bool   // skip generation if exist
	SkipFormat  bool   // skip go format
	rawTemplate []byte // rawTemplate data in bytes
	rawResult   []byte
}

// Set raw template data
func (tmpl *Template) SetRawTemplate(raw []byte) (data []byte) {
	tmpl.rawTemplate = raw
	return raw
}

// Render template by given data
func (tmpl *Template) Render(data interface{}) (result []byte, err error) {
	_, err = tmpl.load()
	if err != nil {
		return nil, err
	}

	return tmpl.render(data)
}

// Fix format of rawResult
func (tmpl *Template) format() (data []byte, err error) {
	data, err = gofmt.Source(tmpl.rawResult)
	if err != nil {
		return nil, err
	}

	tmpl.rawResult = data

	return
}

// Load data from file if rawTemplate is not set
func (tmpl *Template) load() (data []byte, err error) {
	if len(tmpl.rawTemplate) <= 0 {
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
	parsed, parsedErr := template.New(tmpl.Name).Parse(string(tmpl.rawTemplate))
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
