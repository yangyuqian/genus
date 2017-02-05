package genus

import (
	"errors"
	"log"
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
	AbsolutePackage string // absolute package combined with Package and BasePackage
	Imports         Imports
	Templates       []*Template
	SkipFixImports  bool // skip fixing imports
	SkipExists      bool // skip generation if exist
	SkipFormat      bool // skip go format
}

func (tg *TemplateGroup) Render(data interface{}) (err error) {
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

	log.Printf("Rendering template group with BasePackage: %v, Package: %v",
		tg.BasePackage, tg.Package)
	tdata := &templateData{
		Package:         tg.Package,
		AbosultePackage: filepath.Join(tg.BasePackage, tg.Package),
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

func (tg *TemplateGroup) configureTemplates() (err error) {
	for _, t := range tg.Templates {
		t.TargetDir = filepath.Join(tg.BaseDir, tg.Package)
		if idx := strings.LastIndex(t.Name, "/"); idx > 0 {
			t.Filename = t.Name[(idx+1):] + ".go"
		}

		t.SkipExists = tg.SkipExists
		t.SkipFormat = tg.SkipFormat
	}
	log.Printf("Configuring template group with BaseDir: %v, Package: %v",
		tg.BaseDir, tg.Package)

	for k, imp := range tg.Imports {
		// Local import
		// ./p1 ../p1
		if strings.HasPrefix(imp, ".") && !filepath.IsAbs(imp) {
			imp = filepath.Join(tg.BasePackage, imp)
		}

		tg.Imports[k] = imp
	}

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
			log.Printf("Ensure %s under $GOPATH %s", pwd, gopath)
			if tg.BaseDir == "" {
				tg.BaseDir = pwd
			}

			if tg.BasePackage == "" {
				tg.BasePackage = pwd[len(gopath)+5:]
			}
			log.Printf("BasePackage: %s", tg.BasePackage)

			return nil
		}
	}

	return errors.New("Run outside gopath")
}
