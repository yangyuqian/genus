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

func Test_once_Has(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		o    once
		args args
		want bool
	}{
		{"OK - has => true", once{"a": struct{}{}}, args{"a"}, true},
		{"OK - has => false", once{"a": struct{}{}}, args{"b"}, false},
	}
	for _, tt := range tests {
		if got := tt.o.Has(tt.args.s); got != tt.want {
			t.Errorf("%q. once.Has() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_once_Put(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		o    once
		args args
		want bool
	}{
		{"OK - true", once{"a": struct{}{}}, args{"a"}, false},
		{"OK - false", once{"a": struct{}{}}, args{"b"}, true},
	}
	for _, tt := range tests {
		if got := tt.o.Put(tt.args.s); got != tt.want {
			t.Errorf("%q. once.Put() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestBoolWithDefault(t *testing.T) {
	type args struct {
		v bool
		d bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"OK - true + true", args{true, true}, true},
		{"OK - false + true", args{false, true}, true},
	}
	for _, tt := range tests {
		if got := BoolWithDefault(tt.args.v, tt.args.d); got != tt.want {
			t.Errorf("%q. BoolWithDefault() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
