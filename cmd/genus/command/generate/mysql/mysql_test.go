package mysql

import (
	"testing"

	"flag"

	"github.com/yangyuqian/genus"
	cli "gopkg.in/urfave/cli.v1"
)

func TestBaseDir(t *testing.T) {
	type args struct {
		base      string
		framework string
	}
	tests := []struct {
		name         string
		args         args
		wantBasePath string
		wantErr      bool
	}{
		{"OK", args{"./a", "sqlboiler"}, "./a", false},
	}
	for _, tt := range tests {
		gotBasePath, err := BaseDir(tt.args.base, tt.args.framework)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. BaseDir() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if gotBasePath != tt.wantBasePath {
			t.Errorf("%q. BaseDir() = %v, want %v", tt.name, gotBasePath, tt.wantBasePath)
		}
	}
}

func TestInitGPS(t *testing.T) {
	app := cli.NewApp()
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.String("framework", "sqlboiler", "test framework")
	flagSet.String("base", "./", "test base")
	flagSet.String("template-dir", "./", "test template directory")
	ctx := cli.NewContext(app, flagSet, nil)

	type args struct {
		ctx            *cli.Context
		singletonData  []interface{}
		repeatableData []interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantSpec *genus.Spec
		wantErr  bool
	}{
		{"OK", args{ctx, nil, nil}, &genus.Spec{
			Suffix:      ".tpl",
			TemplateDir: "./",
			BaseDir:     "./",
			PlanItems: []*genus.PlanItem{
				&genus.PlanItem{},
				&genus.PlanItem{},
				&genus.PlanItem{},
				&genus.PlanItem{},
				&genus.PlanItem{},
			},
		}, false},
	}
	for _, tt := range tests {
		gotSpec, err := InitGPS(tt.args.ctx, tt.args.singletonData, tt.args.repeatableData)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. InitGPS() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if len(gotSpec.PlanItems) != len(tt.wantSpec.PlanItems) {
			t.Errorf("%q. InitGPS() = %+v, want %+v", tt.name, gotSpec, tt.wantSpec)
		}
	}
}
