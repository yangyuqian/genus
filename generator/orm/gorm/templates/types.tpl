{{if .Table.IsJoinTable -}}
{{else -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $tableNameSingular := .Table.Name | singular | titleCase -}}

var (
	{{$varNameSingular}}Columns               = []string{ {{ range $idx, $col := .Table.Columns }}"{{ $col.Name }}",{{ end }} }
	{{$varNameSingular}}PrimaryKeyColumns     = []string{ {{ range $idx, $colName := .Table.PKey.Columns }}"{{ $colName }}",{{ end }} }
	{{$varNameSingular}}ColumnsWithDefault    = []string{ {{ range $idx, $col := .Table.Columns }}{{ if not ( eq $col.Default "" ) }}"{{ $col.Name }}",{{ end }}{{ end }} }
)
{{end -}}
