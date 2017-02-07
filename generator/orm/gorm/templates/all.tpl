{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase}}
{{- $tableNameSingular := .Table.Name | singular -}}
{{- $modelName := $tableNameSingular | titleCase -}}
{{- $sliceModelName := printf "%sSlice" $modelName -}}

// {{$tableNamePlural}} retrieves all the records using an executor.
func {{$tableNamePlural}}(db *gorm.DB) (o {{ $sliceModelName }}, err error) {
  o = {{ $sliceModelName }}{}

	if db := db.Find(&o).Scan(&o); db.Error != nil {
		return nil, db.Error
	}

	return o, nil
}
