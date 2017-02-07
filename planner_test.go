package genus

import "testing"

func TestPackagePlanner_Plan(t *testing.T) {
	type fields struct {
		RawSpec []byte
		Spec    *Spec
		Plans   []Plan
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"OK", fields{RawSpec: []byte(`{
			"PlanItems": [
				{
					"PlanType": "SINGLETON",
					"Suffix": ".tpl",
					"TemplateDir": "./testdata/plan/success",
					"BaseDir": "./_test",
					"TemplateNames": ["testdata/plan/success/p1/t1"],
					"Data": [
						{
							"Message": "Hello World"
						}
					]
				},
				{
					"PlanType": "SINGLETON",
					"Suffix": ".tpl",
					"TemplateDir": "./testdata/plan/success",
					"Filename": "{{ .Package }}111.go",
					"BaseDir": "./_test",
					"TemplateNames": ["testdata/plan/success/p1/t1"],
					"Data": [
						{
							"Message": "Hello World"
						}
					]
				},
				{
					"PlanType": "REPEATABLE",
					"Suffix": ".tpl",
					"TemplateDir": "./testdata/plan/success",
					"BaseDir": "./_test",
					"Filename": "hello_{{ .Name }}.go",
					"TemplateNames": ["testdata/plan/success/p1/t1"],
					"Data": [
						{
							"Name": "world1",
							"Message": "Hello World"
						}
					]
				}
			]
		}`)}, false},
	}
	for _, tt := range tests {
		pl := &PackagePlanner{
			RawSpec: tt.fields.RawSpec,
			Spec:    tt.fields.Spec,
			Plans:   tt.fields.Plans,
		}
		if err := pl.Plan(); (err != nil) != tt.wantErr {
			t.Errorf("%q. PackagePlanner.Plan() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestPackagePlanner_Perform(t *testing.T) {
	type fields struct {
		RawSpec []byte
		Spec    *Spec
		Plans   []Plan
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"OK", fields{RawSpec: []byte(`{
			"PlanItems": [
				{
					"PlanType": "SINGLETON",
					"Suffix": ".tpl",
					"TemplateDir": "./testdata/plan/success",
					"BaseDir": "./_test",
					"RelativePackage": "p1/p2",
					"TemplateNames": [
						"testdata/plan/success/p1/t1",
						"testdata/plan/success/p1/t3"
					],
					"Data": [
						{
							"Message": "Hello World"
						}
					]
				},
				{
					"PlanType": "SINGLETON",
					"Suffix": ".tpl",
					"Package": "main",
					"TemplateDir": "./testdata/plan/success",
					"BaseDir": "./_test",
					"RelativePackage": "p4/p5",
					"BasePackage": "github.com/user/repo/path1/path2",
					"TemplateNames": [
						"testdata/plan/success/p1/t1",
						"testdata/plan/success/p1/t3",
						"testdata/plan/success/p1/t5"
					],
					"Imports": {
						"./p3": ""
					},
					"Data": [
						{
							"Message": "Hello World"
						}
					]
				},
				{
					"PlanType": "REPEATABLE",
					"Suffix": ".tpl",
					"TemplateDir": "./testdata/plan/success",
					"BaseDir": "./_test",
					"RelativePackage": "p1/p2",
					"Filename": "hello_{{- .Name -}}.go",
					"TemplateNames": ["testdata/plan/success/p1/t1"],
					"Data": [
						{
							"Name": "world1",
							"Message": "Hello World 1"
						},
						{
							"Name": "world2",
							"Message": "Hello World 2"
						}
					]
				}
			]
		}`)}, false},
	}
	for _, tt := range tests {
		pl := &PackagePlanner{
			RawSpec: tt.fields.RawSpec,
			Spec:    tt.fields.Spec,
			Plans:   tt.fields.Plans,
		}
		if err := pl.Plan(); err != nil {
			t.Errorf("%q. PackagePlanner.Plan() error = %v", tt.name, err)
		}

		if err := pl.Perform(); (err != nil) != tt.wantErr {
			t.Errorf("%q. PackagePlanner.Perform() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
