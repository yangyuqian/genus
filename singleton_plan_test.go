package genus

import "testing"
import "github.com/yangyuqian/genus/types"

func TestSingletonPlan_Type(t *testing.T) {
	tests := []struct {
		name string
		p    *SingletonPlan
		want types.PlanType
	}{
		{"OK", nil, types.SINGLETON},
	}
	for _, tt := range tests {
		p := tt.p
		if got := p.Type(); got != tt.want {
			t.Errorf("%q. SingletonPlan.Type() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSingletonPlan_validate(t *testing.T) {
	type fields struct {
		BaseDir         string
		BasePackage     string
		TemplateDir     string
		TemplateNames   []string
		RelativePackage string
		Imports         Imports
		Data            interface{}
		TemplateGroup   *TemplateGroup
		Repo            *Repo
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"OK", fields{BaseDir: "./", RelativePackage: "a/b", TemplateDir: "./", TemplateNames: []string{"a/b"}}, false},
		{"KO - BaseDir not set", fields{RelativePackage: "a/b", TemplateNames: []string{"a/b"}}, true},
	}
	for _, tt := range tests {
		p := &SingletonPlan{
			BaseDir:         tt.fields.BaseDir,
			BasePackage:     tt.fields.BasePackage,
			TemplateDir:     tt.fields.TemplateDir,
			TemplateNames:   tt.fields.TemplateNames,
			RelativePackage: tt.fields.RelativePackage,
			Imports:         tt.fields.Imports,
			Data:            tt.fields.Data,
			TemplateGroup:   tt.fields.TemplateGroup,
			Repo:            tt.fields.Repo,
		}
		if err := p.validate(); (err != nil) != tt.wantErr {
			t.Errorf("%q. SingletonPlan.validate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestSingletonPlan_init(t *testing.T) {
	type fields struct {
		Suffix          string
		BaseDir         string
		TemplateDir     string
		BasePackage     string
		RelativePackage string
		TemplateNames   []string
		Imports         Imports
		Data            interface{}
		TemplateGroup   *TemplateGroup
		Repo            *Repo
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"OK", fields{Suffix: ".tpl", TemplateDir: "./testdata/repo/success", BaseDir: "./_test", RelativePackage: "a/b", TemplateNames: []string{"t1"}}, false},
		{"KO - BaseDir not set", fields{Suffix: ".tpl", TemplateDir: "./testdata", RelativePackage: "a/b", TemplateNames: []string{"t1"}}, true},
		{"KO - TemplateNames not set", fields{RelativePackage: "a/b", BaseDir: "./"}, true},
	}
	for _, tt := range tests {
		p := &SingletonPlan{
			Suffix:          tt.fields.Suffix,
			BaseDir:         tt.fields.BaseDir,
			BasePackage:     tt.fields.BasePackage,
			TemplateDir:     tt.fields.TemplateDir,
			TemplateNames:   tt.fields.TemplateNames,
			RelativePackage: tt.fields.RelativePackage,
			Imports:         tt.fields.Imports,
			Data:            tt.fields.Data,
			TemplateGroup:   tt.fields.TemplateGroup,
			Repo:            tt.fields.Repo,
		}
		if err := p.init(); (err != nil) != tt.wantErr {
			t.Errorf("%q. SingletonPlan.init() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestSingletonPlan_Render(t *testing.T) {
	type fields struct {
		Suffix          string
		TemplateDir     string
		BaseDir         string
		BasePackage     string
		RelativePackage string
		TemplateNames   []string
		Imports         Imports
		Data            interface{}
		SkipExists      bool
		SkipFormat      bool
		SkipFixImports  bool
		Repo            *Repo
		TemplateGroup   *TemplateGroup
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"KO - plan not set", fields{}, true},
		{"OK", fields{
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
			Data: map[string]interface{}{"Package": "p1"},
		}, false},
	}
	for _, tt := range tests {
		p := &SingletonPlan{
			Suffix:          tt.fields.Suffix,
			TemplateDir:     tt.fields.TemplateDir,
			BaseDir:         tt.fields.BaseDir,
			BasePackage:     tt.fields.BasePackage,
			RelativePackage: tt.fields.RelativePackage,
			TemplateNames:   tt.fields.TemplateNames,
			Imports:         tt.fields.Imports,
			Data:            tt.fields.Data,
			SkipExists:      tt.fields.SkipExists,
			SkipFormat:      tt.fields.SkipFormat,
			SkipFixImports:  tt.fields.SkipFixImports,
			Repo:            tt.fields.Repo,
			TemplateGroup:   tt.fields.TemplateGroup,
		}
		if err := p.Render(); (err != nil) != tt.wantErr {
			t.Errorf("%q. SingletonPlan.Render() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
