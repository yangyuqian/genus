package genus

import (
	"errors"
	"io/ioutil"
)

type Template struct {
	Name        string // template name
	Source      string // source path
	TargetDir   string // target directory
	Filename    string // filename of generated code
	SkipExists  bool   // skip generation if exist
	SkipFormat  bool   // skip go format
	rawTemplate []byte // rawTemplate data in bytes
}

// Set raw template data
func (tmpl *Template) SetRawTemplate(raw []byte) (data []byte) {
	tmpl.rawTemplate = raw
	return
}

// load raw template data from file
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
