package mysql

import (
	"github.com/yangyuqian/genus"
	"github.com/yangyuqian/genus/types"
)

type planItemOpts struct {
	PlanType    types.PlanType
	TmplDir     string
	Suffix      string
	RelativePkg string
	BaseDir     string
	Filename    string
	TmplNames   []string
	Data        []interface{}
	Merge       bool
	Imports     genus.Imports
}

func GetPlanItem(opts *planItemOpts) (item *genus.PlanItem) {
	return &genus.PlanItem{
		PlanType:        opts.PlanType.String(),
		Suffix:          opts.Suffix,
		TemplateDir:     opts.TmplDir,
		RelativePackage: opts.RelativePkg,
		BaseDir:         opts.BaseDir,
		TemplateNames:   opts.TmplNames,
		Filename:        opts.Filename,
		Imports:         opts.Imports,
		Data:            opts.Data,
		Merge:           opts.Merge,
	}
}
