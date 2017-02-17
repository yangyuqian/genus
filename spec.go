package genus

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type SpecExtension struct {
	Framework string
	Data      json.RawMessage
}

type Spec struct {
	Suffix         string
	TemplateDir    string
	BaseDir        string
	BasePackage    string
	SkipExists     bool
	SkipFormat     bool
	SkipFixImports bool
	Merge          bool
	PlanItems      []*PlanItem
	Extension      *SpecExtension
	Data           []interface{}
}

func (spec *Spec) Save(baseDir string) (err error) {
	if baseDir == "" {
		baseDir = "./"
	}

	data, err := json.MarshalIndent(spec, " ", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(baseDir, "plan.json"), data, 0666)
}
