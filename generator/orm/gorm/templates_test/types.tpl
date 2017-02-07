{{- $varNameSingular := .Table.Name | singular | camelCase -}}
var (
    {{$varNameSingular}}DBTypes = map[string]string{
      {{ range $idx, $col := .Table.Columns -}}
      "{{ titleCase $col.Name }}":"{{ $col.DBType }}",
      {{- end }}
    }
	_ = bytes.MinRead
)
