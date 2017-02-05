package genus

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// Plan for generations of singleton files
type SingletonPlan struct {
	// Set manually
	Suffix          string
	Package         string // leave empty will use the base directory name
	TemplateDir     string
	BaseDir         string
	BasePackage     string
	RelativePackage string
	Filename        string
	TemplateNames   []string
	Imports         Imports
	Data            interface{}
	SkipExists      bool
	SkipFormat      bool
	SkipFixImports  bool

	// Calculated
	Repo          *Repo
	TemplateGroup *TemplateGroup
}

func (p *SingletonPlan) Render() (err error) {
	if err = p.init(); err != nil {
		return err
	}

	if p.TemplateGroup == nil {
		return errors.New("SingletonPlan is not initialized correctly")
	}

	if err = p.TemplateGroup.Render(p.Data); err != nil {
		return err
	}

	return
}

func (p *SingletonPlan) validate() (err error) {
	if p.BaseDir == "" || p.TemplateDir == "" {
		return errors.New(fmt.Sprintf("Required BaseDir or TemplateDir not set"))
	}

	if len(p.TemplateNames) <= 0 {
		return errors.New("TemplateNames not set")
	}
	return
}

func (p *SingletonPlan) init() (err error) {
	if err = p.validate(); err != nil {
		return err
	}

	if !filepath.IsAbs(p.BaseDir) {
		absDir, err := filepath.Abs(p.BaseDir)
		if err != nil {
			return err
		}
		p.BaseDir = absDir
	}

	if p.BasePackage == "" {
		if idx := strings.Index(p.BaseDir, "/src/"); idx > 0 {
			p.BasePackage = p.BaseDir[(idx + 5):]
		}
	}

	repo := NewRepo(p.TemplateDir, p.Suffix)
	err = repo.Load()
	if err != nil {
		return err
	}

	p.Repo = repo

	p.TemplateGroup, err = repo.BuildGroup(p.TemplateNames...)
	if err != nil {
		return err
	}
	p.TemplateGroup.BaseDir = p.BaseDir
	p.TemplateGroup.BasePackage = p.BasePackage

	if idx := strings.LastIndex(p.RelativePackage, "/"); idx >= 0 {
		p.TemplateGroup.Package = p.RelativePackage[(idx + 1):]
	} else if len(p.RelativePackage) > 0 {
		p.TemplateGroup.Package = p.RelativePackage
	}

	if p.Package != "" {
		p.TemplateGroup.Package = p.Package
	}

	if p.Filename != "" {
		p.TemplateGroup.Filename = p.Filename
	}

	p.TemplateGroup.RelativePackage = p.RelativePackage

	if len(p.Imports) > 0 {
		p.TemplateGroup.Imports = p.Imports
	}

	return
}

// PlanType of SingletonPlan
func (p *SingletonPlan) Type() PlanType {
	return SINGLETON
}
