package genus

import "testing"

func TestTemplateGroup_ensureGOOS(t *testing.T) {
	type fields struct {
		Package        string
		BaseDir        string
		Imports        Imports
		Templates      []*Template
		SkipFixImports bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"OK", fields{}, false},
	}
	for _, tt := range tests {
		tg := &TemplateGroup{
			Package:        tt.fields.Package,
			BaseDir:        tt.fields.BaseDir,
			Imports:        tt.fields.Imports,
			Templates:      tt.fields.Templates,
			SkipFixImports: tt.fields.SkipFixImports,
		}
		if err := tg.ensureGOOS(); (err != nil) != tt.wantErr {
			t.Errorf("%q. TemplateGroup.ensureGOOS() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestTemplateGroup_ensureGopath(t *testing.T) {
	type fields struct {
		Package        string
		BaseDir        string
		BasePackage    string
		Imports        Imports
		Templates      []*Template
		SkipFixImports bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"OK", fields{}, false},
	}
	for _, tt := range tests {
		tg := &TemplateGroup{
			Package:        tt.fields.Package,
			BaseDir:        tt.fields.BaseDir,
			BasePackage:    tt.fields.BasePackage,
			Imports:        tt.fields.Imports,
			Templates:      tt.fields.Templates,
			SkipFixImports: tt.fields.SkipFixImports,
		}
		if err := tg.ensureGopath(); (err != nil) != tt.wantErr {
			t.Errorf("%q. TemplateGroup.ensureGopath() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestTemplateGroup_configureTemplates(t *testing.T) {
	type fields struct {
		Package         string
		BaseDir         string
		BasePackage     string
		AbosultePackage string
		Imports         Imports
		Templates       []*Template
		SkipFixImports  bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"OK", fields{Templates: []*Template{
			&Template{Name: "a/b/c.tpl"},
		}}, false},
		{"OK", fields{}, false},
	}
	for _, tt := range tests {
		tg := &TemplateGroup{
			Package:         tt.fields.Package,
			BaseDir:         tt.fields.BaseDir,
			BasePackage:     tt.fields.BasePackage,
			AbosultePackage: tt.fields.AbosultePackage,
			Imports:         tt.fields.Imports,
			Templates:       tt.fields.Templates,
			SkipFixImports:  tt.fields.SkipFixImports,
		}
		if err := tg.configureTemplates(); (err != nil) != tt.wantErr {
			t.Errorf("%q. TemplateGroup.configureTemplates() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
