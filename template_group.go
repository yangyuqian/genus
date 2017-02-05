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
	Package         string
	BaseDir         string // absolute directory for template group
	BasePackage     string
	AbosultePackage string
	Imports         Imports
	Templates       []*Template
	SkipFixImports  bool
}

func (tg *TemplateGroup) Render(data interface{}) (err error) {
	// tdata := templateData{
	// 	Package:         tg.Package,
	// 	AbosultePackage: filepath.Join(tg.BasePackage, tg.Package),
	// }

	return
}

func (tg *TemplateGroup) configureTemplates() (err error) {
	for _, t := range tg.Templates {
		t.TargetDir = filepath.Join(tg.BaseDir, tg.Package)
		if idx := strings.LastIndex(t.Name, "/"); idx > 0 {
			t.Filename = t.Name[(idx+1):] + ".go"
		}
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
			log.Printf("Ensure %s under $GOPATH", pwd)
			if tg.BaseDir == "" {
				tg.BaseDir = pwd
			}

			if tg.BasePackage == "" {
				tg.BasePackage = strings.TrimRight(pwd, filepath.Join(gopath, "src"))
			}

			return nil
		}
	}

	return errors.New("Run outside gopath")
}
