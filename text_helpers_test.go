package genus

import (
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/bdb"
)

func TestTxtToOne_Uniq(t *testing.T) {
	type fields struct {
		ForeignKey bdb.ForeignKey
		LocalTable struct {
			NameGo       string
			ColumnNameGo string
		}
		ForeignTable struct {
			NameGo       string
			NamePluralGo string
			ColumnNameGo string
			ColumnName   string
		}
		Function struct {
			Name              string
			ForeignName       string
			UsesBytes         bool
			LocalAssignment   string
			ForeignAssignment string
		}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"OK", fields{}, false},
	}
	for _, tt := range tests {
		to := &TxtToOne{
			ForeignKey:   tt.fields.ForeignKey,
			LocalTable:   tt.fields.LocalTable,
			ForeignTable: tt.fields.ForeignTable,
			Function:     tt.fields.Function,
		}
		if err := to.Uniq(); (err != nil) != tt.wantErr {
			t.Errorf("%q. TxtToOne.Uniq() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestTxtToMany_Uniq(t *testing.T) {
	type fields struct {
		LocalTable struct {
			NameGo       string
			ColumnNameGo string
		}
		ForeignTable struct {
			NameGo            string
			NamePluralGo      string
			NameHumanReadable string
			ColumnNameGo      string
			Slice             string
		}
		Function struct {
			Name              string
			ForeignName       string
			UsesBytes         bool
			LocalAssignment   string
			ForeignAssignment string
		}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"OK", fields{}, false},
	}
	for _, tt := range tests {
		tm := &TxtToMany{
			LocalTable:   tt.fields.LocalTable,
			ForeignTable: tt.fields.ForeignTable,
			Function:     tt.fields.Function,
		}
		if err := tm.Uniq(); (err != nil) != tt.wantErr {
			t.Errorf("%q. TxtToMany.Uniq() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_txtsFromToMany(t *testing.T) {
	type args struct {
		tables []bdb.Table
		table  bdb.Table
		rel    bdb.ToManyRelationship
	}
	tests := []struct {
		name string
		args args
		want TxtToMany
	}{
		{"OK", args{
			[]bdb.Table{},
			bdb.Table{
				Columns: []bdb.Column{
					bdb.Column{
						Name: "id",
					},
				},
			},
			bdb.ToManyRelationship{Column: "id"},
		}, TxtToMany{
			LocalTable: struct {
				NameGo       string
				ColumnNameGo string
			}{ColumnNameGo: "ID"},
			ForeignTable: struct {
				NameGo            string
				NamePluralGo      string
				NameHumanReadable string
				ColumnNameGo      string
				Slice             string
			}{Slice: "Slice"},
			Function: struct {
				Name              string
				ForeignName       string
				UsesBytes         bool
				LocalAssignment   string
				ForeignAssignment string
			}{LocalAssignment: "ID", ForeignName: "X"},
		}},
	}
	for _, tt := range tests {
		if got := txtsFromToMany(tt.args.tables, tt.args.table, tt.args.rel); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. txtsFromToMany() = %+v, want %+v", tt.name, got, tt.want)
		}
	}
}

func Test_mkFunctionName(t *testing.T) {
	type args struct {
		fkeyTableSingular    string
		foreignTablePluralGo string
		fkeyColumn           string
		toJoinTable          bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"OK", args{"user", "Users", "name", false}, "NameUsers"},
	}
	for _, tt := range tests {
		if got := mkFunctionName(tt.args.fkeyTableSingular, tt.args.foreignTablePluralGo, tt.args.fkeyColumn, tt.args.toJoinTable); got != tt.want {
			t.Errorf("%q. mkFunctionName() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_trimSuffixes(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"OK", args{"user_id"}, "user"},
		{"OK", args{"user"}, "user"},
	}
	for _, tt := range tests {
		if got := trimSuffixes(tt.args.str); got != tt.want {
			t.Errorf("%q. trimSuffixes() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
