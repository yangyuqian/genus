package main

import "testing"

func Test_validateSpec(t *testing.T) {
	type args struct {
		specPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"OK", args{"testdata/spec/plan-success.json"}, false},
		{"KO - bad Suffix", args{"testdata/spec/plan-fail.json"}, true},
	}
	for _, tt := range tests {
		if err := validateSpec(tt.args.specPath); (err != nil) != tt.wantErr {
			t.Errorf("%q. validateSpec() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
