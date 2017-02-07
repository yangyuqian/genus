package orm

import (
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/bdb/drivers"
	"github.com/yangyuqian/genus/types"
)

func TestSchema_buildTables(t *testing.T) {
	type fields struct {
		user          string
		pass          string
		name          string
		host          string
		port          int
		sslmode       string
		driver        bdb.Interface
		tables        []bdb.Table
		includeTables types.StringSlice
		excludeTables types.StringSlice
	}
	type args struct {
		whitelist []string
		blacklist []string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantTables []bdb.Table
		wantErr    bool
	}{
		{"OK", fields{driver: &drivers.MockDriver{}}, args{[]string{"airports"}, []string{}}, Tables{
			bdb.Table{
				Name: "airports",
				Columns: []bdb.Column{
					{Name: "id", Type: "int", DBType: "integer"},
					{Name: "size", Type: "null.Int", DBType: "integer", Nullable: true},
				},
				PKey: &bdb.PrimaryKey{
					Name:    "airport_id_pkey",
					Columns: []string{"id"},
				},
			},
		}, false},
	}
	for _, tt := range tests {
		s := &Schema{
			user:          tt.fields.user,
			pass:          tt.fields.pass,
			name:          tt.fields.name,
			host:          tt.fields.host,
			port:          tt.fields.port,
			sslmode:       tt.fields.sslmode,
			driver:        tt.fields.driver,
			tables:        tt.fields.tables,
			includeTables: tt.fields.includeTables,
			excludeTables: tt.fields.excludeTables,
		}
		gotTables, err := s.buildTables(tt.args.whitelist, tt.args.blacklist)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Schema.buildTables() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotTables, tt.wantTables) {
			t.Errorf("%q. Schema.buildTables() = %v, want %v", tt.name, gotTables, tt.wantTables)
		}
	}
}

func TestSchema_include(t *testing.T) {
	type fields struct {
		user          string
		pass          string
		name          string
		host          string
		port          int
		sslmode       string
		driver        bdb.Interface
		tables        []bdb.Table
		includeTables types.StringSlice
		excludeTables types.StringSlice
	}
	type args struct {
		schema string
		tables []string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantErr           bool
		wantIncludeTables types.StringSlice
	}{
		{"OK", fields{driver: &drivers.MockDriver{}}, args{"schema1", []string{"jets"}}, false, types.StringSlice{"jets", "pilots", "pilots", "airports", "airports"}},
		{"OK", fields{driver: &drivers.MockDriver{}}, args{"schema1", []string{"jets", "pilot_languages"}}, false, types.StringSlice{"jets", "pilots", "pilots", "airports", "airports", "pilot_languages", "pilots", "pilots", "languages", "languages"}},
	}
	for _, tt := range tests {
		s := &Schema{
			user:          tt.fields.user,
			pass:          tt.fields.pass,
			name:          tt.fields.name,
			host:          tt.fields.host,
			port:          tt.fields.port,
			sslmode:       tt.fields.sslmode,
			driver:        tt.fields.driver,
			tables:        tt.fields.tables,
			includeTables: tt.fields.includeTables,
			excludeTables: tt.fields.excludeTables,
		}
		if err := s.include(tt.args.schema, tt.args.tables...); (err != nil) != tt.wantErr {
			t.Errorf("%q. Schema.include() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}

		if !reflect.DeepEqual(s.includeTables, tt.wantIncludeTables) {
			t.Errorf("%q. Schema.include() includeTables = %v, wantIncludeTables %v", tt.name, s.includeTables, tt.wantIncludeTables)
		}
	}
}

func TestSchema_CollectTables(t *testing.T) {
	type fields struct {
		user          string
		pass          string
		name          string
		host          string
		port          int
		sslmode       string
		driver        bdb.Interface
		tables        []bdb.Table
		includeTables types.StringSlice
		excludeTables types.StringSlice
	}
	type args struct {
		includeTables []string
		blacklist     []string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantTables Tables
		wantErr    bool
	}{
		{"OK", fields{driver: &drivers.MockDriver{}}, args{[]string{"airports"}, []string{}}, Tables{
			bdb.Table{
				Name: "airports",
				Columns: []bdb.Column{
					{Name: "id", Type: "int", DBType: "integer"},
					{Name: "size", Type: "null.Int", DBType: "integer", Nullable: true},
				},
				PKey: &bdb.PrimaryKey{
					Name:    "airport_id_pkey",
					Columns: []string{"id"},
				},
			},
		}, false},
	}
	for _, tt := range tests {
		s := &Schema{
			user:          tt.fields.user,
			pass:          tt.fields.pass,
			name:          tt.fields.name,
			host:          tt.fields.host,
			port:          tt.fields.port,
			sslmode:       tt.fields.sslmode,
			driver:        tt.fields.driver,
			tables:        tt.fields.tables,
			includeTables: tt.fields.includeTables,
			excludeTables: tt.fields.excludeTables,
		}
		gotTables, err := s.CollectTables(tt.args.includeTables, tt.args.blacklist)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Schema.CollectTables() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotTables, tt.wantTables) {
			t.Errorf("%q. Schema.CollectTables() = %v, want %v", tt.name, gotTables, tt.wantTables)
		}
	}
}

func TestSchema_uniqSortIncludeTables(t *testing.T) {
	type fields struct {
		user          string
		pass          string
		name          string
		host          string
		port          int
		sslmode       string
		driver        bdb.Interface
		tables        []bdb.Table
		includeTables types.StringSlice
		excludeTables types.StringSlice
	}
	tests := []struct {
		name         string
		fields       fields
		wantIncludes types.StringSlice
		wantErr      bool
	}{
		{"OK", fields{includeTables: types.StringSlice{"a", "b", "a"}}, types.StringSlice{"a", "b"}, false},
		{"KO", fields{includeTables: types.StringSlice{"a", "a", "a"}}, types.StringSlice{"a"}, false},
		{"KO", fields{includeTables: types.StringSlice{"a", "b"}}, types.StringSlice{"a", "b"}, false},
	}
	for _, tt := range tests {
		s := &Schema{
			user:          tt.fields.user,
			pass:          tt.fields.pass,
			name:          tt.fields.name,
			host:          tt.fields.host,
			port:          tt.fields.port,
			sslmode:       tt.fields.sslmode,
			driver:        tt.fields.driver,
			tables:        tt.fields.tables,
			includeTables: tt.fields.includeTables,
			excludeTables: tt.fields.excludeTables,
		}
		gotIncludes, err := s.uniqSortIncludeTables()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Schema.uniqSortIncludeTables() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotIncludes, tt.wantIncludes) {
			t.Errorf("%q. Schema.uniqSortIncludeTables() = %v, want %v", tt.name, gotIncludes, tt.wantIncludes)
		}
	}
}

func TestDeadDriverByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name       string
		args       args
		wantDriver bdb.Interface
		wantErr    bool
	}{
		{"OK", args{"mysql"}, &drivers.MySQLDriver{}, false},
		{"OK", args{"unknown"}, nil, true},
	}
	for _, tt := range tests {
		gotDriver, err := DeadDriverByName(tt.args.name)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. DeadDriverByName() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotDriver, tt.wantDriver) {
			t.Errorf("%q. DeadDriverByName() = %v, want %v", tt.name, gotDriver, tt.wantDriver)
		}
	}
}
