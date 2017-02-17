package genus

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Imports map[string]string

type templateData struct {
	Package         string
	AbosultePackage string
	Imports         Imports
	Data            interface{}
}

// Represents single Go package
type TemplateGroup struct {
	Package         string // package name
	BaseDir         string // absolute directory for template group
	BasePackage     string // base package
	RelativePackage string
	AbsolutePackage string // absolute package combined with Package and BasePackage
	Filename        string
	Imports         Imports
	Templates       []*Template
	SkipFixImports  bool // skip fixing imports
	SkipExists      bool // skip generation if exist
	SkipFormat      bool // skip go format
	Merge           bool
}

func (tg *TemplateGroup) Render(data interface{}) (err error) {
	if tg.Merge {
		return tg.RenderPartial(data)
	}

	err = tg.ensureGOOS()
	if err != nil {
		return err
	}

	err = tg.ensureGopath()
	if err != nil {
		return err
	}

	err = tg.configureTemplates()
	if err != nil {
		return err
	}

	absPkg := filepath.Join(tg.BasePackage, tg.RelativePackage)
	tdata := &templateData{
		Package:         tg.Package,
		AbosultePackage: absPkg,
		Data:            data,
		Imports:         tg.Imports,
	}

	for _, t := range tg.Templates {
		_, err := t.Render(tdata)
		if err != nil {
			return err
		}
	}

	return
}

func (tg *TemplateGroup) RenderPartial(data interface{}) (err error) {
	err = tg.ensureGOOS()
	if err != nil {
		return err
	}

	err = tg.ensureGopath()
	if err != nil {
		return err
	}

	err = tg.configureTemplates()
	if err != nil {
		return err
	}

	absPkg := filepath.Join(tg.BasePackage, tg.RelativePackage)
	tdata := &templateData{
		Package:         tg.Package,
		AbosultePackage: absPkg,
		Data:            data,
		Imports:         tg.Imports,
	}

	mergedTemplate := Template{
		Filename:       tg.Filename,
		TargetDir:      filepath.Join(tg.BaseDir, tg.RelativePackage),
		SkipExists:     tg.SkipExists,
		SkipFormat:     tg.SkipFormat,
		SkipFixImports: tg.SkipFixImports,
	}

	for _, t := range tg.Templates {
		partial, err := t.RenderPartial(tdata)
		if err != nil {
			return err
		}
		mergedTemplate.rawTemplate = append(mergedTemplate.rawTemplate, '\n')
		mergedTemplate.rawTemplate = append(mergedTemplate.rawTemplate, partial...)
	}

	_, err = mergedTemplate.Render(tdata)
	return
}

func (tg *TemplateGroup) configureTemplates() (err error) {
	if tg.Package == "" {
		if idx := strings.LastIndex(tg.RelativePackage, "/"); idx >= 0 {
			tg.Package = tg.RelativePackage[(idx + 1):]
		} else if len(tg.RelativePackage) > 0 {
			tg.Package = tg.RelativePackage
		}
	}

	for _, t := range tg.Templates {
		t.TargetDir = filepath.Join(tg.BaseDir, tg.RelativePackage)

		if tg.Filename != "" {
			t.Filename = tg.Filename
		}

		if tg.Filename == "" {
			if idx := strings.LastIndex(t.Name, "/"); idx > 0 {
				t.Filename = t.Name[(idx+1):] + ".go"
			} else if len(t.Name) > 0 {
				t.Filename = t.Name + ".go"
			}
		}

		t.SkipExists = tg.SkipExists
		t.SkipFormat = tg.SkipFormat
		t.SkipFixImports = tg.SkipFixImports
	}

	imps := make(Imports)
	for imp, alias := range tg.Imports {
		// Local import
		// ./p1 ../p1
		if strings.HasPrefix(imp, ".") && !filepath.IsAbs(imp) {
			imp = filepath.Join(tg.BasePackage, imp)
		}

		imps[imp] = alias
	}

	tg.Imports = imps

	return
}

// TODO: V1 doesn't support windows
func (tg *TemplateGroup) ensureGOOS() (err error) {
	if runtime.GOOS == "windows" {
		return errors.New("Windows is not supported for now")
	}

	return
}

// Operations must be performed under gopath
func (tg *TemplateGroup) ensureGopath() (err error) {
	gopaths := strings.Split(os.Getenv("GOPATH"), ":")
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, gopath := range gopaths {
		if strings.HasPrefix(pwd, gopath) {
			if tg.BaseDir == "" {
				tg.BaseDir = pwd
			}

			if tg.BasePackage == "" {
				tg.BasePackage = pwd[len(gopath)+5:]
			}

			return nil
		}
	}

	return errors.New("Run outside gopath")
}
