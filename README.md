![License](https://img.shields.io/badge/style-MIT-blue.svg?label=license)
![Build Status](https://api.travis-ci.org/yangyuqian/genus.svg?branch=master)

# genus

A simple wrapper of template engine for code generation tools in Go.

Code generation has been becoming a great idea to do metaprogramming in Go.
When you working on a code generation tool, there are many repeatable works,
such as
  * reading templates
  * handling text transformation, such as capitalization, camelcase
  * handling language specific transformation
  * formating generated code
  * post-processing after generation (fix imports)

With Genus, tools can do code generation in a clean way with a few lines of code.

## Features - V1

- [x] Template Wrapper

  - [x] loading templates from diretories, files and bytes easily
  - [x] grouping templates for Go packages
  - [x] formating, fixing imports after code generation

- [x] Generation Planner

  - [x] planning generation scenarios on template repository

- [x] Language specified helper funcs

  - [x] Golang helper funcs

- [x] Command-line Interface

  - [x] accept context in JSON and do generation without writing any code
