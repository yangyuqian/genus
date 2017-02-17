package genus

import "testing"

func TestRepeatablePlan_Render(t *testing.T) {
	type fields struct {
		Data          []interface{}
		SingletonPlan SingletonPlan
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"KO - plan not set", fields{Data: []interface{}{}}, true},
		{"OK", fields{
			Data: []interface{}{map[string]interface{}{"Package": "p1"}},
			SingletonPlan: SingletonPlan{
				Suffix:          ".tpl",
				TemplateDir:     "./testdata/plan/success",
				BaseDir:         "./_test",
				BasePackage:     "github.com/yangyuqian/myapp",
				RelativePackage: "x/y",
				TemplateNames: []string{
					"p1/t1",
					"p1/p2/t2",
					"p1/t3",
				},
			},
		}, false},
	}
	for _, tt := range tests {
		p := &RepeatablePlan{
			Data:          tt.fields.Data,
			SingletonPlan: tt.fields.SingletonPlan,
		}
		if err := p.Render(); (err != nil) != tt.wantErr {
			t.Errorf("%q. RepeatablePlan.Render() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
