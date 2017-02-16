Generation Plan Specification(GPS)
--------------------

GPS plays as a unified context in `genus`, containing all information used
in code generation, such as paths, settings and data.

Go to [JSON Schema](/cmd/genus/spec/plan_schema.json) to see more details.

GPS consists of some global settings, such as `TemplateDir`, `Suffix` and `PlanItems`.

Also, you can set the `Merge` to `true` to merge the generated code to one single file.

# Global Settings


# PlanItems

PlanItems represents independent code generation plan for files in same package.

For every single plan item, `PlanType` specifies its generation strategy,
it can either be `SINGLETON` or `REPEATABLE`.

## Singleton

When `PlanType` set to `SINGLETON`, code will be generated with `Data`,
it's either be global data or `Data` field in plan item.

For example,

```
# plan.json
{
  "PlanItems": [
    {
      "PlanType": "SINGLETON",
      "TemplateDir": "./",
      "RelativePackage": "myapp",
      "Base": "./_test",
      "TemplateNames": [
        "a",
        "b",
        "c"
      ],
      "Data": {
        ...
      }
    }
  ]
}
```

Above GPS creates following files:

```
./_test
  |- myapp
    |- a.go
    |- b.go
    |- c.go
```

Filenames of your generated code can be customized through set `Filename` to
template with input data of currect plan item:

```
# plan.json
{
  "PlanItems": [
    {
      "PlanType": "SINGLETON",
      "Filename": "{{ .AnyFilenameYouLike }}.go",
      "TemplateDir": "./",
      "RelativePackage": "myapp",
      "Base": "./_test",
      "TemplateNames": [
        "a",
        "b",
        "c"
      ],
      "Data": {
        ...
      }
    }
  ]
}
```

Also, you can merge those 3 files into one single file through setting `Merge`
to true:

```
# plan.json
{
  "PlanItems": [
    {
      "PlanType": "SINGLETON",
      "Filename": "merged.go",
      "Merge": true,
      "TemplateDir": "./",
      "RelativePackage": "myapp",
      "Base": "./_test",
      "TemplateNames": [
        "a",
        "b",
        "c"
      ],
      "Data": {
        ...
      }
    }
  ]
}
```

Merged code contains only 1 file:

```
./_test
  |- myapp
    |- merged.go
```


## Repeatable


