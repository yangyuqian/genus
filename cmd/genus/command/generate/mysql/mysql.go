package mysql

import (
	"go/build"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/yangyuqian/genus"
	"github.com/yangyuqian/genus/generator/orm"
	"github.com/yangyuqian/genus/types"
	cli "gopkg.in/urfave/cli.v1"
)

var basePackage = "github.com/yangyuqian/genus/generator/orm"

func CollectTables(ctx *cli.Context) (tables orm.Tables, err error) {
	s := orm.NewSchema(ctx.String("username"), ctx.String("password"),
		ctx.String("database"), ctx.String("host"), ctx.Int("port"),
		ctx.String("sslmode"))
	defer s.Close()

	return s.CollectTables(ctx.StringSlice("include_table"), ctx.StringSlice("exclude_table"))
}

func BaseDir(base, framework string) (basePath string, err error) {
	if base != "" {
		return base, nil
	}

	base = filepath.Join(basePackage, framework)
	p, err := build.Default.Import(base, "", build.FindOnly)
	if err != nil {
		return "", err
	}

	return p.Dir, nil
}

func InitGPS(ctx *cli.Context, singletonData, repeatableData []interface{}) (spec *genus.Spec, err error) {
	tmplDir, err := BaseDir(ctx.String("template-dir"), ctx.String("framework"))
	if err != nil {
		return nil, err
	}

	switch ctx.String("framework") {
	case "sqlboiler":
		return &genus.Spec{
			Suffix:      ctx.String("template-suffix"),
			BaseDir:     ctx.String("base"),
			TemplateDir: tmplDir,
			PlanItems: []*genus.PlanItem{
				GetPlanItem(
					&planItemOpts{
						PlanType:    types.SINGLETON,
						Suffix:      ctx.String("template-suffix"),
						RelativePkg: ctx.String("relative-pkg"),
						BaseDir:     ctx.String("base"),
						Filename:    "",
						TmplNames: []string{
							"templates/singleton/boil_queries",
							"templates/singleton/boil_types",
						},
						Data:    singletonData,
						Imports: ormImports,
					},
				),

				GetPlanItem(
					&planItemOpts{
						PlanType:    types.SINGLETON,
						Suffix:      ctx.String("template-suffix"),
						RelativePkg: ctx.String("relative-pkg"),
						BaseDir:     ctx.String("base"),
						Filename:    "mysql_main_test.go",
						TmplNames: []string{
							"templates_test/main_test/mysql_main",
						},
						Data:    singletonData,
						Imports: testImports,
					},
				),

				GetPlanItem(
					&planItemOpts{
						PlanType:    types.SINGLETON,
						Suffix:      ctx.String("template-suffix"),
						RelativePkg: ctx.String("relative-pkg"),
						BaseDir:     ctx.String("base"),
						TmplNames: []string{
							"templates_test/singleton/boil_main_test",
							"templates_test/singleton/boil_queries_test",
							"templates_test/singleton/boil_suites_test",
						},
						Data:    singletonData,
						Merge:   false,
						Imports: testImports,
					},
				),

				GetPlanItem(
					&planItemOpts{
						PlanType:    types.REPEATABLE,
						Suffix:      ctx.String("template-suffix"),
						RelativePkg: ctx.String("relative-pkg"),
						BaseDir:     ctx.String("base"),
						Filename:    "{{- with .Table }}{{.Name}}{{ end -}}.go",
						TmplNames: []string{
							"templates/00_struct",
							"templates/01_types",
							"templates/02_hooks",
							"templates/03_finishers",
							"templates/04_relationship_to_one",
							"templates/05_relationship_one_to_one",
							"templates/06_relationship_to_many",
							"templates/07_relationship_to_one_eager",
							"templates/08_relationship_one_to_one_eager",
							"templates/09_relationship_to_many_eager",
							"templates/10_relationship_to_one_setops",
							"templates/11_relationship_one_to_one_setops",
							"templates/12_relationship_to_many_setops",
							"templates/13_all",
							"templates/14_find",
							"templates/15_insert",
							"templates/16_update",
							"templates/17_upsert",
							"templates/18_delete",
							"templates/19_reload",
							"templates/20_exists",
						},
						Data:    repeatableData,
						Merge:   true,
						Imports: ormImports,
					},
				),

				GetPlanItem(
					&planItemOpts{
						PlanType:    types.REPEATABLE,
						Suffix:      ctx.String("template-suffix"),
						RelativePkg: ctx.String("relative-pkg"),
						BaseDir:     ctx.String("base"),
						Filename:    "{{- with .Table }}{{.Name}}{{ end -}}_test.go",
						TmplNames: []string{
							"templates_test/all",
							"templates_test/delete",
							"templates_test/exists",
							"templates_test/find",
							"templates_test/finishers",
							"templates_test/hooks",
							"templates_test/insert",
							"templates_test/relationship_one_to_one",
							"templates_test/relationship_one_to_one_setops",
							"templates_test/relationship_to_many",
							"templates_test/relationship_to_many_setops",
							"templates_test/relationship_to_one",
							"templates_test/relationship_to_one_setops",
							"templates_test/reload",
							"templates_test/select",
							"templates_test/types",
							"templates_test/update",
							"templates_test/upsert",
						},
						Data:    repeatableData,
						Merge:   true,
						Imports: testImports,
					},
				),
			},
		}, nil
	case "gorm":
		return &genus.Spec{
			Suffix:      ctx.String("template-suffix"),
			BaseDir:     ctx.String("base"),
			TemplateDir: tmplDir,
			Data:        repeatableData,
			PlanItems: []*genus.PlanItem{
				GetPlanItem(
					&planItemOpts{
						PlanType:    types.REPEATABLE,
						Suffix:      ctx.String("template-suffix"),
						RelativePkg: ctx.String("relative-pkg"),
						BaseDir:     ctx.String("base"),
						Filename:    "{{- with .Table }}{{.Name}}{{ end -}}.go",
						TmplNames: []string{
							"templates/struct",
						},
						Merge:   true,
						Imports: ormImports,
					},
				),
			},
		}, nil
	default:
		return nil, errors.Errorf("Can not init Generation Plan Specification <%s>", ctx.String("framework"))
	}

	return
}
