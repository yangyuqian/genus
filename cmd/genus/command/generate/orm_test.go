package generate

import (
	"testing"

	"flag"

	cli "gopkg.in/urfave/cli.v1"
)

func Test_doGenerateSqlboiler(t *testing.T) {
	app := cli.NewApp()
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	ctx := cli.NewContext(app, flagSet, nil)

	type args struct {
		ctx *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"OK", args{ctx}, true},
	}
	for _, tt := range tests {
		if err := doGenerateSqlboiler(tt.args.ctx); (err != nil) != tt.wantErr {
			t.Errorf("%q. doGenerateSqlboiler() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
