Guideline of Flexible Templates in Go
-------------------

Thanks to the built-in template engine, many generators,
such as [SQLBoiler](https://github.com/vattle/sqlboiler)
and [Go Swagger](https://github.com/go-swagger/go-swagger),
can parse input to metadata, then output fully functional Go code easily.

However, when you want to extend those generators, like adding some features to
their templates, you will have to go through their implementation, which
prevents many people from using code generation to empower themselves.

After reviewing many popular generators, I found there are many common
use cases and scenarios to perform code generation. Also, it's possible
to organize the code generation in a clean way to provide a flexible generator.

This guideline shows some useful practices when you want to do code generation
with templates:

* always use explicit metadata to do code generation
* use `with` to protect context
* use variables for names
* break big templates into smaller ones
* never use embedded templates

# always use explicit metadata to do code generation

Using explict metadata means that it should be possible to figure the input
context that used to execute templates. In another word, the context of your
templates should be exposed to users directly.

For example, when generating ORM against a given database schema, the metadata
should be the parsed database schema.

```
{
  "Tables": [
    {
      "Columns": [
        {"Column1" ... }
        ...
      ]
    }
  ]
}
```

# use `with` to protect context

`with` wraps a block with it's independent context,
and also it checks the input context is zero value, the block will be executed
only when input context is non-zero value.

For example, consider a template

```
My name is {{ .Name }}
My son is {{ .Son.Name }}
```

with context of

```
{
  "Name": "Bob",
  "Son": null
}
```

null pointer error occurs to execute `{{ .Son.Name }}` when `Son` is null,
to get rid of this, you may need to add a condition block before executing
`{{ .Son.Name }}`

```
My name is {{ .Name }}
{{ if notZero .Son }}
My son is {{ .Son.Name }}
{{ end }}
```

How about to have a grand son in the context? let's add more complexity to
above template:

```
My name is {{ .Name }}
My son is {{ .Son.Name }}
My grand son is {{ .Son.Grandson.Name }}
```

To prevent it from the NPE, you must add if block to all potential nodes

```
My name is {{ .Name }}
{{ if notZero .Son }}
My son is {{ .Son.Name }}
{{ if notZero .Son.Grandson }}
My grand son is {{ .Son.Grandson.Name }}
{{ end }}
{{ end }}
```

It will be too complex to maintain when you have a context in real world,
like database schema.
Instead of using `if` block to protect templates from NPE, you can use `with`
to simply the context and protect it from NPE

```
My name is {{ .Name }}
{{ with .Son }}
My son is {{ .Name }}
{{ with .Grandson }}
My grand son is {{ .Name }}
{{ end }}
{{ end }}
```

# use variables for names

Do not combine names in the templates, instead, use variables for all names.

For example, consider following Go code

```
var catMsg := "Hello, World"
println(catMsg)
```

Suppose you want to make the animal name configurable in generator, a intuitive
template would be

```
var {{ .AnimalName }}Msg := "Hello, World"
println({{ .AnimalName }}Msg)
```

with context of

```
{
  "AnimalName": "cat"
}
```

But when you want to change the name to `{{ .AnimalName }}Words`, you'll have
to update the template everywhere manually.

To get rid of this kind of issue, simply use variables for all names

```
{{ $animalMsg := printf "%sMsg" .AnimalName }}
var {{ $animalMsg }} := "Hello, World"
println({{ $animalMsg }})
```

# break big templates into smaller ones

Big templates are always unreadable to its maintainers,
just break them into smaller ones with no more than one hundred lines.

# never use embedded templates

Embedded templates defines a reuseable template and can be used in other templates.

```
{{ define "tmpl1" }}Meta-Programming{{ end }}

I love {{ template "tmpl1" }}
```

In practice, embedded templates are hard to understand and difficult to debug,
just don't use it in your generator.

Futhermore, if there are embedded templates everywhere, it always means there
are duplicated logic in the generated code. Remove those duplicated logic first.

Again, don't use embedded templates.

