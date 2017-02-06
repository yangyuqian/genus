package genus

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Template repository
type Repo struct {
	Suffix        string
	TemplateDir   string
	Templates     []*Template
	templateNames []string
}

func (r *Repo) Load() (err error) {
	if r.TemplateDir == "" {
		return errors.New("TemplateDir not set")
	}

	return filepath.Walk(r.TemplateDir, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return errors.New(fmt.Sprintf("Directory or file %s not found", path))
		}

		// skip directories
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, r.Suffix) {
			r.templateNames = append(r.templateNames, path)
			r.Templates = append(r.Templates, &Template{
				Name:   strings.TrimRight(path, r.Suffix),
				Source: path,
			})
			log.Printf("Register template %s", path)
		}

		return nil
	}))
}
