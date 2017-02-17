{{- $dot := .}}
// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestInsert(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Insert)
  {{end -}}
  {{- end -}}
}

func TestAll(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}All)
  {{end -}}
  {{- end -}}
}
