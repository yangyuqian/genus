{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $modelName := .Table.Name | singular | titleCase -}}
{{- $sliceModelName := printf "%sSlice" $modelName -}}
func test{{$tableNamePlural}}All(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	{{$varNameSingular}} := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	tx := GetDB().Debug().Begin()
  defer tx.Rollback()
	if err = {{$varNameSingular}}.Insert(tx); err != nil {
		t.Error(err)
	}

  o, err := {{$tableNamePlural}}(tx)
  if err != nil {
    t.Error(err)
  }

  if len(o) != 1 {
    t.Errorf("Can not get all {{ $varNamePlural }}")
  }
}
