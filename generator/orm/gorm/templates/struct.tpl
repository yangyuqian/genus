{{- $dot := . -}}
{{- $tableNameSingular := .Table.Name | singular -}}
{{- $modelName := $tableNameSingular | titleCase -}}
{{- $modelNameCamel := $tableNameSingular | camelCase -}}
{{- $sliceModelName := printf "%sSlice" $modelName -}}

// {{ $modelName }} is an object representing the database table.
type {{ $modelName }} struct {
	{{ range $column := .Table.Columns -}}
	{{ titleCase $column.Name }} {{ $column.Type }} `gorm:"column:{{ $column.Name }}" json:"{{ $column.Name }}{{ if $column.Nullable }},omitempty{{ end }}" boil:"{{$column.Name}}"`
  {{ end -}}
}

func (obj *{{ $modelName }}) TableName() string {
  return "{{ .Table.Name }}"
}

type {{ $sliceModelName }} []*{{ $modelName }}
