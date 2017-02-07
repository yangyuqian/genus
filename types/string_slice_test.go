package types

import (
	"reflect"
	"testing"
)

func TestStringSlice_Uniq(t *testing.T) {
	tests := []struct {
		name    string
		slice   StringSlice
		wantO   StringSlice
		wantErr bool
	}{
		{"OK", StringSlice{"a", "b", "a"}, StringSlice{"a", "b"}, false},
	}
	for _, tt := range tests {
		got, err := tt.slice.Uniq()
		gotO, err := got.Sort()

		if (err != nil) != tt.wantErr {
			t.Errorf("%q. StringSlice.Uniq() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotO, tt.wantO) {
			t.Errorf("%q. StringSlice.Uniq() = %v, want %v", tt.name, gotO, tt.wantO)
		}
	}
}

func TestStringSlice_Sort(t *testing.T) {
	tests := []struct {
		name    string
		slice   StringSlice
		wantO   StringSlice
		wantErr bool
	}{
		{"OK", StringSlice{"b", "a", "c"}, StringSlice{"a", "b", "c"}, false},
	}
	for _, tt := range tests {
		gotO, err := tt.slice.Sort()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. StringSlice.Sort() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotO, tt.wantO) {
			t.Errorf("%q. StringSlice.Sort() = %v, want %v", tt.name, gotO, tt.wantO)
		}
	}
}
