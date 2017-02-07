package genus

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func NewRepo(templateDir, suffix string) *Repo {
	return &Repo{Suffix: suffix, TemplateDir: templateDir}
}

// Template repository
type Repo struct {
	Suffix        string
	TemplateDir   string
	Templates     []*Template
	templateNames []string
}

// Load templates
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
			suffixLen := len(r.Suffix)
			relName, err := filepath.Rel(r.TemplateDir, path)
			if err != nil {
				return err
			}

			r.templateNames = append(r.templateNames, relName)
			tmplName := relName[0:(len(relName) - suffixLen)]
			r.Templates = append(r.Templates, &Template{
				Name:   tmplName,
				Source: path,
			})
			log.Printf("Register template %s at %s", tmplName, path)
		}

		return nil
	}))
}

// Build template group with given template names
func (r *Repo) BuildGroup(names ...string) (tg *TemplateGroup, err error) {
	tg = &TemplateGroup{}
	for _, name := range names {
		t, err := r.Lookup(name)
		if err != nil {
			return nil, err
		}

		tg.Templates = append(tg.Templates, t)
	}

	return
}

// Loakup template by name
func (r *Repo) Lookup(name string) (t *Template, err error) {
	for _, tmpl := range r.Templates {
		if tmpl.Name == name {
			return tmpl, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Template %s not registered", name))
}
