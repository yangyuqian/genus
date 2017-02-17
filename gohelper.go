package genus

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/strmangle"
)

// set is to stop duplication from named enums, allowing a template loop
// to keep some state
type once map[string]struct{}

func newOnce() once {
	return make(once)
}

func (o once) Has(s string) bool {
	_, ok := o[s]
	return ok
}

func (o once) Put(s string) bool {
	if _, ok := o[s]; ok {
		return false
	}

	o[s] = struct{}{}
	return true
}

var GoHelperFuncs = template.FuncMap{
	// String ops
	"quoteWrap": func(s string) string { return fmt.Sprintf(`"%s"`, s) },
	"id":        strmangle.Identifier,

	// Pluralization
	"singular": strmangle.Singular,
	"plural":   strmangle.Plural,

	// Casing
	"titleCase": strmangle.TitleCase,
	"camelCase": strmangle.CamelCase,

	// String Slice ops
	"join":               func(sep string, slice []string) string { return strings.Join(slice, sep) },
	"joinSlices":         strmangle.JoinSlices,
	"stringMap":          strmangle.StringMap,
	"prefixStringSlice":  strmangle.PrefixStringSlice,
	"containsAny":        strmangle.ContainsAny,
	"generateTags":       strmangle.GenerateTags,
	"generateIgnoreTags": strmangle.GenerateIgnoreTags,

	// Enum ops
	"parseEnumName":       strmangle.ParseEnumName,
	"parseEnumVals":       strmangle.ParseEnumVals,
	"isEnumNormal":        strmangle.IsEnumNormal,
	"shouldTitleCaseEnum": strmangle.ShouldTitleCaseEnum,
	"onceNew":             newOnce,
	"oncePut":             once.Put,
	"onceHas":             once.Has,

	// String Map ops
	"makeStringMap": strmangle.MakeStringMap,

	// Set operations
	"setInclude": strmangle.SetInclude,

	// Database related mangling
	"whereClause": strmangle.WhereClause,

	// Relationship text helpers
	"txtsFromFKey":     txtsFromFKey,
	"txtsFromOneToOne": txtsFromOneToOne,
	"txtsFromToMany":   txtsFromToMany,

	// dbdrivers ops
	"filterColumnsByDefault": bdb.FilterColumnsByDefault,
	"filterColumnsByEnum":    bdb.FilterColumnsByEnum,
	"sqlColDefinitions":      bdb.SQLColDefinitions,
	"columnNames":            bdb.ColumnNames,
	"columnDBTypes":          bdb.ColumnDBTypes,
	"getTable":               bdb.GetTable,
	"downcase":               strings.ToLower,
}

func StringWithDefault(v, d string) string {
	if v == "" {
		return d
	}

	return v
}

func BoolWithDefault(v, d bool) bool {
	return v || d
}
