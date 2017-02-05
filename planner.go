package genus

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type PlanItem struct {
	Package         string
	PlanType        string
	Suffix          string
	TemplateDir     string
	BaseDir         string
	BasePackage     string
	RelativePackage string
	Filename        string
	TemplateNames   []string
	SkipExists      bool
	SkipFormat      bool
	SkipFixImports  bool
	Imports         Imports
	Data            []interface{}
}

type Spec struct {
	Suffix         string
	TemplateDir    string
	BaseDir        string
	BasePackage    string
	SkipExists     bool
	SkipFormat     bool
	SkipFixImports bool
	PlanItems      []*PlanItem
}

type Planner interface{}

func NewPackagePlanner(specPath string) (planner *PackagePlanner, err error) {
	raw, err := ioutil.ReadFile(specPath)
	if err != nil {
		return nil, err
	}
	return &PackagePlanner{RawSpec: raw}, nil
}

// Planner for generations for single Go package
// Load plan definitions from a json specification
type PackagePlanner struct {
	RawSpec []byte
	Spec    *Spec
	Plans   []Plan
}

func (pl *PackagePlanner) Perform() (err error) {
	for _, plan := range pl.Plans {
		if err = plan.Render(); err != nil {
			return err
		}
	}

	return
}

func (pl *PackagePlanner) Plan() (err error) {
	pl.Spec = &Spec{}
	if err = json.Unmarshal(pl.RawSpec, pl.Spec); err != nil {
		return err
	}

	for _, planItem := range pl.Spec.PlanItems {
		switch planItem.PlanType {
		case SINGLETON.String():
			plan := &SingletonPlan{
				Suffix:          StringWithDefault(planItem.Suffix, pl.Spec.Suffix),
				TemplateDir:     StringWithDefault(planItem.TemplateDir, pl.Spec.TemplateDir),
				BaseDir:         StringWithDefault(planItem.BaseDir, pl.Spec.BaseDir),
				BasePackage:     StringWithDefault(planItem.BasePackage, pl.Spec.BasePackage),
				SkipExists:      BoolWithDefault(planItem.SkipExists, pl.Spec.SkipExists),
				SkipFormat:      BoolWithDefault(planItem.SkipFormat, pl.Spec.SkipFormat),
				SkipFixImports:  BoolWithDefault(planItem.SkipFixImports, pl.Spec.SkipFixImports),
				RelativePackage: planItem.RelativePackage,
				TemplateNames:   planItem.TemplateNames,
				Filename:        planItem.Filename,
				Package:         planItem.Package,
				Imports:         planItem.Imports,
			}

			if len(planItem.Data) > 0 {
				plan.Data = planItem.Data[0]
			}

			pl.Plans = append(pl.Plans, plan)
		case REPEATABLE.String():
			pl.Plans = append(pl.Plans, &RepeatablePlan{
				SingletonPlan: SingletonPlan{
					Suffix:          StringWithDefault(planItem.Suffix, pl.Spec.Suffix),
					TemplateDir:     StringWithDefault(planItem.TemplateDir, pl.Spec.TemplateDir),
					BaseDir:         StringWithDefault(planItem.BaseDir, pl.Spec.BaseDir),
					BasePackage:     StringWithDefault(planItem.BasePackage, pl.Spec.BasePackage),
					SkipExists:      BoolWithDefault(planItem.SkipExists, pl.Spec.SkipExists),
					SkipFormat:      BoolWithDefault(planItem.SkipFormat, pl.Spec.SkipFormat),
					SkipFixImports:  BoolWithDefault(planItem.SkipFixImports, pl.Spec.SkipFixImports),
					RelativePackage: planItem.RelativePackage,
					Filename:        planItem.Filename,
					Package:         planItem.Package,
					TemplateNames:   planItem.TemplateNames,
					Imports:         planItem.Imports,
				},
				Data: planItem.Data,
			})
		default:
			return errors.New(fmt.Sprintf("Invalid plan type %s, PlanType must be either SINGLETON or REPEATABLE"))
		}
	}

	return
}
