package orm

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/strmangle"
)

type TemplateData struct {
	Tables []bdb.Table
	Table  bdb.Table

	// Controls what names are output
	PkgName string
	Schema  string

	// Controls which code is output (mysql vs postgres ...)
	DriverName      string
	UseLastInsertID bool

	// Turn off auto timestamps or hook generation
	NoHooks          bool
	NoAutoTimestamps bool

	// Tags control which
	Tags []string

	// StringFuncs are usable in templates with stringMap
	StringFuncs map[string]func(string) string `json:"-"`

	// Dialect controls quoting
	Dialect queries.Dialect
	LQ      string
	RQ      string
}

func (t TemplateData) Quotes(s string) string {
	return fmt.Sprintf("%s%s%s", t.LQ, s, t.RQ)
}

func (t TemplateData) SchemaTable(table string) string {
	return strmangle.SchemaTable(t.LQ, t.RQ, t.DriverName, t.Schema, table)
}

var templateStringMappers = map[string]func(string) string{
	// String ops
	"quoteWrap": func(a string) string { return fmt.Sprintf(`"%s"`, a) },

	// Casing
	"titleCase": strmangle.TitleCase,
	"camelCase": strmangle.CamelCase,
}

type DataOpts struct {
	Framework   string
	BaseDir     string
	TemplateDir string
	Schema      string
	DriverName  string
	PkgName     string
	Driver      bdb.Interface `json:"-"`
	Tables      Tables
}

func BuildSingletonData(opts DataOpts) (data []interface{}, err error) {
	switch opts.Framework {
	case "sqlboiler":
		data = append(data, &TemplateData{
			Tables:           opts.Tables,
			Schema:           opts.Schema,
			DriverName:       opts.DriverName,
			UseLastInsertID:  true,
			PkgName:          opts.PkgName,
			NoHooks:          false,
			NoAutoTimestamps: false,
			Dialect: queries.Dialect{
				LQ:                opts.Driver.LeftQuote(),
				RQ:                opts.Driver.RightQuote(),
				IndexPlaceholders: opts.Driver.IndexPlaceholders(),
			},
			LQ: strmangle.QuoteCharacter(opts.Driver.LeftQuote()),
			RQ: strmangle.QuoteCharacter(opts.Driver.RightQuote()),

			StringFuncs: templateStringMappers,
		})
	case "gorm":
		data = append(data, map[string]interface{}{
			"Tables":  opts.Tables,
			"PkgName": opts.PkgName,
		})
	default:
		return nil, errors.Errorf("Can not build singleton data for framework <%s>", opts.Framework)
	}

	return
}
func BuildRepeatableData(opts DataOpts) (data []interface{}, err error) {
	for _, table := range opts.Tables {
		if table.IsJoinTable {
			continue
		}

		switch opts.Framework {
		case "sqlboiler":
			data = append(data, &TemplateData{
				Tables:           opts.Tables,
				Table:            table,
				Schema:           opts.Schema,
				DriverName:       opts.DriverName,
				UseLastInsertID:  true,
				PkgName:          opts.PkgName,
				NoHooks:          false,
				NoAutoTimestamps: false,
				Tags:             nil,
				Dialect: queries.Dialect{
					LQ:                opts.Driver.LeftQuote(),
					RQ:                opts.Driver.RightQuote(),
					IndexPlaceholders: opts.Driver.IndexPlaceholders(),
				},
				LQ: strmangle.QuoteCharacter(opts.Driver.LeftQuote()),
				RQ: strmangle.QuoteCharacter(opts.Driver.RightQuote()),

				StringFuncs: templateStringMappers,
			})
		case "gorm":
			data = append(data, map[string]interface{}{
				"Table":   table,
				"Tables":  opts.Tables,
				"PkgName": opts.PkgName,
			})
		default:
			return nil, errors.Errorf("Can not build repeatable data for framework <%s>", opts.Framework)
		}
	}

	return
}
