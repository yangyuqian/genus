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

// Represents single Go package
type TemplateGroup struct {
	Package        string
	BaseDir        string // absolute directory for template group
	BasePackage    string
	Imports        Imports
	Templates      []*Template
	SkipFixImports bool
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
			tg.BaseDir = pwd
			tg.BasePackage = strings.TrimRight(pwd, filepath.Join(gopath, "src"))

			return nil
		}
	}

	return errors.New("Run outside gopath")
}
