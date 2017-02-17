package genus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/yangyuqian/genus/generator/orm"
	"github.com/yangyuqian/genus/types"
)

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
	if pl.Spec == nil {
		pl.Spec = &Spec{}
		if err = json.Unmarshal(pl.RawSpec, pl.Spec); err != nil {
			return err
		}
	}

	for _, planItem := range pl.Spec.PlanItems {
		switch planItem.PlanType {
		case types.SINGLETON.String():
			plan := &SingletonPlan{
				Suffix:          StringWithDefault(planItem.Suffix, pl.Spec.Suffix),
				TemplateDir:     StringWithDefault(planItem.TemplateDir, pl.Spec.TemplateDir),
				BaseDir:         StringWithDefault(planItem.BaseDir, pl.Spec.BaseDir),
				BasePackage:     StringWithDefault(planItem.BasePackage, pl.Spec.BasePackage),
				SkipExists:      BoolWithDefault(planItem.SkipExists, pl.Spec.SkipExists),
				SkipFormat:      BoolWithDefault(planItem.SkipFormat, pl.Spec.SkipFormat),
				SkipFixImports:  BoolWithDefault(planItem.SkipFixImports, pl.Spec.SkipFixImports),
				Merge:           BoolWithDefault(planItem.Merge, pl.Spec.Merge),
				RelativePackage: planItem.RelativePackage,
				TemplateNames:   planItem.TemplateNames,
				Filename:        planItem.Filename,
				Package:         planItem.Package,
				Imports:         planItem.Imports,
			}

			// Global data
			if plan.Data == nil && len(planItem.Data) <= 0 && pl.Spec.Data != nil {
				planItem.Data = pl.Spec.Data
			}

			// plan item data
			if len(planItem.Data) > 0 && plan.Data == nil {
				plan.Data = planItem.Data[0]
			} else if ext := pl.Spec.Extension; ext != nil {
				// extension data
				// For empty data, try to build data from extension metadata
				data, err := BuildSingletonData(ext.Framework, ext.Data)
				if err != nil {
					return err
				}
				if len(data) > 0 {
					plan.Data = data[0]
				}
			}

			pl.Plans = append(pl.Plans, plan)
		case types.REPEATABLE.String():
			plan := &RepeatablePlan{
				SingletonPlan: SingletonPlan{
					Suffix:          StringWithDefault(planItem.Suffix, pl.Spec.Suffix),
					TemplateDir:     StringWithDefault(planItem.TemplateDir, pl.Spec.TemplateDir),
					BaseDir:         StringWithDefault(planItem.BaseDir, pl.Spec.BaseDir),
					BasePackage:     StringWithDefault(planItem.BasePackage, pl.Spec.BasePackage),
					SkipExists:      BoolWithDefault(planItem.SkipExists, pl.Spec.SkipExists),
					SkipFormat:      BoolWithDefault(planItem.SkipFormat, pl.Spec.SkipFormat),
					SkipFixImports:  BoolWithDefault(planItem.SkipFixImports, pl.Spec.SkipFixImports),
					Merge:           BoolWithDefault(planItem.Merge, pl.Spec.Merge),
					RelativePackage: planItem.RelativePackage,
					Filename:        planItem.Filename,
					Package:         planItem.Package,
					TemplateNames:   planItem.TemplateNames,
					Imports:         planItem.Imports,
				},
			}

			// Global data
			if plan.Data == nil && len(planItem.Data) <= 0 && pl.Spec.Data != nil {
				planItem.Data = pl.Spec.Data
			}

			// plan item data
			if len(planItem.Data) > 0 && plan.Data == nil {
				plan.Data = planItem.Data
			}

			// extension data
			if ext := pl.Spec.Extension; ext != nil && len(planItem.Data) <= 0 {
				data, err := BuildRepeatableData(ext.Framework, ext.Data)
				if err != nil {
					return err
				}

				plan.Data = data
			}
			pl.Plans = append(pl.Plans, plan)
		default:
			return errors.New(fmt.Sprintf("Invalid plan type %s, PlanType must be either SINGLETON or REPEATABLE"))
		}
	}

	return
}

func BuildSingletonData(framework string, data json.RawMessage) (o []interface{}, err error) {
	switch framework {
	case "sqlboiler":
		opts := orm.DataOpts{Framework: framework}
		if err = json.Unmarshal(data, &opts); err != nil {
			return nil, err
		}
		opts.Driver, err = orm.DeadDriverByName(opts.DriverName)
		if err != nil {
			return nil, err
		}

		return orm.BuildSingletonData(opts)
	default:
		return nil, errors.Errorf("Framework <%s> is not supported", framework)
	}
	return
}

func BuildRepeatableData(framework string, data json.RawMessage) (o []interface{}, err error) {
	switch framework {
	case "sqlboiler":
		opts := orm.DataOpts{}
		if err = json.Unmarshal(data, &opts); err != nil {
			return nil, err
		}

		opts.Driver, err = orm.DeadDriverByName(opts.DriverName)
		if err != nil {
			return nil, err
		}
		return orm.BuildRepeatableData(opts)
	default:
		return nil, errors.Errorf("Framework <%s> is not supported", framework)
	}
	return
}
