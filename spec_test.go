package genus

import "testing"

func TestSpec_Save(t *testing.T) {
	type fields struct {
		Suffix         string
		TemplateDir    string
		BaseDir        string
		BasePackage    string
		SkipExists     bool
		SkipFormat     bool
		SkipFixImports bool
		Merge          bool
		PlanItems      []*PlanItem
		Extension      *SpecExtension
	}
	type args struct {
		baseDir string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"OK", fields{}, args{"./_test"}, false},
	}
	for _, tt := range tests {
		spec := &Spec{
			Suffix:         tt.fields.Suffix,
			TemplateDir:    tt.fields.TemplateDir,
			BaseDir:        tt.fields.BaseDir,
			BasePackage:    tt.fields.BasePackage,
			SkipExists:     tt.fields.SkipExists,
			SkipFormat:     tt.fields.SkipFormat,
			SkipFixImports: tt.fields.SkipFixImports,
			Merge:          tt.fields.Merge,
			PlanItems:      tt.fields.PlanItems,
			Extension:      tt.fields.Extension,
		}
		if err := spec.Save(tt.args.baseDir); (err != nil) != tt.wantErr {
			t.Errorf("%q. Spec.Save() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
