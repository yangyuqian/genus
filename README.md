![License](https://img.shields.io/badge/style-MIT-blue.svg?label=license)
![Build Status](https://api.travis-ci.org/yangyuqian/genus.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/yangyuqian/genus)](https://goreportcard.com/report/github.com/yangyuqian/genus)

[中文版](README-cn.md)

List of Contents:

- [Examples](docs/en/examples.md)

- [Generation Plan Specification](docs/en/gps.md)

- [Guideline of Flexible Templates in Go](docs/en/flexible-templates.md)

- [Built-in Generators](docs/en/generators.md)

# genus

A general code generation tools in Go.

New to code generation?
Refer [text/template](https://golang.org/pkg/text/template/) for basic knowledge.

# Why Another Generator?

Code generation has been becoming a popular way of metaprogramming in Go.

For example, when working with a ORM framework, you may want to create models
with a given database schema with json tag support.
it's trivial and burdensome when it comes with a database
with hundreds even thousands of tables.

A code generation tool can help on this kind of issues through retriving
database schema and rendering templates to go code.

However, before it turns into reality, you'll have to build a tool to

* reading and organizing templates properly
* formating the generated source code
* fixing imports, especially removing unused imports
* handle relative imports(imports among the generated code)
* naming of generated file names

Today, generators are handling code generation in their own ways, which means
that you can not generate models of
[Beego ORM](https://beego.me/docs/mvc/model/orm.md) with generators or
[gorm](http://jinzhu.me/gorm). It doesn't make sense to build a new generator
becuase those models are generated from the same database schema.

Genus provides a clean way to perform code generation.
Go to [Examples](docs/en/examples.md) to see more details.

# Installation

Genus has built-in CLI support, you can install it by performing

```
go get -u github.com/yangyuqian/genus/cmd/genus
```

# Features - V1

- [x] Template Wrapper

  - [x] loading templates from diretories, files and bytes easily
  - [x] grouping templates for Go packages
  - [x] formating, fixing imports after code generation

- [x] Generation Planner

  - [x] planning generation scenarios on template repository
  - [x] planner for regular files(non-Golang)

- [x] Language specified helper funcs

  - [x] Golang helper funcs

- [x] Generation Plan Specification

  - [x] JSON schema and validation support
  - [x] Reference and pointer support for complex specfications

- [x] Command-line Interface

  - [x] accept context in JSON and do generation without writing any code
  - [x] create Generation Plan Specification of Database from given database


