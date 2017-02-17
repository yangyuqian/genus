Examples
-------------------

# ORM Generator of GORM

You can perform the orm generation with following commands

```
$ cd $GOPATH/github.com/yangyuqian/genus
$ genus generate --spec examples/orm/gorm/plan.json
```

Generated models are located at `_test/gorm`, structs, fields, tags and other
methods are created with given context in your GPS.

```
// _test/gorm/orm/ants.go
package orm

// Ant is an object representing the database table.
type Ant struct {
	ID      []byte `gorm:"column:id" json:"id" boil:"id"`
	Name    []byte `gorm:"column:name" json:"name" boil:"name"`
	TigerID []byte `gorm:"column:tiger_id" json:"tiger_id" boil:"tiger_id"`
}

func (obj *Ant) TableName() string {
	return "ants"
}

type AntSlice []*Ant

// _test/gorm/orm/bro.go
package orm

import null "gopkg.in/nullbio/null.v6"

// Bro is an object representing the database table.
type Bro struct {
	Bros string      `gorm:"column:bros" json:"bros" boil:"bros"`
	Name null.String `gorm:"column:name" json:"name,omitempty" boil:"name"`
}

func (obj *Bro) TableName() string {
	return "bro"
}

type BroSlice []*Bro
```

Go to [Generation Plan Specification](gps.md) for more details.

# Parse MySQL database schema and generate models for GORM

You may want to parse a given database schema and generate a orm.

In the case of models of GORM, you can do it directly with `genus`:

```
$ genus g orm --host ${your_mysql_host} \
              --port ${your_mysql_port} \
              --username ${your_mysql_username} \
              --password ${your_mysql_password} \
              --database ${your_mysql_database} \
              --framework gorm \
              --mysql \
              --base ${location_of_generated_code} \
              --relative-pkg gorm/orm
```

This command will generate ORM through given database schema,
and save the context to `${location_of_generated_code}/plan.json`.
With the generated specification, you can perform the standard generation
without touching databases to debug your templates effiently.

Generated orm will be located at ${location_of_generated_code},
specified by `--base`. See the usage by performing `genus g orm --help`
