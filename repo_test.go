package genus

import (
	"reflect"
	"testing"
)

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
		{"KO - directory not exist", fields{Suffix: ".tpl", TemplateDir: "./testdata/repo/not-exist"}, true},
		{"KO - template dir not set", fields{Suffix: ".tpl"}, true},
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

func TestRepo_BuildGroup(t *testing.T) {
	type fields struct {
		Suffix        string
		TemplateDir   string
		Templates     []*Template
		templateNames []string
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantTg  *TemplateGroup
		wantErr bool
	}{
		{"OK - build template group with names", fields{
			Suffix:      ".tpl",
			TemplateDir: "./testdata/repo/success"},
			args{[]string{"testdata/repo/success/t1"}},
			&TemplateGroup{
				Templates: []*Template{
					&Template{Name: "testdata/repo/success/t1",
						Source: "testdata/repo/success/t1.tpl"}}},
			false},
		{"KO - template not found", fields{
			Suffix:      ".tpl",
			TemplateDir: "./testdata/repo/success"},
			args{[]string{"testdata/repo/success/non-exist"}},
			nil, true},
	}
	for _, tt := range tests {
		r := &Repo{
			Suffix:        tt.fields.Suffix,
			TemplateDir:   tt.fields.TemplateDir,
			Templates:     tt.fields.Templates,
			templateNames: tt.fields.templateNames,
		}
		if err := r.Load(); err != nil {
			t.Errorf("Load template err %+v", err)
		}

		gotTg, err := r.BuildGroup(tt.args.names...)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Repo.BuildGroup() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotTg, tt.wantTg) {
			t.Errorf("%q. Repo.BuildGroup() = %v, want %v", tt.name, gotTg, tt.wantTg)
		}
	}
}

func TestRepo_Lookup(t *testing.T) {
	type fields struct {
		Suffix        string
		TemplateDir   string
		Templates     []*Template
		templateNames []string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantT   *Template
		wantErr bool
	}{
		{"OK", fields{Suffix: ".tpl", TemplateDir: "./testdata/repo/success"}, args{"testdata/repo/success/t1"}, &Template{Name: "testdata/repo/success/t1", Source: "testdata/repo/success/t1.tpl"}, false},
		{"KO - template not found", fields{Suffix: ".tpl", TemplateDir: "./testdata/repo/success"}, args{"testdata/repo/success/non-exist"}, nil, true},
	}
	for _, tt := range tests {
		r := &Repo{
			Suffix:        tt.fields.Suffix,
			TemplateDir:   tt.fields.TemplateDir,
			Templates:     tt.fields.Templates,
			templateNames: tt.fields.templateNames,
		}

		if err := r.Load(); err != nil {
			t.Errorf("Load template err %+v", err)
		}

		gotT, err := r.Lookup(tt.args.name)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Repo.Lookup() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotT, tt.wantT) {
			t.Errorf("%q. Repo.Lookup() = %v, want %v", tt.name, gotT, tt.wantT)
		}
	}
}
