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
		{"OK", fields{Templates: []*Template{
			&Template{Name: "a/b/c.tpl"},
		}, Imports: Imports{"c": "a/b/c"}}, false},
		{"OK", fields{Templates: []*Template{
			&Template{Name: "a/b/c.tpl"},
		}, Imports: Imports{"c": "./c"}}, false},
		{"OK", fields{}, false},
	}
	for _, tt := range tests {
		tg := &TemplateGroup{
			Package:         tt.fields.Package,
			BaseDir:         tt.fields.BaseDir,
			BasePackage:     tt.fields.BasePackage,
			AbsolutePackage: tt.fields.AbosultePackage,
			Imports:         tt.fields.Imports,
			Templates:       tt.fields.Templates,
			SkipFixImports:  tt.fields.SkipFixImports,
		}
		if err := tg.configureTemplates(); (err != nil) != tt.wantErr {
			t.Errorf("%q. TemplateGroup.configureTemplates() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestTemplateGroup_Render(t *testing.T) {
	type fields struct {
		Package         string
		BaseDir         string
		BasePackage     string
		AbosultePackage string
		Imports         Imports
		Templates       []*Template
		SkipFixImports  bool
		SkipExists      bool
		SkipFormat      bool
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"OK", fields{Package: "p1",
			BaseDir: "./_test",
			Imports: Imports{"./p2": "p2"},
			Templates: []*Template{
				&Template{
					Name:   "success/t1",
					Source: "./testdata/repo/success/t1.tpl",
				},
			}}, args{}, false},
		{"OK", fields{Package: "p2",
			BaseDir: "./_test",
			Imports: Imports{"./p3": "p3"},
			Templates: []*Template{
				&Template{
					Name:   "success/t1/t2",
					Source: "./testdata/repo/success/t1/t2.tpl",
				},
			}}, args{}, false},
	}
	for _, tt := range tests {
		tg := &TemplateGroup{
			Package:         tt.fields.Package,
			BaseDir:         tt.fields.BaseDir,
			BasePackage:     tt.fields.BasePackage,
			AbsolutePackage: tt.fields.AbosultePackage,
			Imports:         tt.fields.Imports,
			Templates:       tt.fields.Templates,
			SkipFixImports:  tt.fields.SkipFixImports,
			SkipExists:      tt.fields.SkipExists,
			SkipFormat:      tt.fields.SkipFormat,
		}
		if err := tg.Render(tt.args.data); (err != nil) != tt.wantErr {
			t.Errorf("%q. TemplateGroup.Render() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestTemplateGroup_RenderPartial(t *testing.T) {
	type fields struct {
		Package         string
		BaseDir         string
		BasePackage     string
		RelativePackage string
		AbsolutePackage string
		Filename        string
		Imports         Imports
		Templates       []*Template
		SkipFixImports  bool
		SkipExists      bool
		SkipFormat      bool
		Merge           bool
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"OK", fields{Package: "p1",
			BaseDir:  "./_test",
			Imports:  Imports{"./p2": "p2"},
			Filename: "merged.go",
			Templates: []*Template{
				&Template{
					Name:   "success/t1",
					Source: "./testdata/repo/success/t1.tpl",
				},
				&Template{
					Name:   "success/t1/t2",
					Source: "./testdata/repo/success/t1/t2.tpl",
				},
			}}, args{map[string]interface{}{"Name": "A"}}, false},
	}
	for _, tt := range tests {
		tg := &TemplateGroup{
			Package:         tt.fields.Package,
			BaseDir:         tt.fields.BaseDir,
			BasePackage:     tt.fields.BasePackage,
			RelativePackage: tt.fields.RelativePackage,
			AbsolutePackage: tt.fields.AbsolutePackage,
			Filename:        tt.fields.Filename,
			Imports:         tt.fields.Imports,
			Templates:       tt.fields.Templates,
			SkipFixImports:  tt.fields.SkipFixImports,
			SkipExists:      tt.fields.SkipExists,
			SkipFormat:      tt.fields.SkipFormat,
			Merge:           tt.fields.Merge,
		}
		if err := tg.RenderPartial(tt.args.data); (err != nil) != tt.wantErr {
			t.Errorf("%q. TemplateGroup.RenderPartial() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
