Built-in Generators
--------------------

A built-in generator parses given metadata to Generation Plan Specification,
then perform the generation in the standard way.

# ORM Generator

ORM Generator parses a given database schema,
then creates orms and their correlated GPS.

For example, you can create models of gorm by performing

```
$ cd $GOPATH/github.com/yangyuqian/genus
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

Specify ORM framework through `--framework`, and generated orm will be located
at `${location_of_generated_code}` specified by `--base`
