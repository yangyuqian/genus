package genus

import (
	"reflect"
	"testing"
)

func TestTemplate_loadFile(t *testing.T) {
	type fields struct {
		Name        string
		Source      string
		TargetDir   string
		Filename    string
		SkipExists  bool
		SkipFormat  bool
		rawTemplate []byte
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []byte
		wantErr  bool
	}{
		{"OK", fields{
			Name:      "t1",
			Source:    "./testdata/template/t1.tpl",
			TargetDir: "_test",
			Filename:  "t1.go",
		}, []byte("package main\nfunc T1(){}\n"), false},
		{"KO", fields{
			Name:      "t1",
			Source:    "./testdata/template/not-exist-t1.tpl",
			TargetDir: "_test",
			Filename:  "t1.go",
		}, nil, true},
		{"KO", fields{
			Name:      "t1",
			TargetDir: "_test",
			Filename:  "t1.go",
		}, nil, true},
	}
	for _, tt := range tests {
		tmpl := &Template{
			Name:        tt.fields.Name,
			Source:      tt.fields.Source,
			TargetDir:   tt.fields.TargetDir,
			Filename:    tt.fields.Filename,
			SkipExists:  tt.fields.SkipExists,
			SkipFormat:  tt.fields.SkipFormat,
			rawTemplate: tt.fields.rawTemplate,
		}
		gotData, err := tmpl.loadFile()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Template.loadFile() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotData, tt.wantData) {
			t.Errorf("%q. Template.loadFile() = %v, want %v", tt.name, string(gotData), string(tt.wantData))
		}
	}
}

func TestTemplate_render(t *testing.T) {
	type fields struct {
		Name        string
		Source      string
		TargetDir   string
		Filename    string
		SkipExists  bool
		SkipFormat  bool
		rawTemplate []byte
	}
	type args struct {
		context interface{}
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData []byte
		wantErr  bool
	}{
		{"OK", fields{
			Name:        "t1",
			Source:      "./testdata/template/t1.tpl",
			TargetDir:   "_test",
			Filename:    "t1.go",
			rawTemplate: []byte("package {{ .Package }}"),
		}, args{map[string]string{"Package": "main"}}, []byte("package main"), false},
		{"KO", fields{
			Name:        "t1",
			Source:      "./testdata/template/t1.tpl",
			TargetDir:   "_test",
			Filename:    "t1.go",
			rawTemplate: []byte("package {{ .Package"),
		}, args{map[string]string{"Package": "main"}}, nil, true},
	}
	for _, tt := range tests {
		tmpl := &Template{
			Name:        tt.fields.Name,
			Source:      tt.fields.Source,
			TargetDir:   tt.fields.TargetDir,
			Filename:    tt.fields.Filename,
			SkipExists:  tt.fields.SkipExists,
			SkipFormat:  tt.fields.SkipFormat,
			rawTemplate: tt.fields.rawTemplate,
		}
		gotData, err := tmpl.render(tt.args.context)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Template.render() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotData, tt.wantData) {
			t.Errorf("%q. Template.render() = %v, want %v", tt.name, string(gotData), string(tt.wantData))
		}
	}
}
