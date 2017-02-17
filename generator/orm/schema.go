package orm

import (
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/bdb/drivers"
	"github.com/yangyuqian/genus/types"
)

type Tables []bdb.Table

func NewSchema(user, pass, dbname, host string, port int, sslmode string) (s *Schema) {
	return &Schema{user, pass, dbname, host, port, sslmode, nil, nil, nil, nil}
}

type Schema struct {
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

func (s *Schema) init() (err error) {
	if s.driver == nil {
		drivers.TinyintAsBool = true
		s.driver = drivers.NewMySQLDriver(s.user, s.pass, s.name, s.host, s.port, s.sslmode)
	}
	opErr := s.driver.Open()
	if opErr != nil {
		return opErr
	}

	return
}

func (s *Schema) Close() {
	if s.driver != nil {
		s.driver.Close()
	}
}

func (s *Schema) buildTables(whitelist, blacklist []string) (tables []bdb.Table, err error) {
	s.tables, err = bdb.Tables(s.driver, s.name, whitelist, blacklist)
	if err != nil {
		return nil, err
	}

	if len(s.tables) > 0 {
		for i, t := range s.tables {
			if t.SchemaName == "" && s.name != "" {
				s.tables[i].SchemaName = s.name
			}
		}
	}

	return s.tables, nil
}

func (s *Schema) CollectTables(includeTables []string, blacklist []string) (tables Tables, err error) {
	if iErr := s.init(); iErr != nil {
		return nil, iErr
	}

	// Collect tables recursively through their foreign keys
	if incErr := s.include(s.name, includeTables...); incErr != nil {
		return nil, incErr
	}

	if _, uniqErr := s.uniqSortIncludeTables(); uniqErr != nil {
		return nil, uniqErr
	}

	// Build tables structures through collected tables
	if _, bErr := s.buildTables(s.includeTables, blacklist); bErr != nil {
		return nil, bErr
	}

	return s.tables, nil
}

func (s *Schema) SetDriver(driver bdb.Interface) {
	s.driver = driver
}

// Schema.include takes a list of tables in a schema to
// seek associated tables recursively through foreign keys
// Since it's a recursive algorithm, it's inefficient to uniq and sort
// on the includeTables on every recursive step.
// There may be duplicated tables in Schema.includeTables
func (s *Schema) include(schema string, tables ...string) (err error) {
	for _, t := range tables {
		s.includeTables = append(s.includeTables, t)
		fks, fksErr := s.driver.ForeignKeyInfo(schema, t)
		if fksErr != nil {
			return fksErr
		}

		if len(fks) > 0 {
			for _, fk := range fks {
				s.includeTables = append(s.includeTables, fk.ForeignTable)
				s.include(schema, fk.ForeignTable)
			}
		}
	}
	return
}

// Ensure tables are unique
func (s *Schema) uniqSortIncludeTables() (includes types.StringSlice, err error) {
	if s.includeTables, err = s.includeTables.Uniq(); err != nil {
		return
	}

	if s.includeTables, err = s.includeTables.Sort(); err != nil {
		return nil, err
	}

	return s.includeTables, nil
}

func DeadDriverByName(name string) (driver bdb.Interface, err error) {
	switch name {
	case "mysql":
		return &drivers.MySQLDriver{}, nil
	case "mock":
		return &drivers.MockDriver{}, nil
	default:
		return nil, errors.Errorf("Driver <%s> not supported", name)
	}

	return
}
