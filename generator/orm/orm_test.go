package orm

import (
	"testing"

	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/bdb/drivers"
)

func TestBuildRepeatableData(t *testing.T) {
	driver := &drivers.MySQLDriver{}
	type args struct {
		opts DataOpts
	}
	tests := []struct {
		name     string
		args     args
		wantData int
		wantErr  bool
	}{
		{"OK - no tables", args{DataOpts{}}, 0, false},
		{"OK - 1 table", args{DataOpts{Tables: Tables{bdb.Table{}}, DriverName: "mysql", Driver: driver}}, 1, false},
	}
	for _, tt := range tests {
		gotData, err := BuildRepeatableData(tt.args.opts)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. BuildRepeatableData() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}

		if len(gotData) != tt.wantData {
			t.Errorf("%q. BuildRepeatableData() len(data) = %v, want %v", tt.name, len(gotData), tt.wantData)
		}
	}
}

func TestBuildSingletonData(t *testing.T) {
	driver := &drivers.MySQLDriver{}
	type args struct {
		opts DataOpts
	}
	tests := []struct {
		name     string
		args     args
		wantData int
		wantErr  bool
	}{
		{"OK - 1 table", args{DataOpts{Tables: Tables{bdb.Table{}}, DriverName: "mysql", Driver: driver}}, 1, false},
	}
	for _, tt := range tests {
		gotData, err := BuildSingletonData(tt.args.opts)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. BuildSingletonData() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if len(gotData) != tt.wantData {
			t.Errorf("%q. BuildSingletonData() len(data) = %v, want %v", tt.name, len(gotData), tt.wantData)
		}
	}
}
