package genus

import "testing"

func TestStringWithDefault(t *testing.T) {
	type args struct {
		v string
		d string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"OK", args{"", "d"}, "d"},
		{"OK", args{"v", "d"}, "v"},
		{"OK", args{"v", ""}, "v"},
	}
	for _, tt := range tests {
		if got := StringWithDefault(tt.args.v, tt.args.d); got != tt.want {
			t.Errorf("%q. StringWithDefault() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
