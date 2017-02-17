package generate

import (
	"encoding/json"

	"github.com/yangyuqian/genus"
	"github.com/yangyuqian/genus/cmd/genus/command/generate/mysql"
	"github.com/yangyuqian/genus/generator/orm"

	cli "gopkg.in/urfave/cli.v1"
)

var drivers = []string{"mysql", "postgres"}

var usages = map[string]string{
	"host":            "Database host",
	"port":            "Database port",
	"username":        "Database username",
	"password":        "Database password",
	"sslmode":         "Connect to database in sslmode",
	"include_table":   "Tables included in the generated ORM and seeked by their foreign keys",
	"database":        "Database name",
	"pkg":             "Package name of generated ORM",
	"template-suffix": "template suffix",
	"relative-pkg":    "relative package name of generated ORM",
	"base":            "Base directory of custom templates and configurations",
	"mysql":           "Generate ORM for MySQL",
	"postgres":        "Generate ORM for PostgrepSQL",
	"models":          "Parse database schema, generate a ORM",
	"force":           "Update cached schema in .lock file and regenerate the ORM",
	"framework":       "Popular frameworks that supported to be generated",
	"template-dir":    "Location of templates",
	"driver-name":     "Driver name",
}

var ORMCmd = cli.Command{
	Name:   "orm",
	Usage:  "Generate ORM against database schema",
	Action: doGenerate,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host", Value: "localhost", Usage: usages["host"]},
		cli.StringFlag{Name: "username", Value: "root", Usage: usages["username"]},
		cli.StringFlag{Name: "password", Value: "root", Usage: usages["password"]},
		cli.IntFlag{Name: "port", Value: 3306, Usage: usages["port"]},
		cli.StringFlag{Name: "database", Usage: usages["database"]},
		cli.StringFlag{Name: "sslmode", Value: "false", Usage: usages["sslmode"]},
		cli.StringSliceFlag{Name: "include_table", Usage: usages["include_table"]},
		cli.StringFlag{Name: "pkg", Usage: usages["pkg"]},
		cli.StringFlag{Name: "relative-pkg", Usage: usages["relative-pkg"]},
		cli.StringFlag{Name: "template-suffix", Usage: usages["template-suffix"], Value: ".tpl"},
		cli.StringFlag{Name: "base", Usage: usages["base"], Value: ""},
		cli.StringFlag{Name: "template-dir", Usage: usages["template-dir"], Value: ""},
		cli.StringFlag{Name: "framework", Usage: usages["framework"], Value: ""},
		cli.StringFlag{Name: "driver-name", Usage: usages[""], Value: "mysql", Hidden: true},
		cli.BoolFlag{Name: "mysql", Usage: usages["mysql"]},
		cli.BoolFlag{Name: "postgres", Usage: usages["postgres"], Hidden: true},
		cli.BoolFlag{Name: "force", Usage: usages["force"]},
	},
}

func doGenerateORM(ctx *cli.Context) (err error) {
	tables, err := mysql.CollectTables(ctx)
	if err != nil {
		return err
	}

	drv, err := orm.DeadDriverByName(ctx.String("driver-name"))
	if err != nil {
		return err
	}

	dataOpts := orm.DataOpts{
		Framework:   ctx.String("framework"),
		BaseDir:     ctx.String("base"),
		TemplateDir: ctx.String("template-dir"),
		Schema:      ctx.String("database"),
		DriverName:  ctx.String("driver-name"),
		PkgName:     ctx.String("pkg"),
		Driver:      drv,
		Tables:      tables,
	}

	singletonData, err := orm.BuildSingletonData(dataOpts)
	if err != nil {
		return err
	}

	repeatableData, err := orm.BuildRepeatableData(dataOpts)
	if err != nil {
		return err
	}

	spec, err := mysql.InitGPS(ctx, singletonData, repeatableData)
	if err != nil {
		return err
	}

	specErr := spec.Save(ctx.String("base"))
	if specErr != nil {
		return specErr
	}

	planner := genus.PackagePlanner{Spec: spec}

	if err := planner.Plan(); err != nil {
		return err
	}

	if err := planner.Perform(); err != nil {
		return err
	}
	return
}

func doGenerateSqlboiler(ctx *cli.Context) (err error) {
	tables, err := mysql.CollectTables(ctx)
	if err != nil {
		return err
	}

	drv, err := orm.DeadDriverByName(ctx.String("driver-name"))
	if err != nil {
		return err
	}

	dataOpts := orm.DataOpts{
		Framework:   ctx.String("framework"),
		BaseDir:     ctx.String("base"),
		TemplateDir: ctx.String("template-dir"),
		Schema:      ctx.String("database"),
		DriverName:  ctx.String("driver-name"),
		PkgName:     ctx.String("pkg"),
		Driver:      drv,
		Tables:      tables,
	}

	singletonData, err := orm.BuildSingletonData(dataOpts)
	if err != nil {
		return err
	}

	repeatableData, err := orm.BuildRepeatableData(dataOpts)
	if err != nil {
		return err
	}

	blankSpec, err := mysql.InitGPS(ctx, nil, nil)
	if err != nil {
		return err
	}

	rawOpts, err := json.Marshal(dataOpts)
	if err != nil {
		return err
	}

	blankSpec.Extension = &genus.SpecExtension{Framework: ctx.String("framework"), Data: json.RawMessage(rawOpts)}
	if berr := blankSpec.Save(ctx.String("base")); berr != nil {
		return berr
	}

	spec, err := mysql.InitGPS(ctx, singletonData, repeatableData)
	if err != nil {
		return err
	}

	planner := genus.PackagePlanner{Spec: spec}

	if err := planner.Plan(); err != nil {
		return err
	}

	if err := planner.Perform(); err != nil {
		return err
	}

	return
}
