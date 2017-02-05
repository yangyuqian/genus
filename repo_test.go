package genus

import "testing"

func TestRepo_Load(t *testing.T) {
	type fields struct {
		Suffix        string
		TemplateDir   string
		Templates     []*Template
		templateNames []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"OK", fields{Suffix: ".tpl", TemplateDir: "./testdata/repo/success"}, false},
	}
	for _, tt := range tests {
		r := &Repo{
			Suffix:        tt.fields.Suffix,
			TemplateDir:   tt.fields.TemplateDir,
			Templates:     tt.fields.Templates,
			templateNames: tt.fields.templateNames,
		}
		if err := r.Load(); (err != nil) != tt.wantErr {
			t.Errorf("%q. Repo.Load() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
