Generation Plan Specification(GPS)
--------------------

GPS plays as a unified context in `genus`, containing all information used
in code generation, such as paths, settings and data.
Go to [JSON Schema](/cmd/genus/spec/plan_schema.json) to see more details.

GPS consists of some global settings, such as `TemplateDir`, `Suffix` and `PlanItems`.

In GPS, code generations are organized in plan items, which represents individual
generation strategy inside a go package.

# Global Settings

Global settings are settings shared across all plan items,
and for every single plan item, global settings can be overrided.

|Setting|Description|
|-------|-----------|
|Suffix|global suffix of templates files, i.e. '.tpl', will be overrided when plan items specifies its own suffix|
|TemplateDir|root directory of template files|
|BaseDir|location you want to put your generated code|
|BasePackge|base go package for local imports in generated code, i.e. github.com/someone/somerepo, default for the package of $PWD|
|SkipExists|set to true if you want to skip generation once file exists|
|SkipFormat|set to true if you want to skip go formating once file exists|
|SkipFixImports|set to true if you want to skip fixing imports|
|Merge|set to true if templates of one plan item should be merged into single file|
|Data|global data shared across all plan items, will be overrided if a plan item has its own data|

Note that data must be an array object.

# PlanItems

PlanItems represents independent code generation plan for files in same package.

|Setting|Description|
|-------|-----------|
|PlanType|specify the generation strategy, can either be 'SINGLETON' or 'REPEATABLE', set to REPEATABLE if plan is executed multiple times with every single item in given data, otherwise set it to SINGLETONE to execute once with 1st item in given data|
|Suffix|template suffix of specified plan item, i.e. '.tpl'|
|TemplateDir|location of templates|
|BaseDir|location of generated code|
|BasePackge|base go package for local imports in generated code|
|SkipExists|set to true if you want to skip generation once file exists|
|SkipFormat|set to true if you want to skip go formating once file exists|
|SkipFixImports|set to true if you want to skip fixing imports|
|Merge|set to true if templates of one plan item should be merged into single file|
|Data|global data shared across all plan items, will be overrided if a plan item has its own data|
|Package|package of generated code|
|RelativePackage|relative package of the generated code under the BaseDir|
|Filename|filename of your generated code, and it can be a Go template accepting the same context of current plan item|
|TemplateNames|specify template names to be included in currenct plan item. It's the relative path of template files under the TemplateDir with Suffix removed.|
|Imports|complete imports for your plan, can be absolute imports or local imports, key of the object is the actual path, and value of the object is alias of the import|
|Data|an array object, for a singleton plan item, first item will be used to execute templates, and for repeatable plan item, templates will be executed with items|

## Singleton



## Repeatable


