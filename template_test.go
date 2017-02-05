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
		}, []byte("func T1(){}\n"), false},
		{"KO - missing file", fields{
			Name:      "t1",
			Source:    "./testdata/template/not-exist-t1.tpl",
			TargetDir: "_test",
			Filename:  "t1.go",
		}, nil, true},
		{"KO - blank template data", fields{
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
			rawTemplate: []byte("type {{ .Name }} struct{}"),
		}, args{map[string]interface{}{"Package": "p1", "Data": map[string]interface{}{"Name": "A"}}}, []byte("package p1\n\n\ntype A struct{}"), false},
		{"KO - bad syntax in template", fields{
			rawTemplate: []byte("type {{ .Name"),
		}, args{map[string]string{"Name": "A"}}, nil, true},
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

func TestTemplate_SetRawTemplate(t *testing.T) {
	type fields struct {
		Name        string
		Source      string
		TargetDir   string
		Filename    string
		SkipExists  bool
		SkipFormat  bool
		rawTemplate []byte
		rawResult   []byte
	}
	type args struct {
		raw []byte
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData []byte
	}{
		{"OK", fields{}, args{[]byte("abc")}, []byte("abc")},
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
			rawResult:   tt.fields.rawResult,
		}
		if gotData := tmpl.SetRawTemplate(tt.args.raw); !reflect.DeepEqual(gotData, tt.wantData) {
			t.Errorf("%q. Template.SetRawTemplate() = %v, want %v", tt.name, gotData, tt.wantData)
		}
	}
}

func TestTemplate_load(t *testing.T) {
	type fields struct {
		Name        string
		Source      string
		TargetDir   string
		Filename    string
		SkipExists  bool
		SkipFormat  bool
		rawTemplate []byte
		rawResult   []byte
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
		}, []byte("func T1(){}\n"), false},
		{"OK", fields{
			Name:        "t1",
			Source:      "./testdata/template/t1.tpl",
			TargetDir:   "_test",
			Filename:    "t1.go",
			rawTemplate: []byte("abc"),
		}, []byte("abc"), false},
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
			rawResult:   tt.fields.rawResult,
		}
		gotData, err := tmpl.load()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Template.load() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotData, tt.wantData) {
			t.Errorf("%q. Template.load() = %v, want %v", tt.name, gotData, tt.wantData)
		}
	}
}

func TestTemplate_Render(t *testing.T) {
	type fields struct {
		Name        string
		Source      string
		TargetDir   string
		Filename    string
		SkipExists  bool
		SkipFormat  bool
		rawTemplate []byte
		rawResult   []byte
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult []byte
		wantErr    bool
	}{
		{"OK - render template from file", fields{
			Name:      "t3",
			Source:    "./testdata/template/t3.tpl",
			TargetDir: "_test",
			Filename:  "t3.go",
		}, args{map[string]interface{}{"Package": "p1", "Data": map[string]interface{}{"Name": "A"}}}, []byte("package p1\n\ntype A struct{}\n"), false},
		{"OK - render template from raw bytes", fields{
			rawTemplate: []byte("type {{ .Name }} struct{}"),
		}, args{map[string]interface{}{"Package": "p1", "Data": map[string]interface{}{"Name": "A"}}}, []byte("package p1\n\ntype A struct{}\n"), false},
		{"OK - set raw template with imports", fields{
			Name:      "t5",
			Source:    "./testdata/template/t5.tpl",
			TargetDir: "_test",
			Filename:  "t5.go",
		}, args{
			map[string]interface{}{
				"Package": "p1",
				"Data": map[string]interface{}{
					"Name":    "A",
					"Imports": map[string]string{"": "p51"}}},
		},
			[]byte("package p1\n\ntype A struct{}\n"), false},
		{"KO - Source not set", fields{
			Name:      "t3",
			TargetDir: "_test",
			Filename:  "t3.go",
		}, args{map[string]string{"Package": "p1"}}, nil, true},
		{"KO - bad syntax in rawTemplate", fields{
			rawTemplate: []byte("package {{ .Package"),
		}, args{map[string]string{"Package": "p1"}}, nil, true},
		{"KO - bad syntax in template file", fields{
			Name:      "t4",
			Source:    "./testdata/template/t4.tpl",
			TargetDir: "_test",
			Filename:  "t4.go",
		}, args{map[string]string{"Package": "p1"}}, nil, true},
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
			rawResult:   tt.fields.rawResult,
		}
		gotResult, err := tmpl.Render(tt.args.data)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Template.Render() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotResult, tt.wantResult) {
			t.Errorf("%q. Template.Render() = %v, want %v", tt.name, string(gotResult), string(tt.wantResult))
		}
	}
}

func TestTemplate_format(t *testing.T) {
	type fields struct {
		Name        string
		Source      string
		TargetDir   string
		Filename    string
		SkipExists  bool
		SkipFormat  bool
		rawTemplate []byte
		rawResult   []byte
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []byte
		wantErr  bool
	}{
		{"OK", fields{
			rawResult: []byte("package main\nfunc T1(){}\n"),
		}, []byte("package main\n\nfunc T1() {}\n"), false},
		{"KO - bad syntax", fields{
			rawResult: []byte("xxx main\nfunc T1(){}\n"),
		}, nil, true},
		{"KO - bad syntax but skip format", fields{
			SkipFormat: true,
			rawResult:  []byte("xxx main\nfunc T1(){}\n"),
		}, []byte("xxx main\nfunc T1(){}\n"), false},
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
			rawResult:   tt.fields.rawResult,
		}
		gotData, err := tmpl.format()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Template.format() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotData, tt.wantData) {
			t.Errorf("%q. Template.format() = %v, want %v", tt.name, gotData, tt.wantData)
		}
	}
}

func TestTemplate_write(t *testing.T) {
	type fields struct {
		Name        string
		Source      string
		TargetDir   string
		Filename    string
		SkipExists  bool
		SkipFormat  bool
		rawTemplate []byte
		rawResult   []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"OK", fields{
			TargetDir: "./_test",
			Filename:  "t1.go",
			rawResult: []byte("abc"),
		}, false},
		{"OK", fields{
			TargetDir:  "./_test",
			Filename:   "t1.go",
			SkipExists: true,
			rawResult:  []byte("abcd"),
		}, false},
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
			rawResult:   tt.fields.rawResult,
		}
		if err := tmpl.write(); (err != nil) != tt.wantErr {
			t.Errorf("%q. Template.write() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestTemplate_fixImports(t *testing.T) {
	type fields struct {
		Name           string
		Source         string
		TargetDir      string
		Filename       string
		SkipExists     bool
		SkipFormat     bool
		SkipFixImports bool
		header         []byte
		rawTemplate    []byte
		rawResult      []byte
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []byte
		wantErr  bool
	}{
		{"OK", fields{
			rawResult: []byte("package p1\nimport \"p2\""),
			TargetDir: "./_test/t_import",
			Filename:  "t1.go",
		}, []byte("package p1\n"), false},
		{"OK - skip fix imports", fields{
			SkipFixImports: true,
			rawResult:      []byte("package p1\nimport \"p2\""),
			TargetDir:      "./_test/t_import",
			Filename:       "t1.go",
		}, []byte("package p1\nimport \"p2\""), false},
		{"KO - incomplete source code", fields{
			rawResult: []byte("import \"p2\""),
			TargetDir: "./_test/t_import",
			Filename:  "t1.go",
		}, nil, true},
	}
	for _, tt := range tests {
		tmpl := &Template{
			Name:           tt.fields.Name,
			Source:         tt.fields.Source,
			TargetDir:      tt.fields.TargetDir,
			Filename:       tt.fields.Filename,
			SkipFixImports: tt.fields.SkipFixImports,
			SkipExists:     tt.fields.SkipExists,
			SkipFormat:     tt.fields.SkipFormat,
			header:         tt.fields.header,
			rawTemplate:    tt.fields.rawTemplate,
			rawResult:      tt.fields.rawResult,
		}
		gotData, err := tmpl.fixImports()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Template.fixImports() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotData, tt.wantData) {
			t.Errorf("%q. Template.fixImports() = %v, want %v", tt.name, string(gotData), string(tt.wantData))
		}
	}
}
