{{- $dot := . -}}
{{- $tableNameSingular := .Table.Name | singular -}}
{{- $modelName := $tableNameSingular | titleCase -}}
{{- $modelNameCamel := $tableNameSingular | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func (obj *{{ $modelName }}) Insert(db *gorm.DB, blacklist ...string) (err error) {
  if db == nil {
    errors.Errorf("Can not insert {{ $modelName }} with nil connection")
  }

  if len(blacklist) > 0 {
    db = db.Omit(blacklist...)
  }

  db.Create(obj)

  return db.Error
}

func (obj *{{ $modelName }}) InsertG(blacklist ...string) (err error) {
  if DB == nil {
    errors.Errorf("Can not insert {{ $modelName }} with nil global connection")
  }

  return obj.Insert(DB)
}

